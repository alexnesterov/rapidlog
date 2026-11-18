package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/alexnesterov/rapidlog-api/internal/service/gateway"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).
		With("service", "gateway")

	if err := gateway.Run(context.Background(), logger); err != nil {
		logger.Error("failed to run service", "error", err)
		os.Exit(1)
	}
}
