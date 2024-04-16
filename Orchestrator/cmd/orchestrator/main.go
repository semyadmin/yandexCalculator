package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	grpcserver "github.com/adminsemy/yandexCalculator/Orchestrator/internal/grpc_server"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/http/server"
	loadfromdb "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/load_from_db"
	sendtocalculate "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/send_to_calculate"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
)

func main() {
	// Создаем конфигурацию
	conf := config.New("config/.env")
	// Создаем сторадж для хранения выражений в памяти
	storage := memory.New(conf)
	// Создаем сторадж для хранения пользователей в памяти
	userStorage := memory.NewUserStorage(conf)
	// Создаем новую очередь
	newQueue := queue.NewMapQueue(queue.NewLockFreeQueue(), conf)
	// Загружаем все из базы
	loadfromdb.LoadFromDB(conf, storage, userStorage, newQueue)
	// Запускаем GRPC сервер
	ctx, cancel := context.WithCancel(context.Background())
	grpcServer := grpcserver.NewServerGRPC(conf, sendtocalculate.NewSendToCalculate(newQueue))
	go grpcServer.Start()
	slog.Info("Оркестратор запущен")
	// Получаем функцию для остановки HTTP сервера
	shutDown, err := server.Run(ctx, conf, newQueue, storage, userStorage)
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
	// Останавливаем GRPC сервер
	grpcServer.Stop()
	shutDown(ctx)
	slog.Info("Оркестратор остановлен")
	os.Exit(0)
}
