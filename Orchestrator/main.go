package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errValidation error = errors.New("Ошибка валидации данных")

func polandNotation(str string) ([]string, error) {
	str = strings.ReplaceAll(str, " ", "")
	result := make([]string, 0, len(str))
	stack := make([]byte, 0, 10)
	operators := map[byte]uint8{'(': 0, ')': 0, '+': 1, '-': 1, '*': 2, '/': 2}
	numbers := map[byte]struct{}{'0': {}, '1': {}, '2': {}, '3': {}, '4': {}, '5': {}, '6': {}, '7': {}, '8': {}, '9': {}}
	index := 0
	indexNumber := 0
	opened := 0
	closed := 0
	res := checkCorrectness(str)
	fmt.Println(res)
	if !res {
		return nil, errValidation
	}
	for index < len(str) {
		if _, ok := numbers[str[index]]; ok {
			indexNumber = index
			for indexNumber < len(str) {
				if _, ok := numbers[str[indexNumber]]; !ok {
					break
				}
				indexNumber++
			}
			result = append(result, str[index:indexNumber])
			index = indexNumber
			continue
		}
		operator := str[index]
		if operator == '(' {
			opened++
			stack = append(stack, operator)
			index++
			continue
		}
		if operator == ')' {
			closed++
			if closed > opened {
				return nil, errValidation
			}
			for i := len(stack) - 1; i >= 0; i-- {
				if stack[i] == '(' {
					stack = stack[:i]
					break
				}
				result = append(result, string(stack[i]))
				stack = stack[:i]
			}
			index++
			continue
		}
		for i := len(stack) - 1; i >= 0; i-- {
			if operators[operator] >= operators[stack[i]] {
				break
			}
			result = append(result, string(stack[i]))
			stack = stack[:i]
		}
		stack = append(stack, operator)
		index++
	}
	if opened != closed {
		return nil, errValidation
	}
	for i := len(stack) - 1; i >= 0; i-- {
		result = append(result, string(stack[i]))
		stack = stack[:i]
	}
	return result, nil
}

func checkCorrectness(str string) bool {
	operators := map[byte]uint8{'(': 0, ')': 0, '+': 1, '-': 1, '*': 2, '/': 2}
	numbers := map[byte]struct{}{'0': {}, '1': {}, '2': {}, '3': {}, '4': {}, '5': {}, '6': {}, '7': {}, '8': {}, '9': {}}
	if str[0] == '+' || str[0] == '-' || str[0] == '*' || str[0] == '/' || str[0] == ')' {
		return false
	}
	if str[len(str)-1] == '+' || str[len(str)-1] == '-' || str[len(str)-1] == '*' || str[len(str)-1] == '/' || str[len(str)-1] == '(' {
		return false
	}
	fmt.Println(str, len(str))
	for i := 0; i < len(str)-1; i++ {
		operator := str[i]
		nextOperator := str[i+1]
		if _, ok := operators[operator]; !ok {
			if _, ok := numbers[operator]; !ok {
				return false
			}
		}
		if operator == '*' || operator == '/' || operator == '+' || operator == '-' {
			if nextOperator == '*' ||
				nextOperator == '/' ||
				nextOperator == '+' ||
				nextOperator == '-' ||
				nextOperator == ')' {
				return false
			}
		}
		if operator == '(' {
			if nextOperator == '*' ||
				nextOperator == '/' ||
				nextOperator == '+' ||
				nextOperator == '-' ||
				nextOperator == ')' {
				return false
			}
		}

		if operator == ')' {
			if nextOperator == '*' ||
				nextOperator == '/' ||
				nextOperator == '+' ||
				nextOperator == '-' ||
				nextOperator == ')' {
				continue
			}
			return false
		}
	}
	return true
}

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
	fmt.Println(stack[0])
	return stack[0], nil
}

func main() {
	res, err := polandNotation("00001+(155+2)*44")
	if err != nil {
		println(err.Error())
		return
	}
	for _, v := range res {
		println(v)
	}
	sum, err := resultPolandNotation(res)
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Println(sum)
}
