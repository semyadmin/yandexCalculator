package newexpression

import (
	"encoding/json"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
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

func NewExpression(conf *config.Config, storage *memory.Storage) *NewExpressionAst {
	return &NewExpressionAst{conf: conf, storage: storage}
}

func (n *NewExpressionAst) NewExpression(expression string) ([]byte, error) {
	exp := entity.NewExpression(expression, validator.Validator)
	n.storage.Set(exp)
	ast, err := arithmetic.NewASTTree(exp, n.conf, n.queue)
	if err != nil {
		return nil, err
	}
	return json.Marshal(ast)
}
