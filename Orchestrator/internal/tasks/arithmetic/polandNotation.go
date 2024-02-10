package arithmetic

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
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
	queue      *queue.MapQueue
	get        chan string
	stack      []float64
	answer     float64
	isAnswer   bool
}

func NewPolandNotation(expression string, config *config.Config, queue *queue.MapQueue) (*PolandNotation, error) {
	slog.Info("Получены данные для польской нотации:", "expression", expression)
	if expression == "" {
		return nil, ErrValidation
	}
	p := &PolandNotation{
		Expression: expression,
		config:     config,
		queue:      queue,
		get:        make(chan string),
		stack:      make([]float64, 0, 10),
	}

	return p, nil
}

func (p *PolandNotation) CreateResult() error {
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

func (p *PolandNotation) Calculate() {
	result := &queue.SendInfo{Id: p.id}
	for _, v := range p.result {
		if v == "+" || v == "-" || v == "*" || v == "/" {
			if len(p.stack) < 2 {
				p.Err = ErrValidation
				return
			}
			right := p.stack[len(p.stack)-1]
			p.stack = p.stack[:len(p.stack)-1]
			left := p.stack[len(p.stack)-1]
			p.stack = p.stack[:len(p.stack)-1]
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
			p.stack = append(p.stack, number)
			continue
		}
		// Отправляем данные для подсчета результата
		p.queue.Enqueue(result)
		result := <-p.get
		number, err := strconv.ParseFloat(result, 64)
		if err != nil {
			p.Err = err
			return
		}
		p.stack = append(p.stack, number)
	}
	if len(p.stack) != 1 {
		p.Err = ErrValidation
		return
	}
	p.answer = p.stack[0]
	p.isAnswer = true
}
