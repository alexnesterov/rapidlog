package gateway

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/config"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/adapter/httpapi"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/adapter/httpapi/middleware"
	"github.com/alexnesterov/rapidlog-api/internal/service/identity/pkg/identitygrpc"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/pkg/journalgrpc"
	"github.com/alexnesterov/rapidlog-api/web"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	identityClient, err := identitygrpc.NewClient(cfg.Identity)
	if err != nil {
		return fmt.Errorf("identity new client: %w", err)
	}
	defer identityClient.Close()

	journalClient, err := journalgrpc.NewClient(cfg.Journal)
	if err != nil {
		return fmt.Errorf("journal new client: %w", err)
	}
	defer journalClient.Close()

	api := http.NewServeMux()
	api.Handle("POST /api/bullets", httpapi.NewCreateBulletHandler(journalClient, logger))
	api.Handle("GET /api/bullets", httpapi.NewListBulletsHandler(journalClient))
	api.Handle("POST /api/bullets/{id}/complete", httpapi.NewCompleteBulletHandler(journalClient, logger))
	api.Handle("POST /api/bullets/{id}/migrate", httpapi.NewMigrateBulletHandler(journalClient, logger))
	api.Handle("POST /api/bullets/{id}/cancel", httpapi.NewCancelBulletHandler(journalClient, logger))

	router := http.NewServeMux()
	router.HandleFunc("/health", httpapi.NewHealthHandler())
	router.Handle("/", web.Handler())

	router.Handle("/api/",
		middleware.Session(identityClient, cfg.Session.CookieName, cfg.Session.CookieTTL, cfg.Session.CookieSecure, logger)(api),
	)

	var handler http.Handler = router
	handler = middleware.Logging(logger)(handler)
	handler = middleware.Recovery(logger)(handler)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	logger.Info("starting server", "port", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
