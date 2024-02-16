package config

import (
	"errors"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var ErrWrongDuration = errors.New("Некорректное время для операций")

// Config - конфигурация приложения
type Config struct {
	Host     string
	HttpPort string
	TCPPort  string
	Db       string
	DbName   string
	DbPort   string
	DbUser   string
	DbPass   string
	Plus     int64
	Minus    int64
	Multiply int64
	Divide   int64
	MaxID    uint64
	sync.Mutex
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
	dbPassword := os.Getenv("ORCHESTRATOR_DB__PASSWORD")
	if httpPort == "" {
		httpPort = "8080"
	}
	if tcpPort == "" {
		tcpPort = "7777"
	}
	if db == "" {
		tcpPort = "localhost"
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
	return &Config{
		Host:     "localhost",
		HttpPort: httpPort,
		TCPPort:  tcpPort,
		Db:       db,
		DbName:   dbName,
		DbPort:   dbPort,
		DbUser:   dbUser,
		DbPass:   dbPassword,
	}
}

func (c *Config) NewDuration(conf *ConfigExpression) error {
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
