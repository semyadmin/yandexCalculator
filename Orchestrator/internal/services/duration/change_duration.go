package duration

import (
	"encoding/json"
	"errors"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	jwttoken "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/jwt_token"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/postgresql/postgresql_config"
)

var errWrongDuration = errors.New("Некорректное время для операций")

type configFromJSON struct {
	Plus     string `json:"plus"`
	Minus    string `json:"minus"`
	Multiply string `json:"multi"`
	Divide   string `json:"divide"`
}

func ChangeDuration(config *config.Config, data []byte, token string, userStorage UserStorage) ([]byte, error) {
	var conf entity.Config
	user, err := jwttoken.ParseToken(token)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	// Сохраняем конфигурацию текущему пользователю
	err = userStorage.SetConfig(user, &conf)
	if err != nil {
		return nil, err
	}
	go config.Db.Config.Add(postgresql_config.Config{
		Plus:     conf.Plus,
		Minus:    conf.Minus,
		Multiply: conf.Multiply,
		Divide:   conf.Divide,
		Login:    user,
	})
	return data, nil
}
