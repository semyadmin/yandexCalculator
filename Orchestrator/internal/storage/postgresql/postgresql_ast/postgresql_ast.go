package postgresql_ast

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

const (
	tableName     = "expressions"
	baseId        = "baseid"
	expression    = "expression"
	value         = "value"
	errColumn     = "err"
	currentResult = "currentresult"
	user          = "login"
)

type Expression struct {
	BaseID        uint64
	Expression    string
	Value         float64
	Err           bool
	User          string
	CurrentResult string
}

type Data struct {
	conn *sql.DB
	sync.Mutex
	updateExp map[string]Expression
}

func NewData(conn *sql.DB) *Data {
	data := &Data{
		conn:      conn,
		updateExp: make(map[string]Expression),
	}
	data.update()
	return data
}

// Add — добавляет запись в базу данных
func (d *Data) Add(exp Expression) {
	go func() {
		var err error
		for {
			err = d.conn.Ping()
			if err == nil {
				break
			}
			time.Sleep(5 * time.Second)
		}
		isExp, err := d.isExpression(exp)
		if err != nil {
			slog.Info("Не удалось проверить наличие выражения в базе данных", "ошибка:", err)
			return
		}
		if isExp {
			slog.Info("Выражение уже существует в базе данных", "выражение:", exp)
			return
		}
		query := fmt.Sprintf(`
			INSERT INTO %s (%s, %s, %s, %s, %s, %s)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			tableName, baseId, expression, user, value, errColumn, currentResult)
		sqlPrepare, err := d.conn.Prepare(query)
		defer sqlPrepare.Close()
		if err != nil {
			return
		}

		_, err = sqlPrepare.Query(
			exp.BaseID,
			exp.Expression,
			exp.User,
			exp.Value,
			exp.Err,
			exp.CurrentResult,
		)
		if err != nil {
			slog.Info("Не удалось добавить запись в базу данных", "ошибка:", err, "выражение:", exp)
			return
		}
		slog.Info("Добавлено выражение в базу данных", "выражение:", exp)
	}()
}

func (d *Data) Update(exp Expression) {
	d.Lock()
	d.updateExp[exp.Expression] = exp
	d.Unlock()
}

// Проверяем наличие выражения по строке выражения
func (d *Data) isExpression(exp Expression) (bool, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1 AND %s = $2 LIMIT 1", expression, tableName, expression, user)
	sqlPrepare, err := d.conn.Prepare(query)
	if err != nil {
		return false, err
	}
	row := sqlPrepare.QueryRow(exp.Expression, exp.User)
	result := ""
	err = row.Scan(&result)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		slog.Info("Не удалось отсканировать запись из базы данных", "ОШИБКА:", err)
		return false, err
	}
	return true, nil
}

// GetAll — получает все записи из базы данных
func GetAll() {
	go func() {
		/* db := postgresql.DbConnect(conf)
		defer db.Close()
		query := fmt.Sprintf("SELECT %s, %s, %s, %s, %s FROM %s", baseId, expression, value, errColumn, currentResult, tableName)
		sql, err := db.Prepare(query)
		if err != nil {
			return
		}
		rows, err := sql.Query()
		if err != nil {
			return
		}
		defer rows.Close()
		currentId := uint64(0)
		currentExp := ""
		currentValue := ""
		currentErr := false
		currentResult := ""
		for rows.Next() {
			err := rows.Scan(&currentId, &currentExp, &currentValue, &currentErr, &currentResult)
			if err != nil {
				return
			}
			entity, _ := arithmetic.NewASTTreeDB(currentId, currentExp, currentValue, currentErr, currentResult, conf, q)
			if entity == nil {
				continue
			}
			entity.Lock()
			entity.ID = currentId
			entity.Expression = currentExp
			entity.Value = currentValue
			if currentErr {
				entity.Err = errors.New("Здесь была ошибка")
			}
			entity.Unlock()
			m.SetFromDb(entity, "saved")
		} */
	}()
}

// Обновляем выражение по его строке выражения
func (d *Data) update() {
	go func() {
		var err error
		items := make([]Expression, 0)
		for {
			err = d.conn.Ping()
			if err != nil {
				time.Sleep(5 * time.Second)
				continue
			}
			items = items[:0]
			d.Lock()
			for _, item := range d.updateExp {
				items = append(items, item)
				delete(d.updateExp, item.Expression)
			}
			d.Unlock()
			for _, item := range items {
				query := fmt.Sprintf(`
				UPDATE %s
				SET %s = $1, %s = $2, %s = $3, %s = $4
				WHERE %s = $5 AND %s = $6;`,
					tableName, baseId, value, errColumn, currentResult, expression, user)
				sqlPrepare, err := d.conn.Prepare(query)
				if err != nil {
					d.updateExp[item.Expression] = item
					continue
				}
				_, err = sqlPrepare.Query(
					item.BaseID,
					item.Value,
					item.Err,
					item.CurrentResult,
					item.Expression,
					item.User,
				)
				if err != nil {
					d.updateExp[item.Expression] = item
					continue
				}
				sqlPrepare.Close()
				slog.Info("Обновление записи в базе данных", "выражение:", item)
			}
			time.Sleep(1 * time.Second)
		}
	}()
}
