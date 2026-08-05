package port

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/google/uuid"
)

type BulletRepository interface {
	Create(ctx context.Context, bullet *entity.Bullet) error
	List(ctx context.Context) ([]*entity.Bullet, error)
	Get(ctx context.Context, id uuid.UUID) (*entity.Bullet, error)
	Update(ctx context.Context, bullet *entity.Bullet) error
}
