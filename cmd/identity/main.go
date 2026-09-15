package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexnesterov/rapidlog-api/internal/service/identity"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := identity.Run(ctx, logger); err != nil {
		logger.Error("failed to run service", "error", err)
		os.Exit(1)
	}
}
