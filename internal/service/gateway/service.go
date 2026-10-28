package gateway

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/config"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/controller/httpapi"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/controller/httpapi/middleware"
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

	router := http.NewServeMux()
	router.HandleFunc("/health", httpapi.NewHealthHandler())
	router.HandleFunc("POST /api/bullets", httpapi.NewCreateBulletHandler(journalClient, logger))
	router.HandleFunc("GET /api/bullets", httpapi.NewListBulletsHandler(journalClient))
	router.HandleFunc("POST /api/bullets/{id}/complete", httpapi.EmptyHandler)
	router.HandleFunc("POST /api/bullets/{id}/migrate", httpapi.EmptyHandler)
	router.HandleFunc("POST /api/bullets/{id}/cancel", httpapi.EmptyHandler)

	var handler http.Handler = router
	handler = middleware.Session(identityClient, cfg.Session.CookieName, cfg.Session.CookieTTL, cfg.Session.CookieSecure)(handler)
	handler = middleware.Logging(logger)(handler)
	handler = middleware.Recovery(logger)(handler)

	frontend, err := fs.Sub(web.Dist, "dist")
	if err != nil {
		return fmt.Errorf("failed to load embedded frontend: %w", err)
	}
	router.Handle("/", http.FileServer(http.FS(frontend)))

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
