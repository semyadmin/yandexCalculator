package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/http/server"
)

func main() {
	slog.Info("Orchestrator started")
	ctx, cancel := context.WithCancel(context.Background())
	shutDown, err := server.Run(ctx)
	if err != nil {
		slog.Error("Server error", "error", err)
		os.Exit(1)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-c
	cancel()
	shutDown(ctx)
	slog.Info("Orchestrator stopped")
	os.Exit(0)
}
