package calculate

import (
	"time"

	"github.com/adminsemy/yandexCalculator/Agent/intenal/entity/expression"
)

type Expression interface {
	Id() string
	First() float64
	Second() float64
	Operation() string
	Duration() uint64
	Result() float64
	Error() string
	SetResult(float64)
	SetError(string)
}

// Вычисляем операцию в зависимости от оператора
// и времени на его обработку
func CalculateGRPC(exp Expression) Expression {
	first := exp.First()
	second := exp.Second()
	duration := exp.Duration()
	id := exp.Id()
	operation := exp.Operation()
	newExp := expression.New(id, first, second, operation, duration)
	switch operation {
	case "+":
		time.Sleep(time.Duration(int64(duration)) * time.Second)
		newExp.SetResult(first + second)
	case "-":
		time.Sleep(time.Duration(int64(duration)) * time.Second)
		newExp.SetResult(first - second)
	case "*":
		time.Sleep(time.Duration(int64(duration)) * time.Second)
		newExp.SetResult(first * second)
	case "/":
		time.Sleep(time.Duration(int64(duration)) * time.Second)
		if second == 0 {
			newExp.SetError("делить на ноль нельзя")
		} else {
			newExp.SetResult(first / second)
		}

	}
	return newExp
}
