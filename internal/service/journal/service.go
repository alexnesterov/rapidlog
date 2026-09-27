package journal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/config"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger.Info("journal service dsn", "dsn", cfg.DSN)
	logger.Info("journal service port", "port", cfg.Port)

	return nil
}
