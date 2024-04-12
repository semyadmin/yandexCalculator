package expression

import (
	"log/slog"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	jwttoken "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/jwt_token"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
)

func GetAllExpressions(storage *memory.Storage, token string) []*entity.Expression {
	name, err := jwttoken.ParseToken(token)
	if err != nil {
		slog.Error("Невозможно расшифровать токен:", "ОШИБКА:", err)
		return nil
	}
	return storage.GetAll(name)
}
