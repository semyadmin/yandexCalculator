package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/http/server"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
	serverTCP "github.com/adminsemy/yandexCalculator/Orchestrator/internal/tcp/server"
)

func main() {
	config := config.New()
	newQueue := queue.NewLockFreeQueue()
	tcpServer, err := serverTCP.NewServer(":7777", config, newQueue)
	if err != nil {
		slog.Error("Ошибка запуска TCP/IP сервера:", "ошибка:", err)
		os.Exit(1)
	}
	slog.Info("Запуск TCP/IP сервера на порту 7777")
	tcpServer.Start()
	slog.Info("Оркестратор запущен")
	ctx, cancel := context.WithCancel(context.Background())
	shutDown, err := server.Run(ctx, config, newQueue)
	if err != nil {
		slog.Error("Ошибка запуска сервера:", "ошибка:", err)
		os.Exit(1)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-c
	cancel()
	tcpServer.Stop()
	shutDown(ctx)
	slog.Info("Сервер TCP/IP остановлен")
	slog.Info("Оркестратор остановлен")
	os.Exit(0)
}
