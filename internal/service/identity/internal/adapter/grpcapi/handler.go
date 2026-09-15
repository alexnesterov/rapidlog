package grpcapi

import (
	"context"

	identityv1 "github.com/alexnesterov/rapidlog-api/gen/identity/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/domain/port"
	"github.com/google/uuid"
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
	id, _ := uuid.Parse(req.SessionId)

	userID, err := s.usecase.ResolveSession(ctx, id)
	if err != nil {
		return nil, err
	}

	return &identityv1.ResolveSessionResponse{UserId: userID.String()}, nil
}
