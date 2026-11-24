package gateway

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/config"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/adapter/kratos"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/controller/httpapi"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/controller/httpapi/middleware"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/pkg/journalgrpc"
	"github.com/alexnesterov/rapidlog-api/web"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	kratosClient, err := kratos.NewClient(cfg.Kratos.PublicURL)
	if err != nil {
		return fmt.Errorf("kratos new client: %w", err)
	}

	kratosProxy, err := kratos.NewProxy(cfg.Kratos.PublicURL)
	if err != nil {
		return fmt.Errorf("kratos new proxy: %w", err)
	}

	journalClient, err := journalgrpc.NewClient(cfg.Journal)
	if err != nil {
		return fmt.Errorf("journal new client: %w", err)
	}
	defer func() { _ = journalClient.Close() }()

	api := http.NewServeMux()
	api.Handle("POST /api/bullets", httpapi.NewCreateBullet(journalClient, logger))
	api.Handle("GET /api/bullets", httpapi.NewListBullets(journalClient))
	api.Handle("POST /api/bullets/{id}/complete", httpapi.NewCompleteBullet(journalClient, logger))
	api.Handle("POST /api/bullets/{id}/migrate", httpapi.NewMigrateBullet(journalClient, logger))
	api.Handle("POST /api/bullets/{id}/cancel", httpapi.NewCancelBullet(journalClient, logger))

	router := http.NewServeMux()
	router.HandleFunc("/health", httpapi.NewHealthHandler())
	router.Handle("/", web.Handler())
	// публичный API Kratos на том же origin, что и SPA: cookie сессии и CSRF без CORS
	router.Handle("/auth/", http.StripPrefix("/auth", kratosProxy))

	router.Handle("/api/", middleware.Auth(kratosClient, logger)(api))

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

	errCh := make(chan error, 1)
	go func() {
		logger.Info("starting server", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("failed to start server: %w", err)
	case <-ctx.Done():
		logger.Info("server shutting down")
	}

	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	return nil
}
