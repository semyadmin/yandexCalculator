package expression

import (
	"log/slog"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	jwttoken "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/jwt_token"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
)

func GetAllExpressions(storage *memory.Storage, token string) []entity.ResponseExpression {
	name, err := jwttoken.ParseToken(token)
	if err != nil {
		slog.Error("Невозможно расшифровать токен:", "ОШИБКА:", err)
		return nil
	}
	allExpressions := storage.GetAll(name)
	result := make([]entity.ResponseExpression, len(allExpressions))
	for i, expression := range allExpressions {
		result[i] = entity.NewResponseExpression(expression.ID, expression.Expression, expression.Start, expression.Duration, expression.IsCalc, expression.Result, expression.Err)
	}
	return result
}
