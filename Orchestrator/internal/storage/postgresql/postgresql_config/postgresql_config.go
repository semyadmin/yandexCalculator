package postgresql_config

import (
	"log/slog"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
)

const (
	tableName = "configs"
	plus      = "plus"
	minus     = "minus"
	multiply  = "multiply"
	divide    = "divide"
	maxid     = "maxid"
)

type Conf struct {
	Plus     int64
	Minus    int64
	Multiply int64
	Divide   int64
	MaxID    uint64
}

// Add — добавляет запись в базу данных
func Save(conf *config.Config) {
	/* go func() {
		db := postgresql.DbConnect(conf)
		defer db.Close()
		query := fmt.Sprintf(`
		UPDATE %s
		SET %s = $1, %s = $2, %s = $3, %s = $4, %s = $5
		WHERE id = 1`, tableName, plus, minus, multiply, divide, maxid)
		sqlPrepare, err := db.Prepare(query)
		defer sqlPrepare.Close()
		if err != nil {
			return
		}
		newConf := Conf{}
		newConf.Init(conf)
		_, err = sqlPrepare.Query(
			newConf.Plus,
			newConf.Minus,
			newConf.Multiply,
			newConf.Divide,
			newConf.MaxID,
		)
		if err != nil {
			slog.Info("Не удалось обновить конфигурацию", "ошибка:", err)
			return
		}

		slog.Info("Обновлена запись конфигурации", "конфиг:", newConf)
	}() */
}

// Ищем созданную конфигурацию в базе данных. Ищем по ID = 1
func GetByIdOne(conf *config.Config) (Conf, error) {
	/* db := postgresql.DbConnect(conf)
	defer db.Close()
	query := fmt.Sprintf("SELECT %s, %s, %s, %s, %s FROM %s WHERE id = $1", plus, minus, multiply, divide, maxid, tableName)
	prepare, err := db.Prepare(query)
	if err != nil {
		return Conf{}, err
	}
	result := Conf{}
	row := prepare.QueryRow(1)
	err = row.Scan(
		&result.Plus,
		&result.Minus,
		&result.Multiply,
		&result.Divide,
		&result.MaxID,
	)
	if errors.Is(sql.ErrNoRows, err) {
		create(conf)
		slog.Info("Создание новой конфигурации")
		return GetByIdOne(conf)
	}
	if err != nil {
		return Conf{}, err
	}
	*/
	return Conf{}, nil
}

// Если нет конфигурации с ID = 1, то создаем
func create(conf *config.Config) {
	/* db := postgresql.DbConnect(conf)
	defer db.Close()
	query := fmt.Sprintf(`
			INSERT INTO %s (%s, %s, %s, %s, %s)
			VALUES ($1, $2, $3, $4, $5)`,
		tableName, plus, minus, multiply, divide, maxid)
	sqlPrepare, err := db.Prepare(query)
	if err != nil {
		return
	}
	defer sqlPrepare.Close()
	sqlPrepare.Query(0, 0, 0, 0, 0) */
}

// Инициализируем конфиг
func (c *Conf) Init(conf *config.Config) {
	conf.Lock()
	defer conf.Unlock()
	c.Plus = conf.Plus
	c.Minus = conf.Minus
	c.Multiply = conf.Multiply
	c.Divide = conf.Divide
	c.MaxID = conf.MaxID
}

// Загружаем сохраненную конфигурацию во время старта приложения
// Пытаемся подключиться к базе данных. Если данные будут изменены и потом будет
// подключена к базе данных, то конфиг будет перезаписан
func Load(conf *config.Config) {
	go func() {
		duration, err := GetByIdOne(conf)
		if err != nil {
			slog.Error("Не удалось загрузить конфигурацию", "ошибка:", err)
		}
		conf.Lock()
		defer conf.Unlock()
		conf.Plus = duration.Plus
		conf.Minus = duration.Minus
		conf.Multiply = duration.Multiply
		conf.Divide = duration.Divide
		conf.MaxID = duration.MaxID
		slog.Info("Загружена конфигурация", "конфиг:", conf)
	}()
}
