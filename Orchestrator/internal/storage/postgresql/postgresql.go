package postgresql

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	_ "github.com/lib/pq"
)

func DbConnect(conf *config.Config) *sql.DB {
	var db *sql.DB
	var err error
	for {
		connect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			conf.Db, conf.DbPort, conf.DbUser, conf.DbPass, conf.DbName)
		db, err = sql.Open("postgres", connect)
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}
		err = db.Ping()
		if err == nil {
			break
		}
		time.Sleep(time.Second * 5)
	}
	slog.Info("Соединение с базой данных установлено")
	return db
}
