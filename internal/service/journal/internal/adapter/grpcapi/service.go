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

type Usecase interface {
	CreateBullet(ctx context.Context, input usecase.CreateBulletInput) (*entity.Bullet, error)
	ListBullets(ctx context.Context, userID uuid.UUID) ([]*entity.Bullet, error)
	CompleteBullet(ctx context.Context, id, userID uuid.UUID) (*entity.Bullet, error)
	MigrateBullet(ctx context.Context, id, userID uuid.UUID) (*entity.Bullet, error)
	CancelBullet(ctx context.Context, id, userID uuid.UUID) (*entity.Bullet, error)
}

type JournalService struct {
	journalv1.UnimplementedJournalServiceServer
	usecase Usecase
	logger  *slog.Logger
}

func NewJournalService(usecase Usecase, logger *slog.Logger) *JournalService {
	return &JournalService{
		usecase: usecase,
		logger:  logger,
	}
}

func RegisterJournalService(srv *grpc.Server, uc Usecase, logger *slog.Logger) {
	journalService := NewJournalService(uc, logger)
	journalv1.RegisterJournalServiceServer(srv, journalService)
}
