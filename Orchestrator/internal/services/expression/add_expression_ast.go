package expression

import (
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	jwttoken "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/jwt_token"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/postgresql/postgresql_expression"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/arithmetic"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/validator"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/web_socket/client"
)

type NewExpressionAst struct {
	conf    *config.Config
	storage *memory.Storage
	queue   *queue.MapQueue
}

func NewExpression(conf *config.Config,
	storage *memory.Storage,
	queue *queue.MapQueue,
	expression string,
	token string,
) ([]byte, error) {
	user, err := jwttoken.ParseToken(token)
	if err != nil {
		slog.Error("Невозможно распарсить токен:", "ОШИБКА:", err)
		return nil, err
	}
	exp, err := storage.GeByExpression(expression, user)
	var resp entity.ResponseExpression
	if errors.Is(err, memory.ErrExpressionNotExists) {
		exp = entity.NewExpression(expression, "", validator.Validator, user)
		conf.Lock()
		conf.MaxID++
		nextId := conf.MaxID
		conf.Unlock()
		exp.SetId(nextId)
		storage.Set(exp)
		ast, err := arithmetic.NewASTTree(exp, conf, queue)
		if exp.Err == nil {
			exp.Duration = duration(exp.Expression, conf)
		}
		expDb := postgresql_expression.Expression{
			BaseID:        exp.ID,
			Expression:    exp.Expression,
			User:          exp.User,
			Value:         exp.Result,
			CurrentResult: ast.PrintExpression(),
		}
		if exp.Err != nil {
			expDb.Err = true
		}
		conf.Db.Expression.Add(expDb)
		if err != nil {
			resp = entity.NewResponseExpression(exp.ID, exp.Expression, time.Now(), 0, false, 0, err)
		} else {
			resp = entity.NewResponseExpression(exp.ID, exp.Expression, time.Now(), exp.Duration, true, exp.Result, nil)
		}
	}
	data, e := json.Marshal(resp)
	if e != nil {
		return nil, e
	}
	go func() {
		conf.WSmanager.MessageCh <- &client.Message{
			Payload: data,
			Type:    client.ClientExpression,
		}
	}()
	return data, nil
}

func duration(exp string, config *config.Config) int64 {
	res := int64(0)
	for i := 0; i < len(exp); i++ {
		if exp[i] == '+' {
			res += config.Plus
		}
		if exp[i] == '-' {
			if i == 0 ||
				exp[i-1] == '(' ||
				exp[i-1] == '+' ||
				exp[i-1] == '-' ||
				exp[i-1] == '*' ||
				exp[i-1] == '/' {
				continue
			}
			res += config.Minus
		}
		if exp[i] == '*' {
			res += config.Multiply
		}
		if exp[i] == '/' {
			res += config.Divide
		}
	}
	return res
}
