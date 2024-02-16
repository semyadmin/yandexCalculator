package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/http/server/handlers"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
)

func Run(ctx context.Context,
	config *config.Config,
	queue *queue.MapQueue,
	storage *memory.Storage,
) (func(context.Context) error, error) {
	// Инициализируем маршрутизатор
	serveMux, err := handlers.NewServeMux(config, queue, storage)
	if err != nil {
		return nil, err
	}
	serveMux = handlers.Decorate(serveMux, logMiddleware())
	// Инициализируем HTTP сервер
	httpServer := &http.Server{Addr: ":" + config.HttpPort, Handler: serveMux}

	slog.Info("Http сервер запущен на порту " + config.HttpPort)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Ошибка сервера:", "error", err)
		}
	}()
	// Возвращаем функцию для остановки сервера
	return httpServer.Shutdown, nil
}

// Логгируем все приходящие запросы
func logMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			slog.Info("Запрос получен", "Метод:", r.Method, "Путь:", r.URL.Path, "Продолжительность:", duration)
		})
	}
}
