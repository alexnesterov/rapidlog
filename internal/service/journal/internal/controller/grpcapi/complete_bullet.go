package grpcapi

import (
	"context"
	"errors"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *JournalService) CompleteBullet(
	ctx context.Context,
	req *journalv1.CompleteBulletRequest,
) (*journalv1.CompleteBulletResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		s.logger.Warn("invalid user id", "user_id", req.GetUserId(), "error", err)
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	bulletID, err := uuid.Parse(req.GetId())
	if err != nil {
		s.logger.Warn("invalid bullet id", "bullet_id", req.GetId(), "error", err)
		return nil, status.Error(codes.InvalidArgument, "invalid bullet id")
	}

	bullet, err := s.usecase.CompleteBullet(ctx, bulletID, userID)
	if err != nil {
		var validationErr *entity.ValidationError
		if errors.As(err, &validationErr) {
			s.logger.Warn("invalid complete bullet request", "error", err)
			return nil, status.Error(codes.InvalidArgument, validationErr.Error())
		}

		s.logger.Error("failed to complete bullet", "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	data, err := ToProtoBullet(bullet)
	if err != nil {
		s.logger.Error("failed to convert bullet to proto",
			"bullet_id", bullet.ID,
			"user_id", bullet.UserID,
			"error", err,
		)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &journalv1.CompleteBulletResponse{
		Data: data,
	}, nil
}
