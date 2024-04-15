package duration

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	jwttoken "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/jwt_token"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
)

var errWrongDuration = errors.New("Некорректное время для операций")

type configFromJSON struct {
	Plus     string `json:"plus"`
	Minus    string `json:"minus"`
	Multiply string `json:"multi"`
	Divide   string `json:"divide"`
}

func ChangeDuration(data []byte, token string, userStorage *memory.UserStorage) ([]byte, error) {
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
	data, err = json.Marshal(conf)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func newDuration(conf configFromJSON, c *entity.Config) error {
	num, err := parseStringToInt(conf.Plus)
	if err != nil || num < 0 {
		return errWrongDuration
	}
	c.Plus = num
	num, err = parseStringToInt(conf.Minus)
	if err != nil || num < 0 {
		return errWrongDuration
	}
	c.Minus = num
	num, err = parseStringToInt(conf.Multiply)
	if err != nil || num < 0 {
		return errWrongDuration
	}
	c.Multiply = num
	num, err = parseStringToInt(conf.Divide)
	if err != nil || num < 0 {
		return errWrongDuration
	}
	c.Divide = num

	return nil
}

func parseStringToInt(str string) (int64, error) {
	num64, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		return 0, err
	}
	return num64, nil
}
