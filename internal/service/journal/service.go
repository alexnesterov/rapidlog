package journal

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/config"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/adapter/postgres"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/controller/grpcapi"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase"
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

	txMgr := postgres.NewTransactionManager(pool)
	repository := postgres.NewBulletRepository(pool)
	service := usecase.NewBulletService(repository, txMgr)

	journalServer := &grpcapi.JournalServer{
		Usecase: service,
	}

	ln, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	journalv1.RegisterJournalServiceServer(grpcServer, journalServer)

	go func() {
		logger.Info("gRPC server listening", "port", cfg.Port)
		if err := grpcServer.Serve(ln); err != nil {
			logger.Error("failed to serve gRPC server", "error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("gRPC server shutting down")

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		logger.Info("gRPC server stopped")
	case <-time.After(15 * time.Second):
		logger.Info("gRPC server stop timed out, forcing stop")
		grpcServer.Stop()
	}

	return nil
}
