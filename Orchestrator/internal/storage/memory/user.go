package memory

import (
	"sync"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
)

type UserStorage struct {
	conf  *config.Config
	users map[string]*entity.User
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
	u.users[user.Login] = user
}

func (u *UserStorage) GetId(user string) *entity.User {
	u.Lock()
	defer u.Unlock()
	return u.users[user]
}
