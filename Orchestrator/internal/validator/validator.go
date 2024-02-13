package validator

import (
	"strings"
)

func Validator(str string) (string, bool) {
	str = strings.ReplaceAll(str, " ", "")
	operators := map[byte]uint8{'(': 0, ')': 0, '+': 1, '-': 1, '*': 2, '/': 2, '.': 0}
	numbers := map[byte]struct{}{'0': {}, '1': {}, '2': {}, '3': {}, '4': {}, '5': {}, '6': {}, '7': {}, '8': {}, '9': {}}
	if str[0] == '+' ||
		str[0] == '-' ||
		str[0] == '*' ||
		str[0] == '/' ||
		str[0] == ')' ||
		str[0] == '.' {
		return "", false
	}
	if str[len(str)-1] == '+' ||
		str[len(str)-1] == '-' ||
		str[len(str)-1] == '*' ||
		str[len(str)-1] == '/' ||
		str[len(str)-1] == '(' ||
		str[len(str)-1] == '.' {
		return "", false
	}
	for i := 0; i < len(str)-1; i++ {
		operator := str[i]
		nextOperator := str[i+1]
		if _, ok := operators[operator]; !ok {
			if _, ok := numbers[operator]; !ok {
				return "", false
			}
		}
		if operator == '*' || operator == '/' || operator == '+' || operator == '-' {
			if nextOperator == '*' ||
				nextOperator == '/' ||
				nextOperator == '+' ||
				nextOperator == '-' ||
				nextOperator == ')' {
				return "", false
			}
		}
		if operator == '(' {
			if nextOperator == '*' ||
				nextOperator == '/' ||
				nextOperator == '+' ||
				nextOperator == '-' ||
				nextOperator == ')' {
				return "", false
			}
		}
		if operator == '.' {
			if _, ok := numbers[nextOperator]; !ok {
				return "", false
			}
			if _, ok := numbers[str[i-1]]; !ok {
				return "", false
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
			return "", false
		}
	}
	return str, true
}
