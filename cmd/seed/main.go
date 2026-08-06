package main

import (
	"context"
	_ "embed"
	"log/slog"
	"os"

	"github.com/alexnesterov/rapidlog-api/internal/config"
	"github.com/jackc/pgx/v5"
)

//go:embed seed.sql
var seedSQL []byte

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := pgx.Connect(ctx, cfg.DB.DSN)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close(ctx)

	_, err = db.Exec(ctx, string(seedSQL))
	if err != nil {
		slog.Error("failed to execute seed.sql", "error", err)
		os.Exit(1)
	}

	slog.Info("seed.sql executed successfully")
}
