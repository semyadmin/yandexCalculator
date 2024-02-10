package calculate

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func Calculate(expression string) (string, error) {
	array := strings.Split(expression, " ")
	id := array[0]
	first, err := strconv.ParseFloat(array[1], 64)
	if err != nil {
		return "", err
	}
	operation := array[2]
	second, err := strconv.ParseFloat(array[3], 64)
	if err != nil {
		return "", err
	}
	duration, err := strconv.ParseFloat(array[4], 64)
	if err != nil {
		return "", err
	}
	result := float64(0)
	switch operation {
	case "+":
		time.Sleep(time.Duration(duration) * time.Minute)
		result = first + second
	case "-":
		time.Sleep(time.Duration(duration) * time.Minute)
		result = first - second
	case "*":
		time.Sleep(time.Duration(duration) * time.Minute)
		result = first * second
	case "/":
		if second == 0 {
			return "", errors.New("делить на ноль нельзя")
		}
		time.Sleep(time.Duration(duration) * time.Minute)
		result = first / second
	}
	return id + " " + strconv.FormatFloat(result, 'f', -1, 64), nil
}
