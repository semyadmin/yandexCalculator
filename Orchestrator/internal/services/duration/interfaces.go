package duration

import "github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"

type UserStorage interface {
	SetConfig(user string, conf *entity.Config) error
	GetConfig(user string) (*entity.Config, error)
}
