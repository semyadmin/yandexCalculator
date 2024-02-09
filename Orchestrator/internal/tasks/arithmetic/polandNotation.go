package arithmetic

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/validator"
)

type Expression interface {
	GetExpression() string
	SetID(uint64)
	Result() []string
}

var (
	ErrValidation     error = errors.New("Ошибка валидации данных")
	errNullExpression error = errors.New("Выражение не может быть пустым")

	ErrDivideByZero = errors.New("на ноль делить нельзя")
)

type PolandNotation struct {
	Expression string
	result     []string
	id         uint64
	Err        error
	config     *config.Config
	send       chan SendInfo
	get        chan string
}

type SendInfo struct {
	Id       uint64
	Result   string
	Deadline uint64
}

func NewPolandNotation(expression string, config *config.Config) (*PolandNotation, error) {
	slog.Info("Получены данные для польской нотации:", "expression", expression)
	if expression == "" {
		return nil, ErrValidation
	}
	p := &PolandNotation{
		Expression: expression,
		config:     config,
		send:       make(chan SendInfo),
		get:        make(chan string),
	}
	err := p.createResult()
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *PolandNotation) createResult() error {
	str := strings.ReplaceAll(p.Expression, " ", "")
	if !validator.IsValid(str) {
		return ErrValidation
	}
	result := make([]string, 0, len(str))
	stack := make([]byte, 0, 10)
	operators := map[byte]uint8{'(': 0, ')': 0, '+': 1, '-': 1, '*': 2, '/': 2}
	numbers := map[byte]struct{}{'0': {}, '1': {}, '2': {}, '3': {}, '4': {}, '5': {}, '6': {}, '7': {}, '8': {}, '9': {}}
	index := 0
	indexNumber := 0
	opened := 0
	closed := 0

	for index < len(str) {
		if _, ok := numbers[str[index]]; ok {
			indexNumber = index
			for indexNumber < len(str) {
				if _, ok := numbers[str[indexNumber]]; !ok {
					break
				}
				indexNumber++
			}
			result = append(result, str[index:indexNumber])
			index = indexNumber
			continue
		}
		operator := str[index]
		if operator == '(' {
			opened++
			stack = append(stack, operator)
			index++
			continue
		}
		if operator == ')' {
			closed++
			if closed > opened {
				return ErrValidation
			}
			for i := len(stack) - 1; i >= 0; i-- {
				if stack[i] == '(' {
					stack = stack[:i]
					break
				}
				result = append(result, string(stack[i]))
				stack = stack[:i]
			}
			index++
			continue
		}
		for i := len(stack) - 1; i >= 0; i-- {
			if operators[operator] >= operators[stack[i]] {
				break
			}
			result = append(result, string(stack[i]))
			stack = stack[:i]
		}
		stack = append(stack, operator)
		index++
	}
	if opened != closed {
		return ErrValidation
	}
	for i := len(stack) - 1; i >= 0; i-- {
		result = append(result, string(stack[i]))
		stack = stack[:i]
	}
	p.result = result
	return nil
}

func (p *PolandNotation) GetExpression() string {
	return p.Expression
}

func (p *PolandNotation) SetID(id uint64) {
	p.id = id
}

func (p *PolandNotation) Result() {
	stack := make([]float64, 0, 10)
	result := SendInfo{Id: p.id}
	for _, v := range p.result {
		if v == "+" || v == "-" || v == "*" || v == "/" {
			if len(stack) < 2 {
				p.Err = ErrValidation
				return
			}
			right := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			left := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			switch v {
			case "+":
				result.Result = strconv.FormatFloat(left, 'f', -1, 64) + " + " + strconv.FormatFloat(right, 'f', -1, 64)
				result.Deadline = uint64(p.config.Plus)
			case "-":
				result.Result = strconv.FormatFloat(left, 'f', -1, 64) + " - " + strconv.FormatFloat(right, 'f', -1, 64)
				result.Deadline = uint64(p.config.Minus)
			case "*":
				result.Result = strconv.FormatFloat(left, 'f', -1, 64) + " * " + strconv.FormatFloat(right, 'f', -1, 64)
				result.Deadline = uint64(p.config.Multiply)
			case "/":
				if right == 0 {
					p.Err = ErrDivideByZero
					return
				}
				result.Result = strconv.FormatFloat(left, 'f', -1, 64) + " / " + strconv.FormatFloat(right, 'f', -1, 64)
				result.Deadline = uint64(p.config.Divide)
			}
		} else {
			number, err := strconv.ParseFloat(v, 64)
			if err != nil {
				p.Err = err
				return
			}
			stack = append(stack, number)
			continue
		}
		// Отправляем данные для подсчета результата
		p.send <- result
		result := <-p.get
		number, err := strconv.ParseFloat(result, 64)
		if err != nil {
			p.Err = err
			return
		}
		stack = append(stack, number)
	}
	if len(stack) != 1 {
		p.Err = ErrValidation
		return
	}
}
