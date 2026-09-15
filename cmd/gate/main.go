package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/alexnesterov/rapidlog-api/internal/service/gate"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := gate.Run(ctx, logger); err != nil {
		logger.Error("failed to run service", "error", err)
		os.Exit(1)
	}
}
