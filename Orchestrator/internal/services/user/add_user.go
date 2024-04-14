package user

import (
	"crypto/sha256"
	"encoding/json"
	"errors"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	jwttoken "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/jwt_token"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
)

const (
	errLoginOrPassword = "Неверные логин или пароль"
)
const saltPassword = "s3cr3t"

func User(userStorage *memory.UserStorage, data []byte) (string, error) {
	user := new(entity.User)
	err := json.Unmarshal(data, user)
	if err != nil {
		return "", err
	}
	if user.Login == "" || user.Password == "" {
		return "", errors.New(errLoginOrPassword)
	}
	foundUser, err := userStorage.FindUser(user.Login)
	if err != nil {
		user.Password = hashPassword(user.Password)
		userStorage.Add(user)
		return jwttoken.GenerateToken(user.Login)
	}
	if foundUser.Password != hashPassword(user.Password) {
		return "", errors.New(errLoginOrPassword)
	}
	return jwttoken.GenerateToken(foundUser.Login)
}

func hashPassword(password string) string {
	res := sha256.New()
	res.Write([]byte(password + saltPassword))
	return string(res.Sum(nil))
}
