package user

import (
	"crypto/sha1"
	"errors"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
)

const (
	errLoginOrPassword = "Невалидные логин или пароль"
	saltPassword       = "saltPassword"
)

type User struct {
	conf    *config.Config
	storage *memory.UserStorage
	maxId   uint64
}

func (u *User) Add(login string, password string) error {
	if login == "" || password == "" {
		return errors.New(errLoginOrPassword)
	}
	password = hashPassword(password)
	u.storage.Add(&entity.User{
		Login:    login,
		Password: password,
		Id:       u.maxId,
	})
	u.maxId++
	return nil
}

func hashPassword(password string) string {
	res := sha1.New().Sum([]byte(password + saltPassword))
	return string(res)
}
