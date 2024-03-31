package validator

import (
	"strings"
)

// Валидируем выражение
func Validator(str string) (string, bool) {
	str = strings.ReplaceAll(str, " ", "")
	operators := map[byte]uint8{'(': 0, ')': 0, '+': 1, '-': 1, '*': 2, '/': 2, '.': 0}
	numbers := map[byte]struct{}{'0': {}, '1': {}, '2': {}, '3': {}, '4': {}, '5': {}, '6': {}, '7': {}, '8': {}, '9': {}}
	isNegative := false
	countBrackets := 0
	if str[0] == '-' {
		ok := false
		for i := 1; i < len(str); i++ {
			if _, ok = operators[str[i]]; ok {
				break
			}
		}
		if !ok {
			return "", false
		}
	}
	if str[0] == '+' ||
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
		currentByte := str[i]
		if currentByte != '-' {
			isNegative = false
		}
		nextByte := str[i+1]
		if _, ok := operators[currentByte]; !ok {
			if _, ok := numbers[currentByte]; !ok {
				return "", false
			}
		}
		if currentByte == '*' || currentByte == '/' || currentByte == '+' || currentByte == '-' {
			if nextByte == '*' ||
				nextByte == '/' ||
				nextByte == '+' ||
				nextByte == ')' ||
				isNegative {
				return "", false
			}
			if currentByte == '-' {
				if i == 0 {
					isNegative = true
					continue
				}
				// Проверяем предыдущее значение
				if str[i-1] == '(' || str[i-1] == '-' || str[i-1] == '+' || str[i-1] == '*' || str[i-1] == '/' {
					isNegative = true
					continue
				}
			}
		}
		if currentByte == '(' {
			if i > 0 {
				if _, ok := operators[str[i-1]]; !ok {
					return "", false
				}
			}
			countBrackets++
			if nextByte == '*' ||
				nextByte == '/' ||
				nextByte == '+' ||
				nextByte == ')' {
				return "", false
			}
		}
		if currentByte == '.' {
			if _, ok := numbers[nextByte]; !ok {
				return "", false
			}
		}

		if currentByte == ')' {
			countBrackets--
			if countBrackets < 0 {
				return "", false
			}
			if nextByte != '*' &&
				nextByte != '/' &&
				nextByte != '+' &&
				nextByte != '-' {
				return "", false
			}
		}
	}
	if str[len(str)-1] == ')' {
		countBrackets--
	}
	if countBrackets != 0 {
		return "", false
	}
	return str, true
}
