package grpcapi

import (
	"context"

	identityv1 "github.com/alexnesterov/rapidlog-api/api/gen/identity/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Usecase interface {
	ResolveSession(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
}

type IdentityService struct {
	identityv1.UnimplementedIdentityServiceServer
	usecase Usecase
}

func NewIdentityService(usecase Usecase) *IdentityService {
	return &IdentityService{
		usecase: usecase,
	}
}

func RegisterIdentityService(srv *grpc.Server, uc Usecase) {
	identityService := NewIdentityService(uc)
	identityv1.RegisterIdentityServiceServer(srv, identityService)
}
