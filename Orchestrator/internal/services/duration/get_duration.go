package duration

import (
	"encoding/json"

	jwttoken "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/jwt_token"
)

func GetDuration(token string, userStorage UserStorage) ([]byte, error) {
	user, err := jwttoken.ParseToken(token)
	if err != nil {
		return nil, err
	}
	currentDuration, err := userStorage.GetConfig(user)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(currentDuration)
	if err != nil {
		return nil, err
	}
	return data, nil
}
