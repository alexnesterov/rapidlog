package port

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/google/uuid"
)

type BulletRepository interface {
	Create(ctx context.Context, bullet *entity.Bullet) error
	List(ctx context.Context, userID uuid.UUID) ([]*entity.Bullet, error)
	Get(ctx context.Context, id, userID uuid.UUID) (*entity.Bullet, error)
	Update(ctx context.Context, bullet *entity.Bullet) error
}
