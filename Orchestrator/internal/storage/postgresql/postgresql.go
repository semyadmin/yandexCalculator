package postgresql

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

// Создаем подключение к базе данных
func DbConnect(Db, DbPort, DbUser, DbPass, DbName string) *sql.DB {
	var db *sql.DB
	var err error
	connect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		Db, DbPort, DbUser, DbPass, DbName)
	db, err = sql.Open("postgres", connect)
	if err != nil {
		slog.Error("Неверные данные для подключения к базе данных", "ОШИБКА:", err)
		panic(err)
	}
	return db
}
