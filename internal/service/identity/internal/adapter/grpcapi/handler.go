package grpcapi

import (
	"context"

	identityv1 "github.com/alexnesterov/rapidlog-api/gen/identity/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/domain/port"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	identityv1.UnimplementedIdentityServiceServer
	usecase port.IdentityService
}

func NewHandler(service port.IdentityService) *Handler {
	return &Handler{
		usecase: service,
	}
}

func (s *Handler) ResolveSession(ctx context.Context, req *identityv1.ResolveSessionRequest) (*identityv1.ResolveSessionResponse, error) {
	id, _ := uuid.Parse(req.GetSessionId())

	userID, err := s.usecase.ResolveSession(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &identityv1.ResolveSessionResponse{UserId: userID.String()}, nil
}
