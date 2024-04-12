package entity

import (
	"errors"
	"strings"
	"time"
)

type Expression struct {
	ID                   uint64
	Start                time.Time
	Duration             int64
	Expression           string
	CalculatedExpression string
	Result               float64
	IsCalc               bool
	Err                  error
	User                 string
}

func NewExpression(exp string, calcExp string, validator func(string) bool, user string) *Expression {
	exp = strings.ReplaceAll(exp, " ", "")
	if !validator(exp) {
		return &Expression{
			Expression: exp,
			Err:        errors.New("invalid expression"),
			Start:      time.Now(),
			User:       user,
		}
	}
	return &Expression{
		Expression:           exp,
		CalculatedExpression: calcExp,
		Start:                time.Now(),
		User:                 user,
	}
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
