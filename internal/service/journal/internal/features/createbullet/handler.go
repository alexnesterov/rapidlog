package createbullet

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/pkg/protoenum"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/protomap"
)

type Executor interface {
	Execute(ctx context.Context, in Input) (*entity.Bullet, error)
}

type Handler struct {
	uc     Executor
	logger *slog.Logger
}

func NewHandler(uc Executor, logger *slog.Logger) *Handler {
	return &Handler{uc: uc, logger: logger}
}

func (h *Handler) Handle(
	ctx context.Context,
	req *journalv1.CreateBulletRequest,
) (*journalv1.CreateBulletResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		h.logger.Warn("invalid user id", "user_id", req.GetUserId(), "error", err)
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	kind, err := entity.ParseKind(protoenum.ToCleanString(req.GetKind()))
	if err != nil {
		h.logger.Warn("invalid kind", "kind", req.GetKind(), "error", err)
		return nil, status.Error(codes.InvalidArgument, "invalid kind")
	}

	input := Input{
		UserID:  userID,
		Kind:    kind,
		Content: req.GetContent(),
	}

	bullet, err := h.uc.Execute(ctx, input)
	if err != nil {
		var validationErr *entity.ValidationError
		if errors.As(err, &validationErr) {
			h.logger.Warn("invalid create bullet request", "error", err)
			return nil, status.Error(codes.InvalidArgument, validationErr.Error())
		}

		h.logger.Error("failed to create bullet", "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	data, err := protomap.ToProtoBullet(bullet)
	if err != nil {
		h.logger.Error("failed to convert bullet to proto",
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
