package postgresql_config

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

const (
	tableName = "configs"
	plus      = "plus"
	minus     = "minus"
	multiply  = "multiply"
	divide    = "divide"
	login     = "login"
)

type Config struct {
	Plus     int64
	Minus    int64
	Multiply int64
	Divide   int64
	Login    string
}

type Data struct {
	conn *sql.DB
}

func New(conn *sql.DB) *Data {
	return &Data{
		conn: conn,
	}
}

func (d *Data) Add(conf Config) error {
	var err error
	for {
		err = d.conn.Ping()
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	query := fmt.Sprintf(`
		INSERT INTO %s (%s, %s, %s, %s, %s)
		VALUES ($1, $2, $3, $4, $5)`, tableName, plus, minus, multiply, divide, login)
	sqlPrepare, err := d.conn.Prepare(query)
	defer sqlPrepare.Close()
	if err != nil {
		return err
	}
	_, err = sqlPrepare.Query(
		conf.Plus,
		conf.Minus,
		conf.Multiply,
		conf.Divide,
		conf.Login,
	)
	if err != nil {
		slog.Info("Не удалось сохранить длительность выполнения операции в базу данных", "ошибка:", err, "конфиг:", conf)
		return err
	}
	return nil
}

func (d *Data) GetAll() ([]Config, error) {
	var err error
	for {
		err = d.conn.Ping()
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	query := fmt.Sprintf("SELECT %s, %s, %s, %s, %s FROM %s", plus, minus, multiply, divide, login, tableName)
	rows, err := d.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var configs []Config
	for rows.Next() {
		var conf Config
		if err := rows.Scan(&conf.Plus, &conf.Minus, &conf.Multiply, &conf.Divide, &conf.Login); err != nil {
			slog.Error("Не удалось получить конфиг из базы данных", "ошибка:", err)
			return nil, err
		}
		configs = append(configs, conf)
	}
	return configs, nil
}
