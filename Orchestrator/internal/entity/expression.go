package entity

import (
	"errors"
	"strings"
	"time"
)

const (
	errInvalidExpression = "invalid expression"
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

func NewExpression(exp string,
	calcExp string,
	validator func(string) bool,
	user string,
	start time.Time,
	updater func([]byte) []byte,
) *Expression {
	exp = strings.ReplaceAll(exp, " ", "")
	newExp := &Expression{
		Expression:           exp,
		CalculatedExpression: string(updater([]byte(calcExp))),
		Start:                start,
		User:                 user,
	}
	if !validator(exp) {
		newExp.Err = errors.New(errInvalidExpression)
	}
	return newExp
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
