package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/alexnesterov/rapidlog-api/internal/adapter/httpapi"
	"github.com/alexnesterov/rapidlog-api/internal/adapter/httpapi/middleware"
	"github.com/alexnesterov/rapidlog-api/internal/config"
	"github.com/alexnesterov/rapidlog-api/internal/domain/usecase"
	"github.com/alexnesterov/rapidlog-api/internal/infrastructure/postgres"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	pool, err := postgres.Connect(ctx, cfg.DB.DSN)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	logger.Info("connected to postgres")

	if err := postgres.Migrate(cfg.DB.DSN); err != nil {
		logger.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	if err := postgres.Seed(ctx, pool); err != nil {
		logger.Error("failed to run seed", "error", err)
		os.Exit(1)
	}

	txMgr := postgres.NewTransactionManager(pool)

	router := http.NewServeMux()
	router.HandleFunc("/health", httpapi.NewHealthHandler(pool))

	var handler http.Handler = router
	handler = middleware.Logging(logger)(handler)
	handler = middleware.Recovery(logger)(handler)

	bulletRepository := postgres.NewBulletRepository(pool)
	bulletService := usecase.NewBulletService(bulletRepository, txMgr)
	bulletHandler := httpapi.NewBulletHandler(bulletService)

	router.HandleFunc("POST /api/bullets", bulletHandler.CreateBullet)
	router.HandleFunc("GET /api/bullets", bulletHandler.ListBullets)
	router.HandleFunc("POST /api/bullets/{id}/complete", bulletHandler.CompleteBullet)
	router.HandleFunc("POST /api/bullets/{id}/migrate", bulletHandler.MigrateBullet)

	logger.Info("starting server", "name", cfg.App.Name, "port", cfg.HTTP.Port)
	if err := http.ListenAndServe(":"+cfg.HTTP.Port, handler); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
