package usecase

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain"
	"github.com/google/uuid"
)

type BulletRepository interface {
	Get(ctx context.Context, id, userID uuid.UUID) (*domain.Bullet, error)
	List(ctx context.Context, userID uuid.UUID) ([]*domain.Bullet, error)
	Save(ctx context.Context, bullet *domain.Bullet) error
}
