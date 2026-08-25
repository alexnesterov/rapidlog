package port

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/google/uuid"
)

type CreateBulletInput struct {
	Type    entity.BulletType
	Content string
}

type BulletService interface {
	CreateBullet(ctx context.Context, input CreateBulletInput) (*entity.Bullet, error)
	ListBullets(ctx context.Context) ([]*entity.Bullet, error)
	CompleteBullet(ctx context.Context, id uuid.UUID) (*entity.Bullet, error)
	MigrateBullet(ctx context.Context, id uuid.UUID) (*entity.Bullet, error)
}
