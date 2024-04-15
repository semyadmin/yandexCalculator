package user

import (
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
)

func GetAllUsers(conf *config.Config, userStorage *memory.UserStorage) {
	users := conf.Db.User.GetAll()
	for _, user := range users {
		memUser := memory.User{}
		memUser.User.Id = user.Id
		memUser.User.Login = user.Login
		memUser.User.Password = user.Password
		userStorage.Add(&memUser)
	}
}
