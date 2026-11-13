package grpcapi

import (
	"context"

	identityv1 "github.com/alexnesterov/rapidlog-api/api/gen/identity/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *IdentityService) ResolveSession(ctx context.Context, req *identityv1.ResolveSessionRequest) (*identityv1.ResolveSessionResponse, error) {
	id, _ := uuid.Parse(req.GetSessionId())

	userID, err := s.usecase.ResolveSession.Execute(ctx, id)
	if err != nil {
		s.logger.Error("failed to resolve session", "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &identityv1.ResolveSessionResponse{UserId: userID.String()}, nil
}
