package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/http/server/handlers"
)

func Run(ctx context.Context) (func(context.Context) error, error) {
	// Инициализируем маршрутизатор
	serveMux, err := handlers.NewServeMux()
	if err != nil {
		return nil, err
	}
	// Инициализируем HTTP сервер
	httpServer := &http.Server{Addr: ":8080", Handler: serveMux}

	slog.Info("Server started")

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "error", err)
		}
	}()

	return httpServer.Shutdown, nil
}

func logMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			slog.Info("Request processed", "Method:", r.Method, "Path:", r.URL.Path, "Duration:", duration)
		})
	}
}
