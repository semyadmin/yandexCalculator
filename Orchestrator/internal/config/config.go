package config

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/postgresql"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/web_socket/manager"
	"github.com/joho/godotenv"
)

// Config - конфигурация приложения
type Config struct {
	Host        string
	HttpPort    string
	TCPPort     string
	Db          *postgresql.Storage
	Plus        int64
	Minus       int64
	Multiply    int64
	Divide      int64
	MaxID       uint64
	TokenLimit  uint64
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

func New(confFile string) *Config {
	err := godotenv.Load(confFile)
	if err != nil {
		slog.Error("Failed to load .env file", err)
	}
	httpPort := os.Getenv("ORCHESTRATOR_HTTP_PORT")
	tcpPort := os.Getenv("ORCHESTRATOR_TCP_PORT")
	db := os.Getenv("ORCHESTRATOR_DB")
	dbName := os.Getenv("ORCHESTRATOR_DB_NAME")
	dbPort := os.Getenv("ORCHESTRATOR_DB_PORT")
	dbUser := os.Getenv("ORCHESTRATOR_DB_USER")
	dbPassword := os.Getenv("ORCHESTRATOR_DB_PASSWORD")
	host := os.Getenv("ORCHESTRATOR_HOST")
	tokenLimit := os.Getenv("ORCHESTRATOR_TOKEN_LIMIT")
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
	if tokenLimit == "" {
		tokenLimit = "15"
	}
	t, _ := strconv.ParseUint(tokenLimit, 10, 64)
	return &Config{
		Host:       host,
		HttpPort:   httpPort,
		TCPPort:    tcpPort,
		TokenLimit: t,
		Db:         postgresql.NewPostgresConnect(db, dbPort, dbUser, dbPassword, dbName),
		WSmanager:  manager.NewManager(context.Background()),
	}
}
