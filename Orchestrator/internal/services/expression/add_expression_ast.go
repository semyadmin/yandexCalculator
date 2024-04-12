package expression

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	jwttoken "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/jwt_token"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/arithmetic"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/validator"
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
		return nil, err
	}
	exp, err := storage.GeByExpression(expression, user)
	if errors.Is(err, memory.ErrExpressionNotExists) {
		exp = entity.NewExpression(expression, "", validator.Validator, user)
		storage.Set(exp)
		_, err := arithmetic.NewASTTree(exp, conf, queue)
		if exp.Err == nil {
			exp.Duration = duration(exp.Expression, conf)
		}
		if err != nil {
			resp := entity.NewResponseExpression(exp.ID, exp.Expression, time.Now(), 0, false, 0, err)
			data, e := json.Marshal(resp)
			if e != nil {
				return nil, e
			}
			return data, nil
		}
	}
	resp := entity.NewResponseExpression(exp.ID, exp.Expression, exp.Start, exp.Duration, exp.IsCalc, exp.Result, exp.Err)
	data, e := json.Marshal(resp)
	if e != nil {
		return nil, e
	}
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
