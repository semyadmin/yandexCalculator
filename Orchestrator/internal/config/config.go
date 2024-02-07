package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config - конфигурация приложения
type Config struct {
	Host     string
	HttpPort string
	TCPPort  string
}

type ConfigExpression struct {
	Plus     int `json:"plus"`
	Minus    int `json:"minus"`
	Multiply int `json:"multiply"`
	Divide   int `json:"divide"`
}

func New() *Config {
	httpPort := os.Getenv("ORCHESTRATOR_HTTP_PORT")
	tcpPort := os.Getenv("ORCHESTRATOR_TCP_PORT")
	godotenv.Load("./config/.env")
	if httpPort == "" {
		httpPort = os.Getenv("ORCHESTRATOR_HTTP_PORT")
		if httpPort == "" {
			httpPort = "8080"
		}
	}
	if tcpPort == "" {
		tcpPort = os.Getenv("ORCHESTRATOR_TCP_PORT")
		if tcpPort == "" {
			tcpPort = "7777"
		}
	}
	return &Config{
		Host:     "localhost",
		HttpPort: httpPort,
		TCPPort:  tcpPort,
	}
}
