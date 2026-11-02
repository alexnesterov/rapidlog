package grpcapi

import (
	"context"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *JournalService) ListBullets(ctx context.Context, req *journalv1.ListBulletsRequest) (*journalv1.ListBulletsResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		s.logger.Warn("invalid user id", "user_id", req.GetUserId(), "error", err)
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	bullets, err := s.usecase.ListBullets(ctx, userID)
	if err != nil {
		s.logger.Error("failed to list bullets", "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	data, err := ToProtoBullets(bullets)
	if err != nil {
		s.logger.Error("failed to convert bullets to proto",
			"user_id", userID,
			"error", err,
		)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &journalv1.ListBulletsResponse{
		Data: data,
	}, nil
}
