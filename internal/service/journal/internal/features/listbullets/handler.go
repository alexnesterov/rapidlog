package listbullets

import (
	"context"
	"log/slog"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/protomap"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Executor interface {
	Execute(ctx context.Context, in Input) ([]*entity.Bullet, error)
}

type Handler struct {
	uc     Executor
	logger *slog.Logger
}

func NewHandler(uc Executor, logger *slog.Logger) *Handler {
	return &Handler{uc: uc, logger: logger}
}

func (h *Handler) Handle(ctx context.Context, req *journalv1.ListBulletsRequest) (*journalv1.ListBulletsResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		h.logger.Warn("invalid user id", "user_id", req.GetUserId(), "error", err)
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	bullets, err := h.uc.Execute(ctx, Input{UserID: userID})
	if err != nil {
		h.logger.Error("failed to list bullets", "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	data, err := protomap.ToProtoBullets(bullets)
	if err != nil {
		h.logger.Error("failed to convert bullets to proto",
			"user_id", userID,
			"error", err,
		)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &journalv1.ListBulletsResponse{
		Data: data,
	}, nil
}
