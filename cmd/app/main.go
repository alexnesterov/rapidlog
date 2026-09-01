package main

import (
	"context"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"github.com/alexnesterov/rapidlog-api/internal/adapter/httpapi"
	"github.com/alexnesterov/rapidlog-api/internal/adapter/httpapi/middleware"
	"github.com/alexnesterov/rapidlog-api/internal/config"
	"github.com/alexnesterov/rapidlog-api/internal/domain/usecase"
	"github.com/alexnesterov/rapidlog-api/internal/infrastructure/postgres"
	"github.com/alexnesterov/rapidlog-api/web"
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

	txMgr := postgres.NewTransactionManager(pool)

	bulletRepository := postgres.NewBulletRepository(pool)
	bulletService := usecase.NewBulletService(bulletRepository, txMgr)
	bulletHandler := httpapi.NewBulletHandler(bulletService)

	userRepository := postgres.NewUserRepository(pool)
	userService := usecase.NewUserService(
		userRepository,
		bulletRepository,
		txMgr,
	)

	router := http.NewServeMux()
	router.HandleFunc("/health", httpapi.NewHealthHandler(pool))
	router.HandleFunc("POST /api/bullets", bulletHandler.CreateBullet)
	router.HandleFunc("GET /api/bullets", bulletHandler.ListBullets)
	router.HandleFunc("POST /api/bullets/{id}/complete", bulletHandler.CompleteBullet)
	router.HandleFunc("POST /api/bullets/{id}/migrate", bulletHandler.MigrateBullet)

	var handler http.Handler = router
	handler = middleware.Session(userService, cfg.Session.CookieName, cfg.Session.CookieTTL, cfg.Session.CookieSecure)(handler)
	handler = middleware.Logging(logger)(handler)
	handler = middleware.Recovery(logger)(handler)

	frontend, err := fs.Sub(web.Dist, "dist")
	if err != nil {
		logger.Error("failed to load embedded frontend", "error", err)
		os.Exit(1)
	}
	router.Handle("/", http.FileServer(http.FS(frontend)))

	server := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      handler,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	logger.Info("starting server", "name", cfg.App.Name, "port", cfg.HTTP.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
