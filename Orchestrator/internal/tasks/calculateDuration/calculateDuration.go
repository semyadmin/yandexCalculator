package calculateDuration

import "github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"

func Calc(expression []string, config *config.ConfigExpression) uint64 {
	res := uint64(0)
	for i := 0; i < len(expression); i++ {
		if expression[i] == "+" {
			res += uint64(config.Plus)
		}
		if expression[i] == "-" {
			res += uint64(config.Minus)
		}
		if expression[i] == "*" {
			res += uint64(config.Multiply)
		}
		if expression[i] == "*" {
			res += uint64(config.Divide)
		}
	}
	return res
}
