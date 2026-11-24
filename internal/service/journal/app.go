package journal

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/pkg/transactor"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/config"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/adapter/postgres"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/controller/grpcapi"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase"
	"google.golang.org/grpc"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	txMgr := transactor.NewTxManager(pool)
	repository := postgres.NewBulletRepository(pool)
	uc := usecase.New(repository, txMgr)

	grpcServer := grpc.NewServer()
	grpcapi.RegisterJournalService(grpcServer, uc, logger)

	ln, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("gRPC server listening", "port", cfg.Port)
		if err := grpcServer.Serve(ln); err != nil {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	case <-ctx.Done():
		logger.Info("gRPC server shutting down")
	}

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		logger.Info("gRPC server stopped")
	case <-time.After(5 * time.Second):
		logger.Info("gRPC server stop timed out, forcing stop")
		grpcServer.Stop()
	}

	return nil
}
