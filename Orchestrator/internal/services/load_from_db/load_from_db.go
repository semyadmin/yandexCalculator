package loadfromdb

import (
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/expression"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/user"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
)

func LoadFromDB(conf *config.Config, store *memory.Storage, userStorage *memory.UserStorage, queue *queue.MapQueue) {
	go func() {
		user.GetAllUsers(conf, userStorage)
		expression.LoadFromDb(conf, store, queue, userStorage)
	}()
}
