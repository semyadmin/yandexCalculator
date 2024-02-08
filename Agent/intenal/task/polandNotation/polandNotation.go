package polandNotation

import (
	"errors"
	"strconv"
	"time"

	"github.com/adminsemy/yandexCalculator/Agent/intenal/config"
)

var (
	ErrValidation   = errors.New("неверное выражение")
	ErrDivideByZero = errors.New("на ноль делить нельзя")
)

type PolandNotation struct {
	ID         string
	Expression []string
	Result     int64
	config     *config.ConfigExpression
	Stack      []float64
	Err        error
}

func New(id string, expression []string, config *config.ConfigExpression) *PolandNotation {
	return &PolandNotation{
		ID:         id,
		Expression: expression,
		config:     config,
		Stack:      make([]float64, 0),
		Err:        nil,
	}
}

func (p *PolandNotation) Calculate(v string) {
	if v == "+" || v == "-" || v == "*" || v == "/" {
		if len(p.Stack) < 2 {
			p.Err = ErrValidation
			return
		}
		right := p.Stack[len(p.Stack)-1]
		p.Stack = p.Stack[:len(p.Stack)-1]
		left := p.Stack[len(p.Stack)-1]
		p.Stack = p.Stack[:len(p.Stack)-1]
		switch v {
		case "+":
			time.Sleep(time.Duration(p.config.Plus) * time.Minute)
			p.Stack = append(p.Stack, left+right)
		case "-":
			time.Sleep(time.Duration(p.config.Minus) * time.Minute)
			p.Stack = append(p.Stack, left-right)
		case "*":
			time.Sleep(time.Duration(p.config.Multiply) * time.Minute)
			p.Stack = append(p.Stack, left*right)
		case "/":
			if right == 0 {
				p.Err = ErrDivideByZero
				return
			}
			time.Sleep(time.Duration(p.config.Divide) * time.Minute)
			p.Stack = append(p.Stack, left/right)
		}
	} else {
		number, err := strconv.ParseFloat(v, 64)
		if err != nil {
			p.Err = err
			return
		}
		p.Stack = append(p.Stack, number)
	}
}
