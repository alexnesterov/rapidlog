package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).
		With("service", "journal")

	if err := journal.Run(context.Background(), logger); err != nil {
		logger.Error("failed to run service", "error", err)
		os.Exit(1)
	}
}
