package config

import "os"

// Config - конфигурация приложения
type Config struct {
	Host     string
	Port     string
	Plus     int
	Minus    int
	Multiply int
	Divide   int
}

func New() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return &Config{
		Host:     "localhost",
		Port:     port,
		Plus:     1,
		Minus:    1,
		Multiply: 1,
		Divide:   1,
	}
}
