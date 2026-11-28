package listbullets

import (
	"context"
	"log/slog"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/adapter/postgres"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/platform/server"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListBulletsUsecase interface {
	Execute(ctx context.Context, in usecase.ListBulletsInput) ([]*domain.Bullet, error)
}

type Handler struct {
	uc     ListBulletsUsecase
	logger *slog.Logger
}

func NewHandler(uc ListBulletsUsecase, logger *slog.Logger) *Handler {
	return &Handler{
		uc:     uc,
		logger: logger,
	}
}

func Setup(pool *pgxpool.Pool, logger *slog.Logger) *Handler {
	r := postgres.NewBulletRepository(pool)
	uc := usecase.NewListBullets(r)
	return NewHandler(uc, logger)
}

func (h *Handler) Handle(ctx context.Context, req *journalv1.ListBulletsRequest) (*journalv1.ListBulletsResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		h.logger.Warn("invalid user id", "user_id", req.GetUserId(), "error", err)
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	bullets, err := h.uc.Execute(ctx, usecase.ListBulletsInput{UserID: userID})
	if err != nil {
		h.logger.Error("failed to list bullets", "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	data, err := server.ToProtoBullets(bullets)
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
