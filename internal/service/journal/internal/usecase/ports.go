package usecase

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/google/uuid"
)

type BulletRepository interface {
	Get(ctx context.Context, id, userID uuid.UUID) (*entity.Bullet, error)
	List(ctx context.Context, userID uuid.UUID) ([]*entity.Bullet, error)
	Save(ctx context.Context, bullet *entity.Bullet) error
}
