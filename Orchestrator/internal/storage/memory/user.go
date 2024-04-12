package memory

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
)

const errUserNotExists = "Пользователь не существует"

type UserStorage struct {
	conf  *config.Config
	users map[string]*entity.User
	MaxId uint64
	sync.Mutex
}

func NewUserStorage(conf *config.Config) *UserStorage {
	return &UserStorage{
		conf:  conf,
		users: make(map[string]*entity.User),
	}
}

func (u *UserStorage) Add(user *entity.User) {
	u.Lock()
	defer u.Unlock()
	u.MaxId++
	user.Id = u.MaxId
	u.users[user.Login] = user
	slog.Info("Добавлен новый пользователь", "пользователь:", user)
}

func (u *UserStorage) FindUser(login string) (*entity.User, error) {
	u.Lock()
	defer u.Unlock()
	user := u.users[login]
	if user == nil {
		return nil, errors.New(errUserNotExists)
	}
	return user, nil
}

func (u *UserStorage) GetId(user string) *entity.User {
	u.Lock()
	defer u.Unlock()
	return u.users[user]
}
