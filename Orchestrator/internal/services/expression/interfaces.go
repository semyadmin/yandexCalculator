package expression

import (
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
)

type Storage interface {
	GeByExpression(string, string) (*entity.Expression, error)
	GetById(uint64, string) (*entity.Expression, error)
	Set(*entity.Expression) error
}

type UserStorage interface {
	SetConfig(string, *entity.Config) error
	GetConfig(string) (*entity.Config, error)
}

type Queue interface {
	Enqueue(queue.Expression)
}

type StorageGetAll interface {
	GetAll(string) []*entity.Expression
}
