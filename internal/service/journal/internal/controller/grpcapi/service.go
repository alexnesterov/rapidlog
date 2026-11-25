package grpcapi

import (
	"context"
	"log/slog"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type CreateBulletUsecase interface {
	Execute(ctx context.Context, in usecase.CreateBulletInput) (*entity.Bullet, error)
}

type CancelBulletUsecase interface {
	Execute(ctx context.Context, in usecase.CancelBulletInput) (*entity.Bullet, error)
}

type LegacyUsecase interface {
	ListBullets(ctx context.Context, userID uuid.UUID) ([]*entity.Bullet, error)
	CompleteBullet(ctx context.Context, id, userID uuid.UUID) (*entity.Bullet, error)
	MigrateBullet(ctx context.Context, id, userID uuid.UUID) (*entity.Bullet, error)
}

type Usecases struct {
	CreateBullet CreateBulletUsecase
	CancelBullet CancelBulletUsecase
	Legacy       LegacyUsecase
}

type JournalService struct {
	journalv1.UnimplementedJournalServiceServer
	uc     Usecases
	logger *slog.Logger
}

func NewJournalService(uc Usecases, logger *slog.Logger) *JournalService {
	return &JournalService{
		uc:     uc,
		logger: logger,
	}
}

func RegisterJournalService(srv *grpc.Server, uc Usecases, logger *slog.Logger) {
	journalService := NewJournalService(uc, logger)
	journalv1.RegisterJournalServiceServer(srv, journalService)
}
