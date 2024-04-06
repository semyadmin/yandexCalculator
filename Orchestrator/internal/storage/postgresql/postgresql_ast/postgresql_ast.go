package postgresql_ast

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/postgresql"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/arithmetic"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
)

const (
	tableName     = "expressions"
	baseId        = "baseid"
	expression    = "expression"
	value         = "value"
	errColumn     = "err"
	currentResult = "currentresult"
)

// Add — добавляет запись в базу данных
func Add(model *arithmetic.ASTTree, conf *config.Config) {
	go func() {
		/* db := postgresql.DbConnect(conf)
		defer db.Close()
		ok, err := GetByExpression(model.Expression, conf)
		if err != nil {
			return
		}
		if ok {
			return
		}
		query := fmt.Sprintf(`
			INSERT INTO %s (%s, %s, %s, %s, %s)
			VALUES ($1, $2, $3, $4, $5)`,
			tableName, baseId, expression, value, errColumn, currentResult)
		sqlPrepare, err := db.Prepare(query)
		defer sqlPrepare.Close()
		if err != nil {
			return
		}
		model.Lock()
		astBaseID := model.ID
		astExp := model.Expression
		astValue := model.Value
		astErr := model.Err
		astCurrentRes := arithmetic.PrintExpression(model)
		model.Unlock()

		currentErr := false
		if astErr != nil {
			currentErr = true
		}
		_, err = sqlPrepare.Query(
			astBaseID,
			astExp,
			astValue,
			currentErr,
			astCurrentRes,
		)
		if err != nil {
			slog.Info("Не удалось добавить запись в базу данных", "ошибка:", err)
			return
		}
		slog.Info("Добавление записи в базу данных", "выражение:", model.Expression) */
	}()
}

// Проверяем наличие выражения по строке выражения
func GetByExpression(exp string, conf *config.Config) (bool, error) {
	db := postgresql.DbConnect(conf)
	defer db.Close()
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1 LIMIT 1", expression, tableName, expression)
	sqlPrepare, err := db.Prepare(query)
	if err != nil {
		return false, err
	}
	row := sqlPrepare.QueryRow(exp)
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
func GetAll(conf *config.Config, q *queue.MapQueue, m *memory.Storage) {
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
func Update(conf *config.Config, q *queue.MapQueue, m *memory.Storage) {
	go func() {
		/* for {
			q.Lock()
			items := []string{}
			for item := range q.Update {
				items = append(items, item)
				delete(q.Update, item)
			}
			q.Unlock()
			if len(items) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}
			db := postgresql.DbConnect(conf)

			for _, item := range items {
				data, err := m.GeByExpression(item)
				if err != nil {
					continue
				}
				ast := data.Expression
				if ast == nil {
					continue
				}
				ast.Lock()
				slog.Info("Получено выражение для обновления в базе", "выражение:", ast.Expression)
				astBaseID := ast.ID
				astExp := ast.Expression
				astValue := ast.Value
				astErr := ast.Err
				ast.Unlock()
				astCurrentRes := arithmetic.PrintExpression(ast)
				query := fmt.Sprintf(`
				UPDATE %s
				SET %s = $1, %s = $2, %s = $3, %s = $4
				WHERE %s = $5`,
					tableName, baseId, value, errColumn, currentResult, expression)
				sqlPrepare, err := db.Prepare(query)
				if err != nil {
					continue
				}
				currentErr := false

				if astErr != nil {
					currentErr = true
				}
				_, err = sqlPrepare.Query(
					astBaseID,
					astValue,
					currentErr,
					astCurrentRes,
					astExp,
				)
				if err != nil {
					continue
				}
				sqlPrepare.Close()
				slog.Info("Обновление записи в базе данных", "выражение:", astExp)
			}
			db.Close()
			time.Sleep(1 * time.Second)
		} */
	}()
}
