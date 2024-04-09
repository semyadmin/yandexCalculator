package calculate

import (
	"errors"
	"strconv"
	"strings"
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

func Calculate(expression string) (string, error) {
	array := strings.Split(expression, " ")
	if len(array) != 3 {
		return "", errors.New("неверное выражение")
	}
	id := array[0]
	first, second, operation, err := parseExpression(array[1])
	if err != nil {
		return id, err
	}
	duration, err := strconv.ParseFloat(array[2], 64)
	if err != nil {
		return id, err
	}
	result := float64(0)
	switch operation {
	case "+":
		time.Sleep(time.Duration(duration) * time.Second)
		result = first + second
	case "-":
		time.Sleep(time.Duration(duration) * time.Second)
		result = first - second
	case "*":
		time.Sleep(time.Duration(duration) * time.Second)
		result = first * second
	case "/":
		time.Sleep(time.Duration(duration) * time.Second)
		if second == 0 {
			return id, errors.New("делить на ноль нельзя")
		}
		result = first / second
	}
	return id + " " + strconv.FormatFloat(result, 'f', -1, 64) + " " + array[2], nil
}

// Проверяем корректность значения
func parseExpression(str string) (float64, float64, string, error) {
	split := 0
	for i := 1; i < len(str); i++ {
		if str[i] == '+' || str[i] == '-' || str[i] == '*' || str[i] == '/' {
			split = i
			break
		}
	}
	first, err := strconv.ParseFloat(str[:split], 64)
	if err != nil {
		return 0, 0, "", err
	}
	second, err := strconv.ParseFloat(str[split+1:], 64)
	if err != nil {
		return 0, 0, "", err
	}
	return first, second, string(str[split]), nil
}
