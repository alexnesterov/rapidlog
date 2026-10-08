package gate

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"

	identityv1 "github.com/alexnesterov/rapidlog-api/api/gen/identity/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/adapter/httpapi"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/adapter/httpapi/middleware"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/config"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/infra/identity"
	"github.com/alexnesterov/rapidlog-api/web"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	identityConn, err := grpc.NewClient(cfg.Identity, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to identity service: %w", err)
	}
	defer identityConn.Close()

	identityClient := identity.NewClient(identityv1.NewIdentityServiceClient(identityConn))

	router := http.NewServeMux()
	router.HandleFunc("/health", httpapi.NewHealthHandler())
	router.HandleFunc("POST /api/bullets", nil)
	router.HandleFunc("GET /api/bullets", nil)
	router.HandleFunc("POST /api/bullets/{id}/complete", nil)
	router.HandleFunc("POST /api/bullets/{id}/migrate", nil)
	router.HandleFunc("POST /api/bullets/{id}/cancel", nil)

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

	logger.Info("starting server", "name", "gate", "port", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
