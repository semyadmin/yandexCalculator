package postgresql_user

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

const (
	tableName = "users"
	id        = "id"
	login     = "login"
	password  = "password"
)

type UserStorage struct {
	Id       uint64
	Login    string
	Password string
}

type Data struct {
	conn *sql.DB
}

func New(conn *sql.DB) *Data {
	return &Data{
		conn: conn,
	}
}

func (d *Data) Add(user UserStorage) {
	var err error
	for {
		err = d.conn.Ping()
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	query := fmt.Sprintf(`
		INSERT INTO %s (%s, %s)
		VALUES ($1, $2)`, tableName, login, password)
	sqlPrepare, err := d.conn.Prepare(query)
	defer sqlPrepare.Close()
	if err != nil {
		return
	}
	_, err = sqlPrepare.Query(
		user.Login,
		user.Password,
	)
	if err != nil {
		slog.Info("Не удалось добавить пользователя в базу данных", "пользователь:", user)
		return
	}
}

func (d *Data) GetAll() []UserStorage {
	var err error
	for {
		err = d.conn.Ping()
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	result := make([]UserStorage, 0)
	rows, err := d.conn.Query(fmt.Sprintf("SELECT %s, %s, %s FROM %s", id, login, password, tableName))
	if err != nil {
		return result
	}
	for rows.Next() {
		user := UserStorage{}
		err = rows.Scan(&user.Id, &user.Login, &user.Password)
		if err != nil {
			return result
		}
		result = append(result, user)
	}
	return result
}
