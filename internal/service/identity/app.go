package identity

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/service/identity/config"
	"github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/adapter/grpcapi"
	"github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/adapter/postgres"
	"github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/usecase"
	"google.golang.org/grpc"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	pool, err := postgres.Connect(ctx, cfg.DSN)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	logger.Info("connected to postgres")

	if err := postgres.Migrate(cfg.DSN); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	repository := postgres.NewUserRepository(pool)
	uc := usecase.New(repository)

	grpcServer := grpc.NewServer()
	grpcapi.RegisterIdentityService(grpcServer, uc)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	go func() {
		<-ctx.Done()
		logger.Info("shutting down identity gRPC server")

		stopped := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:
			logger.Info("graceful shutdown completed")
		case <-time.After(15 * time.Second):
			logger.Info("graceful shutdown timed out, forcing stop")
			grpcServer.Stop()
		}
	}()

	logger.Info("starting identity grpc server")
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to start grpc server: %w", err)
	}

	return nil
}
