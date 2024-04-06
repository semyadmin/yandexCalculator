package entity

import (
	"errors"
	"strings"
)

type Expression struct {
	ID         uint64
	Expression string
	Result     float64
	IsCalc     bool
	Err        error
}

func NewExpression(exp string, validator func(string) bool) *Expression {
	exp = strings.ReplaceAll(exp, " ", "")
	if !validator(exp) {
		return &Expression{Expression: exp, Err: errors.New("invalid expression")}
	}
	return &Expression{Expression: exp}
}

func (e *Expression) SetResult(r float64, err error) {
	if err != nil {
		e.Err = err
		return
	}
	e.Result = r
	e.IsCalc = true
}

func (e *Expression) SetId(id uint64) {
	e.ID = id
}
