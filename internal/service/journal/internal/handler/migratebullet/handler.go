package migratebullet

import (
	"context"
	"errors"
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

type MigrateBulletUsecase interface {
	Execute(ctx context.Context, in usecase.MigrateBulletInput) (*domain.Bullet, error)
}

type Handler struct {
	uc     MigrateBulletUsecase
	logger *slog.Logger
}

func NewHandler(uc MigrateBulletUsecase, logger *slog.Logger) *Handler {
	return &Handler{
		uc:     uc,
		logger: logger,
	}
}

func Setup(pool *pgxpool.Pool, logger *slog.Logger) *Handler {
	r := postgres.NewBulletRepository(pool)
	uc := usecase.NewMigrateBullet(r)
	return NewHandler(uc, logger)
}

func (h *Handler) Handle(
	ctx context.Context,
	req *journalv1.MigrateBulletRequest,
) (*journalv1.MigrateBulletResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		h.logger.Warn("invalid user id", "user_id", req.GetUserId(), "error", err)
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	bulletID, err := uuid.Parse(req.GetId())
	if err != nil {
		h.logger.Warn("invalid bullet id", "bullet_id", req.GetId(), "error", err)
		return nil, status.Error(codes.InvalidArgument, "invalid bullet id")
	}

	bullet, err := h.uc.Execute(ctx, usecase.MigrateBulletInput{
		ID:     bulletID,
		UserID: userID,
	})
	if err != nil {
		var validationErr *domain.ValidationError
		if errors.As(err, &validationErr) {
			h.logger.Warn("invalid migrate bullet request", "error", err)
			return nil, status.Error(codes.InvalidArgument, validationErr.Error())
		}

		h.logger.Error("failed to migrate bullet", "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	data, err := server.ToProtoBullet(bullet)
	if err != nil {
		h.logger.Error("failed to convert bullet to proto",
			"bullet_id", bullet.ID,
			"user_id", bullet.UserID,
			"error", err,
		)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &journalv1.MigrateBulletResponse{
		Data: data,
	}, nil
}
