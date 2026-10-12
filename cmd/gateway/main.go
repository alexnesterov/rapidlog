package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/alexnesterov/rapidlog-api/internal/service/gateway"
)

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	gatewayLogger := logger.With("service", "gateway")

	if err := gateway.Run(ctx, gatewayLogger); err != nil {
		logger.Error("failed to run service", "error", err)
		os.Exit(1)
	}
}
