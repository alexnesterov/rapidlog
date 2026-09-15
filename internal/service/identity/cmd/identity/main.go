package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	identityv1 "github.com/alexnesterov/rapidlog-api/gen/identity/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/adapter/grpcapi"
	"github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/domain/usecase"
	"github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/infrastructure/postgres"
	"google.golang.org/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dsn := os.Getenv("IDENTITY_DB_DSN")

	pool, err := postgres.Connect(ctx, dsn)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	logger.Info("connected to postgres")

	if err := postgres.Migrate(dsn); err != nil {
		logger.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	userRepository := postgres.NewUserRepository(pool)
	identityService := usecase.NewIdentityService(userRepository)

	handler := grpcapi.NewHandler(identityService)

	grpcServer := grpc.NewServer()
	identityv1.RegisterIdentityServer(grpcServer, handler)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("failed to listen", "error", err)
		os.Exit(1)
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
		logger.Error("failed to start grpc server", "error", err)
		os.Exit(1)
	}
}
