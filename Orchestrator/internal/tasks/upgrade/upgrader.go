package upgrade

// Обновление выражения с добавлением скобок
func Upgrade(exp []byte) []byte {
	exp = upgradePlusMinus(exp)
	exp = upgradeMultiDivide(exp)
	exp = updateDoubleMinus(exp)
	return exp
}

// Добавляем скобки в умножение и деление
func upgradeMultiDivide(exp []byte) []byte {
	result := make([]byte, 0)
	left := 0
	prevOperator := byte(0)
	countBrackets := 0

	for i := left; i < len(exp); i++ {
		if exp[i] == '(' {
			countBrackets++
			result = append(result, exp[left:i]...)
			result = append(result, '(')
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
			i++
			result = append(result, array...)
			result = append(result, ')')
			left = i
			if i < len(exp) {
				result = append(result, exp[i])
				left = i + 1
			}
			prevOperator = 0
			continue
		}
		if exp[i] == '*' {
			if prevOperator == 1 {
				result = append(result, exp[i])
				left = i + 1
				prevOperator = 0
				continue
			}
			if prevOperator == 0 {
				prevOperator = exp[i]
				continue
			}
			if prevOperator == exp[i] {
				result = append(result, '(')
				result = append(result, exp[left:i]...)
				result = append(result, ')', exp[i])
				left = i + 1
				prevOperator = 0
				continue
			}
			result = append(result, exp[left:i+1]...)
			left = i + 1
			if prevOperator == '/' {
				prevOperator = 0
			} else {
				prevOperator = exp[i]
			}
		}
		if exp[i] == '+' || exp[i] == '-' {
			if exp[i] == '-' {
				if i == 0 || exp[i-1] < 48 {
					continue
				}
			}
			if prevOperator == '*' {
				result = append(result, '(')
				result = append(result, exp[left:i]...)
				result = append(result, ')', exp[i])
				left = i + 1
				prevOperator = 0
				continue
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
			if prevOperator == '*' {
				result = append(result, '(')
				result = append(result, exp[left:i]...)
				result = append(result, ')')
				return result
			}
			result = append(result, exp[left:i]...)
			return result
		}
	}
	if prevOperator == '*' {
		result = append(result, '(')
		result = append(result, exp[left:]...)
		result = append(result, ')')
		return result
	}
	result = append(result, exp[left:]...)
	return result
}

// Добавляем скобки в умножение и деление
func upgradePlusMinus(exp []byte) []byte {
	result := make([]byte, 0)
	left := 0
	prevOperator := byte(0)
	countBrackets := 0
	for i := left; i < len(exp); i++ {
		if exp[i] == '(' {
			countBrackets++
			result = append(result, exp[left:i]...)
			result = append(result, '(')
			array := upgradePlusMinus(exp[i+1:])
			for countBrackets > 0 {
				i++
				if exp[i] == '(' {
					countBrackets++
				}
				if exp[i] == ')' {
					countBrackets--
				}
			}
			i++
			result = append(result, array...)
			result = append(result, ')')
			left = i
			if i < len(exp) {
				result = append(result, exp[i])
				left = i + 1
			}
			prevOperator = 0
			continue
		}
		if exp[i] == '+' {
			if prevOperator == '*' || prevOperator == '/' {
				result = append(result, exp[left:i+1]...)
				left = i + 1
				prevOperator = 0
				continue
			}
			if prevOperator == 0 {
				prevOperator = exp[i]
				continue
			}
			result = append(result, '(')
			result = append(result, exp[left:i]...)
			result = append(result, ')', exp[i])
			left = i + 1
			prevOperator = 0
		}
		if exp[i] == '*' || exp[i] == '/' {
			result = append(result, exp[left:i+1]...)
			left = i + 1
			prevOperator = exp[i]
		}
		if exp[i] == '-' {
			if i == 0 || exp[i-1] < 48 {
				continue
			}
			if prevOperator == '*' || prevOperator == '/' {
				result = append(result, exp[left:i]...)
				result = append(result, '+')
				left = i
				prevOperator = 0
				continue
			}
			if prevOperator == 0 {
				prevOperator = exp[i]
				continue
			}
			result = append(result, '(')
			result = append(result, exp[left:i]...)
			result = append(result, ')', '+')
			left = i
			prevOperator = 0
			continue
		}
		if exp[i] == ')' {
			if prevOperator == '+' || prevOperator == '-' {
				result = append(result, '(')
				result = append(result, exp[left:i]...)
				result = append(result, ')')
				return result
			}
			result = append(result, exp[left:i]...)
			return result
		}
	}
	if prevOperator == '+' || prevOperator == '-' {
		result = append(result, '(')
		result = append(result, exp[left:]...)
		result = append(result, ')')
		return result
	}
	result = append(result, exp[left:]...)
	return result
}

func updateDoubleMinus(exp []byte) []byte {
	result := make([]byte, 0)
	left := 0
	for i := 0; i < len(exp)-1; i++ {
		if exp[i] == '-' && exp[i+1] == '-' {
			result = append(result, exp[left:i]...)
			result = append(result, '+')
			left = i + 2
		}
	}
	result = append(result, exp[left:]...)
	return result
}
