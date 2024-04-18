package expression

import (
	"encoding/json"
	"log/slog"
	"strconv"

	jwttoken "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/jwt_token"
	responseexpression "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/response_expression"
)

func GetById(storage Storage, number string, token string) ([]byte, error) {
	user, err := jwttoken.ParseToken(token)
	if err != nil {
		return nil, err
	}
	id, err := strconv.ParseUint(number, 10, 64)
	if err != nil {
		slog.Error("Невозможно распарсить ID:", "ОШИБКА:", err)
		return nil, err
	}
	exp, err := storage.GetById(id, user)
	if err != nil {
		slog.Error("Невозможно получить данные по ID:", "ОШИБКА:", err, "ID:", id)
		return nil, err
	}
	resp := responseexpression.NewResponseExpression(exp.ID, exp.Expression, exp.Start, exp.Duration, exp.IsCalc, exp.Result, exp.Err)
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Невозможно сериализовать данные:", "ОШИБКА:", err)
		return nil, err
	}
	return data, nil
}
