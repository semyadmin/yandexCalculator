package auth

import (
	"crypto/sha1"
	"encoding/json"
	"errors"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	jwttoken "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/jwt_token"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
)

const (
	errUserNotExists   = "Пользователь не существует"
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
	return jwttoken.GenerateToken(foundUser.Login)
}

func hashPassword(password string) string {
	res := sha1.New().Sum([]byte(password + saltPassword))
	return string(res)
}
