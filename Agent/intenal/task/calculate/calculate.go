package calculate

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func Calculate(expression string) (string, error) {
	array := strings.Split(expression, " ")
	if len(array) != 3 {
		return "", errors.New("неверное выражение")
	}
	id := array[0]
	first, second, operation, err := parseExpression(array[1])
	if err != nil {
		return "", err
	}
	duration, err := strconv.ParseFloat(array[2], 64)
	if err != nil {
		return "", err
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
		if second == 0 {
			return "", errors.New("делить на ноль нельзя")
		}
		time.Sleep(time.Duration(duration) * time.Second)
		result = first / second
	}
	return id + " " + strconv.FormatFloat(result, 'f', -1, 64) + " " + array[2], nil
}

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
