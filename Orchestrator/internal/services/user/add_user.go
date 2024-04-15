package user

import (
	"crypto/sha256"
	"encoding/json"
	"errors"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	jwttoken "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/jwt_token"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/postgresql/postgresql_user"
)

const (
	errLoginOrPassword = "Неверные логин или пароль"
)
const saltPassword = "s3cr3t"

func User(userStorage *memory.UserStorage, data []byte, conf *config.Config) (string, error) {
	user := new(entity.User)
	err := json.Unmarshal(data, user)
	if err != nil {
		return "", err
	}
	if user.Login == "" || user.Password == "" {
		return "", errors.New(errLoginOrPassword)
	}
	foundUser, err := userStorage.FindUser(user.Login)
	userConfig := new(memory.User)
	userConfig.User = user
	userConfig.Config = &entity.Config{
		Plus:     conf.Plus,
		Minus:    conf.Minus,
		Multiply: conf.Multiply,
		Divide:   conf.Divide,
	}
	if err != nil {
		user.Password = hashPassword(user.Password)
		userStorage.Add(userConfig)
		go conf.Db.User.Add(postgresql_user.UserStorage{
			Login:    user.Login,
			Password: user.Password,
		})
		return jwttoken.GenerateToken(user.Login)
	}
	if foundUser.User.Password != hashPassword(user.Password) {
		return "", errors.New(errLoginOrPassword)
	}
	return jwttoken.GenerateToken(foundUser.User.Login)
}

func hashPassword(password string) string {
	res := sha256.New()
	res.Write([]byte(password + saltPassword))
	return string(res.Sum(nil))
}
