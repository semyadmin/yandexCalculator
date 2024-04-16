package expression

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	jwttoken "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/jwt_token"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/postgresql/postgresql_expression"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/arithmetic"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/web_socket/client"
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

func LoadFromDb(
	conf *config.Config,
	storage *memory.Storage,
	queue *queue.MapQueue,
	userStorage *memory.UserStorage,
) {
	exp := make(chan postgresql_expression.Expression)
	conf.Db.Expression.GetAll(exp)
	for expression := range exp {
		newExp := &entity.Expression{
			ID:                   expression.BaseID,
			Start:                time.Now(),
			Expression:           expression.Expression,
			CalculatedExpression: expression.CurrentResult,
			Result:               expression.Value,
			User:                 expression.User,
		}
		if expression.Err {
			newExp.Err = errors.New("ошибка вычисления")
		}
		_, err := strconv.ParseFloat(expression.CurrentResult, 64)
		if err == nil {
			newExp.IsCalc = true
		}
		conf.MaxID = max(conf.MaxID, newExp.ID)
		err = storage.Set(newExp)
		if err == nil {
			arithmetic.NewASTTree(newExp, conf, queue, userStorage)
		}
		slog.Info("Загружено выражение:", "выражение", newExp)
		resp := entity.NewResponseExpression(
			newExp.ID,
			newExp.Expression,
			newExp.Start,
			newExp.Duration,
			newExp.IsCalc,
			newExp.Result,
			newExp.Err)
		data, e := json.Marshal(resp)
		if e != nil {
			slog.Error("Невозможно сериализовать ответ:", "ОШИБКА:", e)
			continue
		}
		go func() {
			conf.WSmanager.MessageCh <- &client.Message{
				Payload: data,
				Type:    client.ClientExpression,
			}
		}()
	}
	slog.Info("Загрузка выражений из базы данных завершена")
}
