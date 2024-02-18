package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/http/server"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/postgresql/postgresql_ast"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/postgresql/postgresql_config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
	serverTCP "github.com/adminsemy/yandexCalculator/Orchestrator/internal/tcp/server"
)

func main() {
	// Создаем конфигурацию
	conf := config.New()
	// Загружаем конфигурацию из базы
	postgresql_config.Load(conf)
	// Создаем сторадж для хранения выражений в памяти
	storage := memory.New(conf)
	// Создаем новую очередь
	newQueue := queue.NewMapQueue(queue.NewLockFreeQueue(), conf)
	// Загружаем все выражения из базы
	postgresql_ast.GetAll(conf, newQueue, storage)
	// Горутина для обновления результатов выражений
	postgresql_ast.Update(conf, newQueue, storage)
	// Запускаем TCP/IP сервер
	tcpServer, err := serverTCP.NewServer(":"+conf.TCPPort, conf, newQueue, storage)
	if err != nil {
		slog.Error("Ошибка запуска TCP/IP сервера:", "ошибка:", err)
		os.Exit(1)
	}
	slog.Info("Запуск TCP/IP сервера на порту " + conf.TCPPort)
	tcpServer.Start()
	slog.Info("Оркестратор запущен")
	ctx, cancel := context.WithCancel(context.Background())
	// Получаем функцию для остановки HTTP сервера
	shutDown, err := server.Run(ctx, conf, newQueue, storage)
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
