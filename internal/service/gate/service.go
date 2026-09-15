package gate

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"

	identityv1 "github.com/alexnesterov/rapidlog-api/gen/identity/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/gate/internal/adapter/httpapi"
	"github.com/alexnesterov/rapidlog-api/internal/service/gate/internal/adapter/httpapi/middleware"
	"github.com/alexnesterov/rapidlog-api/internal/service/gate/internal/config"
	"github.com/alexnesterov/rapidlog-api/internal/service/gate/internal/domain/usecase"
	"github.com/alexnesterov/rapidlog-api/internal/service/gate/internal/infrastructure/identity"
	"github.com/alexnesterov/rapidlog-api/internal/service/gate/internal/infrastructure/postgres"
	"github.com/alexnesterov/rapidlog-api/web"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	pool, err := postgres.Connect(ctx, cfg.DB.DSN)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	logger.Info("connected to postgres")

	if err := postgres.Migrate(cfg.DB.DSN); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	txMgr := postgres.NewTransactionManager(pool)

	bulletRepository := postgres.NewBulletRepository(pool)
	bulletService := usecase.NewBulletService(bulletRepository, txMgr)
	bulletHandler := httpapi.NewBulletHandler(bulletService)

	// userRepository := postgres.NewUserRepository(pool)
	// userService := usecase.NewUserService(
	// 	userRepository,
	// 	bulletRepository,
	// 	txMgr,
	// )
	identityConn, err := grpc.NewClient(cfg.Identity.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to identity service: %w", err)
	}
	defer identityConn.Close()

	identityClient := identity.NewClient(identityv1.NewIdentityServiceClient(identityConn))

	router := http.NewServeMux()
	router.HandleFunc("/health", httpapi.NewHealthHandler(pool))
	router.HandleFunc("POST /api/bullets", bulletHandler.CreateBullet)
	router.HandleFunc("GET /api/bullets", bulletHandler.ListBullets)
	router.HandleFunc("POST /api/bullets/{id}/complete", bulletHandler.CompleteBullet)
	router.HandleFunc("POST /api/bullets/{id}/migrate", bulletHandler.MigrateBullet)
	router.HandleFunc("POST /api/bullets/{id}/cancel", bulletHandler.CancelBullet)

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
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      handler,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	logger.Info("starting server", "name", cfg.App.Name, "port", cfg.HTTP.Port)
	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
