package polandNotation

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/adminsemy/yandexCalculator/Agent/intenal/config"
)

var (
	ErrValidation   = errors.New("неверное выражение")
	ErrDivideByZero = errors.New("на ноль делить нельзя")
)

type PolandNotation struct {
	expression string
	result     int64
	config     *config.ConfigExpression
}

func New(expression string, config *config.ConfigExpression) *PolandNotation {
	return &PolandNotation{
		expression: expression,
		config:     config,
	}
}

func (p *PolandNotation) Calculate() error {
	stack := make([]float64, 0)
	str := strings.Split(p.expression, " ")
	slog.Info("Старт вычисления выражения", "выражение", p.expression)
	for _, v := range str {
		if v == "+" || v == "-" || v == "*" || v == "/" {
			if len(stack) < 2 {
				return ErrValidation
			}
			right := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			left := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			switch v {
			case "+":
				time.Sleep(time.Duration(p.config.Plus) * time.Minute)
				stack = append(stack, left+right)
			case "-":
				time.Sleep(time.Duration(p.config.Minus) * time.Minute)
				stack = append(stack, left-right)
			case "*":
				time.Sleep(time.Duration(p.config.Multiply) * time.Minute)
				stack = append(stack, left*right)
			case "/":
				if right == 0 {
					return ErrDivideByZero
				}
				time.Sleep(time.Duration(p.config.Divide) * time.Minute)
				stack = append(stack, left/right)
			}
		} else {
			number, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return err
			}
			stack = append(stack, number)
		}
	}
	if len(stack) != 1 {
		return ErrValidation
	}
	p.result = int64(stack[0])
	slog.Info("Конец вычисления выражения", "результат", p.result)
	return nil
}

func (p *PolandNotation) Result() int64 {
	return p.result
}
