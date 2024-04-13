package config

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/postgresql"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/web_socket/manager"
	"github.com/joho/godotenv"
)

var ErrWrongDuration = errors.New("Некорректное время для операций")

// Config - конфигурация приложения
type Config struct {
	Id          int64
	Host        string
	HttpPort    string
	TCPPort     string
	Db          *sql.DB
	Plus        int64
	Minus       int64
	Multiply    int64
	Divide      int64
	MaxID       uint64
	AgentsAll   atomic.Int64
	WorkersAll  atomic.Int64
	WorkersBusy atomic.Int64
	WSmanager   *manager.Manager
	sync.Mutex
}

type Workers struct {
	Agents      int64    `json:"agents"`
	Workers     int64    `json:"workers"`
	WorkersBusy int64    `json:"workersBusy"`
	Expressions []string `json:"expressions"`
}

type ConfigExpression struct {
	Plus     string `json:"plus"`
	Minus    string `json:"minus"`
	Multiply string `json:"multi"`
	Divide   string `json:"divide"`
}

func New() *Config {
	godotenv.Load("./config/.env")
	httpPort := os.Getenv("ORCHESTRATOR_HTTP_PORT")
	tcpPort := os.Getenv("ORCHESTRATOR_TCP_PORT")
	db := os.Getenv("ORCHESTRATOR_DB")
	dbName := os.Getenv("ORCHESTRATOR_DB_NAME")
	dbPort := os.Getenv("ORCHESTRATOR_DB_PORT")
	dbUser := os.Getenv("ORCHESTRATOR_DB_USER")
	dbPassword := os.Getenv("ORCHESTRATOR_DB_PASSWORD")
	host := os.Getenv("ORCHESTRATOR_HOST")
	if httpPort == "" {
		httpPort = "8080"
	}
	if tcpPort == "" {
		tcpPort = "7777"
	}
	if db == "" {
		db = "localhost"
	}
	if dbName == "" {
		dbName = "orchestrator"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbPassword == "" {
		dbPassword = "postgres"
	}
	if host == "" {
		host = "localhost"
	}
	return &Config{
		Host:      host,
		HttpPort:  httpPort,
		TCPPort:   tcpPort,
		Db:        postgresql.DbConnect(db, dbPort, dbUser, dbPassword, dbName),
		WSmanager: manager.NewManager(context.Background()),
	}
}

// Копируем конфигурацию из конфигурации от пользователя в нашу конфигурацию
func (c *Config) NewDuration(conf *ConfigExpression) error {
	c.Lock()
	defer c.Unlock()
	num, err := parseStringToInt(conf.Plus)
	if err != nil || num < 0 {
		return ErrWrongDuration
	}
	c.Plus = num
	num, err = parseStringToInt(conf.Minus)
	if err != nil || num < 0 {
		return ErrWrongDuration
	}
	c.Minus = num
	num, err = parseStringToInt(conf.Multiply)
	if err != nil || num < 0 {
		return ErrWrongDuration
	}
	c.Multiply = num
	num, err = parseStringToInt(conf.Divide)
	if err != nil || num < 0 {
		return ErrWrongDuration
	}
	c.Divide = num

	return nil
}

func (c *ConfigExpression) Init(conf *Config) {
	conf.Lock()
	defer conf.Unlock()
	c.Plus = strconv.FormatInt(int64(conf.Plus), 10)
	c.Minus = strconv.FormatInt(int64(conf.Minus), 10)
	c.Multiply = strconv.FormatInt(int64(conf.Multiply), 10)
	c.Divide = strconv.FormatInt(int64(conf.Divide), 10)
}

func parseStringToInt(str string) (int64, error) {
	num64, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		return 0, err
	}
	return num64, nil
}
