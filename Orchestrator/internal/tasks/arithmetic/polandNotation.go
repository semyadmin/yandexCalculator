package arithmetic

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/validator"
)

type Expression interface {
	GetExpression() string
	SetID(uint64)
	Result() []string
}

var (
	errValidation     error = errors.New("Ошибка валидации данных")
	errNullExpression error = errors.New("Выражение не может быть пустым")
)

type PolandNotation struct {
	expression string
	result     []string
	id         uint64
}

func NewPolandNotation(expression string) (*PolandNotation, error) {
	slog.Info("Получены данные для польской нотации:", "expression", expression)
	if expression == "" {
		return nil, errValidation
	}
	p := &PolandNotation{
		expression: expression,
	}
	err := p.createResult()
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *PolandNotation) createResult() error {
	str := strings.ReplaceAll(p.expression, " ", "")
	if !validator.IsValid(str) {
		return errValidation
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
				return errValidation
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
		return errValidation
	}
	for i := len(stack) - 1; i >= 0; i-- {
		result = append(result, string(stack[i]))
		stack = stack[:i]
	}
	p.result = result
	return nil
}

func (p *PolandNotation) GetExpression() string {
	return p.expression
}

func (p *PolandNotation) SetID(id uint64) {
	p.id = id
}

func (p *PolandNotation) Result() []string {
	return p.result
}
