package user

import (
	"crypto/sha1"
	"errors"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	jwttoken "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/jwt_token"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
)

const (
	errLoginOrPassword = "Невалидные логин или пароль"
	errUserNotExists   = "Невалидные логин или пароль"
	saltPassword       = "saltPassword"
)

type User struct {
	conf    *config.Config
	storage *memory.UserStorage
	maxId   uint64
}

func (u *User) Add(login string, password string) (string, error) {
	if login == "" || password == "" {
		return "", errors.New(errLoginOrPassword)
	}
	foundUser, err := u.FindUser(login)
	if errors.Is(err, errors.New(errUserNotExists)) {
		u.maxId++
		password = hashPassword(password)
		foundUser = &entity.User{
			Login:    login,
			Password: password,
			Id:       u.maxId,
		}
		u.storage.Add(foundUser)
	} else if err != nil {
		return "", err
	}
	password = hashPassword(password)
	if password != foundUser.Password {
		return "", errors.New(errLoginOrPassword)
	}
	token, err := jwttoken.GenerateToken(login)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *User) FindUser(login string) (*entity.User, error) {
	user := u.storage.GetId(login)
	if user == nil {
		return nil, errors.New(errUserNotExists)
	}
	return user, nil
}

func hashPassword(password string) string {
	res := sha1.New().Sum([]byte(password + saltPassword))
	return string(res)
}
