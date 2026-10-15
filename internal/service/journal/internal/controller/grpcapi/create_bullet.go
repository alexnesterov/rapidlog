package grpcapi

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/pkg/protoenum"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/port"
)

func (s *JournalServer) CreateBullet(
	ctx context.Context,
	req *journalv1.CreateBulletRequest,
) (*journalv1.CreateBulletResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		s.Logger.Warn("invalid user id", "user_id", req.GetUserId(), "error", err)
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	input := port.CreateBulletInput{
		UserID:  userID,
		Type:    entity.BulletType(protoenum.ToCleanString(req.GetType())),
		Content: req.GetContent(),
	}

	bullet, err := s.Usecase.CreateBullet(ctx, input)
	if err != nil {
		var validationErr *entity.ValidationError
		if errors.As(err, &validationErr) {
			s.Logger.Warn("invalid create bullet request", "error", err)
			return nil, status.Error(codes.InvalidArgument, validationErr.Error())
		}

		s.Logger.Error("failed to create bullet", "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	data, err := ToProtoBullet(bullet)
	if err != nil {
		s.Logger.Error("failed to convert bullet to proto",
			"bullet_id", bullet.ID,
			"user_id", bullet.UserID,
			"error", err,
		)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &journalv1.CreateBulletResponse{
		Data: data,
	}, nil
}
