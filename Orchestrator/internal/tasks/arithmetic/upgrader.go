package arithmetic

// Обновление выражения с добавлением скобок
func Upgrade(exp []byte) []byte {
	operatorIndex := 0
	result := make([]byte, 0)
	left := 0
	for i := 0; i < len(exp); i++ {
		if exp[i] == '+' || exp[i] == '-' {
			if exp[operatorIndex] == '+' || exp[operatorIndex] == '-' {
				array := append([]byte{'('}, exp[left:i]...)
				array = append(array, ')')
				array = append(array, exp[i])
				result = append(result, array...)
				left = i + 1
				operatorIndex = 0
				continue
			}
			operatorIndex = i
			continue
		}
		/* if exp[i] == '*' || exp[i] == '/' {
			if operatorIndex > left {
				result = append(result, exp[left:operatorIndex+1]...)
				left = operatorIndex + 1
			}
			array, index := upgradeMultiDivide(exp, i+1, left)
			result = append(result, array...)
			if index == len(exp) {
				return result
			}
			operatorIndex = 0
			if exp[index] == ')' {
				return result
			}
			if exp[index] == '(' {
				left = index
				i = index - 1
				continue
			}
			i = index
			left = i + 1
			if index == len(exp) {
				break
			}
			result = append(result, exp[index])
			continue
		} */
		if exp[i] == '(' {
			result = append(result, exp[left:i+1]...)
			array := Upgrade(exp[i+1:])

			result = append(result, array...)
			closeBoarder := -1
			for j := i + 1; j < len(exp); j++ {
				if exp[j] == ')' {
					closeBoarder++
				}
				if exp[j] == '(' {
					closeBoarder--
				}
				if closeBoarder == 0 {
					i = j
					break
				}
			}
			result = append(result, exp[i])
			if i == len(exp)-1 {
				return result
			}
			if exp[i+1] == '+' || exp[i+1] == '-' ||
				exp[i+1] == '*' || exp[i+1] == '/' {
				i = i + 1
				result = append(result, exp[i])
			}
			left = i + 1
			operatorIndex = 0
			continue
		}
		if exp[i] == ')' {
			result = append(result, exp[left:i]...)
			return result
		}
	}
	if operatorIndex != 0 {
		array := append([]byte{'('}, exp[left:]...)
		array = append(array, ')')
		result = append(result, array...)
	} else {
		result = append(result, exp[left:]...)
	}

	return result
}

// Добавляем скобки в умножение и деление
func upgradeMultiDivide(exp []byte) []byte {
	result := make([]byte, 0)
	left := 0
	prevOperator := byte(0)
	countBrackets := 0
	for i := 0; i < len(exp); i++ {
		if exp[i] == '*' {
			if prevOperator == 1 {
				result = append(result, exp[i])
				left = i + 1
				prevOperator = exp[i]
				continue
			}
			if exp[i+1] == '(' {
				countBrackets++
				i++
				result = append(result, exp[left:i+1]...)
				array := upgradeMultiDivide(exp[i+1:])
				for countBrackets > 0 {
					i++
					if exp[i] == '(' {
						countBrackets++
					}
					if exp[i] == ')' {
						countBrackets--
					}
				}
				result = append(result, array...)
				result = append(result, ')')
				left = i + 1
				prevOperator = 1
				continue
			}
			if prevOperator == exp[i] || prevOperator == 0 {
				if exp[i+1] == '-' {
					i++
				}
				var j int
				for j = i + 1; j < len(exp); j++ {
					if exp[j] < 48 {
						break
					}
				}
				i = j - 1
				array := append([]byte{'('}, exp[left:i+1]...)
				array = append(array, ')')
				result = append(result, array...)
				left = i + 1
				prevOperator = 1
				continue
			}
			result = append(result, exp[left:i+1]...)
			left = i + 1
			prevOperator = exp[i]
		}
		if exp[i] == '+' || exp[i] == '-' {
			if exp[i] == '-' {
				if i == 0 || exp[i-1] < 48 {
					continue
				}
			}
			result = append(result, exp[left:i+1]...)
			left = i + 1
			prevOperator = 0
		}
		if exp[i] == '/' {
			result = append(result, exp[left:i+1]...)
			left = i + 1
			prevOperator = exp[i]
		}
		if exp[i] == ')' {
			result = append(result, exp[left:i]...)
			return result
		}
	}
	result = append(result, exp[left:]...)
	return result
}
