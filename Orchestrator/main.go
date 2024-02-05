package main

import (
	"errors"
	"strings"
)

var errValidation error = errors.New("Ошибка валидации данных")

func polangNotation(str string) ([]string, error) {
	str = strings.ReplaceAll(str, " ", "")
	println(str)
	result := make([]string, 0, len(str))
	stack := make([]byte, 0, 10)
	operators := map[byte]uint8{'(': 0, ')': 0, '+': 1, '-': 1, '*': 2, '/': 2}
	numbers := map[byte]struct{}{'0': {}, '1': {}, '2': {}, '3': {}, '4': {}, '5': {}, '6': {}, '7': {}, '8': {}, '9': {}}
	index := 0
	indexNumber := 0
	closed := false
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
		if _, ok := operators[operator]; !ok {
			return nil, errValidation
		}
		if index == 0 && operator != '(' {
			return nil, errValidation
		}
		if index == len(str)-1 && operator != ')' {
			return nil, errValidation
		}
		if index != len(str)-1 {
			if _, ok := operators[str[index+1]]; ok {
				return nil, errValidation
			}
		}
		if operator == '(' {
			if closed {
				return nil, errValidation
			}
			closed = true
			stack = append(stack, operator)
		}
		if operator == ')' {
			if !closed {
				return nil, errValidation
			}
			closed = false
			for i := len(stack) - 1; i >= 0; i-- {
				if stack[i] == '(' {
					stack = stack[:i]
					index++
					continue
				}
				result = append(result, string(stack[i]))
				stack = stack[:i]
			}
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
	for i := len(stack) - 1; i >= 0; i-- {
		println(111111111111111)
		result = append(result, string(stack[i]))
		stack = stack[:i]
	}
	return result, nil
}

func main() {
	res, err := polangNotation("(   1 +     2 * 3 )   ")
	if err != nil {
		println(err.Error())
		return
	}
	for _, v := range res {
		println(v)
	}
}
