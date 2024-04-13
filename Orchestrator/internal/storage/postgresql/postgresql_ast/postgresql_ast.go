package postgresql_ast

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
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
	user          = "login"
)

// Add — добавляет запись в базу данных
func Add(exp *entity.Expression, ast *arithmetic.ASTTree, conn *sql.DB) {
	go func() {
		var err error
		for {
			err = conn.Ping()
			if err == nil {
				break
			}
			time.Sleep(5 * time.Second)
		}
		query := fmt.Sprintf(`
			INSERT INTO %s (%s, %s, %s, %s, %s, %s)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			tableName, baseId, expression, user, value, errColumn, currentResult)
		sqlPrepare, err := conn.Prepare(query)
		defer sqlPrepare.Close()
		if err != nil {
			return
		}
		astBaseID := exp.ID
		astExp := exp.Expression
		astUser := exp.User
		astValue := exp.Result
		astErr := exp.Err
		astCurrentRes := ast.PrintExpression()

		currentErr := false
		if astErr != nil {
			currentErr = true
		}
		_, err = sqlPrepare.Query(
			astBaseID,
			astExp,
			astUser,
			astValue,
			currentErr,
			astCurrentRes,
		)
		if err != nil {
			slog.Info("Не удалось добавить запись в базу данных", "ошибка:", err)
			return
		}
		slog.Info("Добавлено выражение в базу данных", "выражение:", exp.Expression)
	}()
}

// Проверяем наличие выражения по строке выражения
func GetByExpression(exp string, conf *config.Config) (bool, error) {
	/* db := postgresql.DbConnect(conf)
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
	*/
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
