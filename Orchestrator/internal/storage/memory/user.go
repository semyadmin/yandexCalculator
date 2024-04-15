package memory

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
)

type User struct {
	User   *entity.User
	Config *entity.Config
}

const errUserNotExists = "Пользователь не существует"

type UserStorage struct {
	conf  *config.Config
	users map[string]*User
	MaxId uint64
	sync.Mutex
}

func NewUserStorage(conf *config.Config) *UserStorage {
	return &UserStorage{
		conf:  conf,
		users: make(map[string]*User),
	}
}

func (u *UserStorage) Add(user *User) {
	u.Lock()
	defer u.Unlock()
	if user.User.Id == 0 {
		u.MaxId++
		user.User.Id = u.MaxId
	}
	u.MaxId = max(u.MaxId, user.User.Id)
	u.users[user.User.Login] = user
	slog.Info("Добавлен новый пользователь", "пользователь:", user)
}

func (u *UserStorage) FindUser(login string) (*User, error) {
	u.Lock()
	defer u.Unlock()
	user := u.users[login]
	if user == nil {
		return nil, errors.New(errUserNotExists)
	}
	return user, nil
}

func (u *UserStorage) GetId(user string) *User {
	u.Lock()
	defer u.Unlock()
	return u.users[user]
}

func (u *UserStorage) SetConfig(user string, config *entity.Config) error {
	u.Lock()
	defer u.Unlock()
	if _, ok := u.users[user]; !ok {
		return errors.New(errUserNotExists)
	}
	u.users[user].Config = config
	return nil
}

func (u *UserStorage) GetConfig(user string) (*entity.Config, error) {
	u.Lock()
	defer u.Unlock()
	if _, ok := u.users[user]; !ok {
		return nil, errors.New(errUserNotExists)
	}
	return u.users[user].Config, nil
}
