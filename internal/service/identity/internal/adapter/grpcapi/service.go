package grpcapi

import (
	"context"
	"log/slog"

	identityv1 "github.com/alexnesterov/rapidlog-api/api/gen/identity/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type ResolveSessionUsecase interface {
	Execute(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
}

type Usecase struct {
	ResolveSession ResolveSessionUsecase
}

type IdentityService struct {
	identityv1.UnimplementedIdentityServiceServer
	usecase Usecase
	logger  *slog.Logger
}

func NewIdentityService(usecase Usecase, logger *slog.Logger) *IdentityService {
	return &IdentityService{
		usecase: usecase,
		logger:  logger,
	}
}

func RegisterIdentityService(srv *grpc.Server, uc Usecase, logger *slog.Logger) {
	identityService := NewIdentityService(uc, logger)
	identityv1.RegisterIdentityServiceServer(srv, identityService)
}
