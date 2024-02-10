package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Host          string
	Port          string
	MaxGoroutines int
}

type ConfigExpression struct {
	Plus     int `json:"plus"`
	Minus    int `json:"minus"`
	Multiply int `json:"multiply"`
	Divide   int `json:"divide"`
}

func New() *Config {
	godotenv.Load("./config/.env")
	maxGoroutines := os.Getenv("MAX_GOROUTINES_AGENT")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "7777"
	}
	if maxGoroutines == "" {
		maxGoroutines = "5"
	}
	config := &Config{
		Host:          host,
		Port:          port,
		MaxGoroutines: parseInt(maxGoroutines),
	}
	slog.Info("Установлен новый конфиг:", "конфиг", config)
	return config
}

func parseInt(val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return i
}
