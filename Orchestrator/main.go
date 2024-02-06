package main

import (
	"errors"
	"strconv"
)

var errValidation = errors.New("Ошибка валидации данных")

func resultPolandNotation(str []string) (float64, error) {
	stack := make([]float64, 0)
	for _, v := range str {
		if v == "+" || v == "-" || v == "*" || v == "/" {
			if len(stack) < 2 {
				return 0, errValidation
			}
			right := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			left := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			switch v {
			case "+":
				stack = append(stack, left+right)
			case "-":
				stack = append(stack, left-right)
			case "*":
				stack = append(stack, left*right)
			case "/":
				stack = append(stack, left/right)
			}
		} else {
			number, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, number)
		}
	}
	if len(stack) != 1 {
		return 0, errValidation
	}
	return stack[0], nil
}

func main() {
	stack := []int{}
	res := stack[len(stack)-1]
	println(res)
}
