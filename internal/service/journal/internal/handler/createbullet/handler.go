package createbullet

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/pkg/protoenum"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/adapter/postgres"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/platform/server"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase"
)

type CreateBulletUsecase interface {
	Execute(ctx context.Context, in usecase.CreateBulletInput) (*domain.Bullet, error)
}

type Handler struct {
	uc     CreateBulletUsecase
	logger *slog.Logger
}

func NewHandler(uc CreateBulletUsecase, logger *slog.Logger) *Handler {
	return &Handler{
		uc:     uc,
		logger: logger,
	}
}

func Setup(pool *pgxpool.Pool, logger *slog.Logger) *Handler {
	r := postgres.NewBulletRepository(pool)
	uc := usecase.NewCreateBullet(r)
	return NewHandler(uc, logger)
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

	kind, err := domain.ParseKind(protoenum.ToCleanString(req.GetKind()))
	if err != nil {
		h.logger.Warn("invalid kind", "kind", req.GetKind(), "error", err)
		return nil, status.Error(codes.InvalidArgument, "invalid kind")
	}

	input := usecase.CreateBulletInput{
		UserID:  userID,
		Kind:    kind,
		Content: req.GetContent(),
	}

	bullet, err := h.uc.Execute(ctx, input)
	if err != nil {
		var validationErr *domain.ValidationError
		if errors.As(err, &validationErr) {
			h.logger.Warn("invalid create bullet request", "error", err)
			return nil, status.Error(codes.InvalidArgument, validationErr.Error())
		}

		h.logger.Error("failed to create bullet", "error", err)
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

	return &journalv1.CreateBulletResponse{
		Data: data,
	}, nil
}
