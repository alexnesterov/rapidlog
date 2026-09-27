package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	journalLogger := logger.With("service", "journal")

	if err := journal.Run(ctx, journalLogger); err != nil {
		logger.Error("failed to run service", "error", err)
		os.Exit(1)
	}
}
