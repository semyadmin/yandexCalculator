package duration

import (
	"log/slog"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
)

func LoadFromDB(conf *config.Config, userStorage UserStorage) error {
	allConfigs, err := conf.Db.Config.GetAll()
	if err != nil {
		return err
	}
	for _, v := range allConfigs {
		newConf := &entity.Config{
			Plus:     v.Plus,
			Minus:    v.Minus,
			Multiply: v.Multiply,
			Divide:   v.Divide,
		}
		slog.Info("Загружена конфигурация пользователя", "пользователь:", v.Login, "конфигурация:", newConf)
		userStorage.SetConfig(v.Login, newConf)
	}
	return nil
}
