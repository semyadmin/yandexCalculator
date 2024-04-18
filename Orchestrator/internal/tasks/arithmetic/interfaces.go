package arithmetic

import (
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
)

type UserStorage interface {
	SetConfig(login string, conf *entity.Config) error
	GetConfig(login string) (*entity.Config, error)
}

type Queue interface {
	Enqueue(queue.Expression)
}
