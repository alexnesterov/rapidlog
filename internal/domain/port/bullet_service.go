package port

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/google/uuid"
)

type CreateBulletInput struct {
	Type    entity.BulletType
	Content string
	UserID  uuid.UUID `json:"-"`
}

type BulletService interface {
	CreateBullet(ctx context.Context, input CreateBulletInput) (*entity.Bullet, error)
	ListBullets(ctx context.Context, userID uuid.UUID) ([]*entity.Bullet, error)
	CompleteBullet(ctx context.Context, id, userID uuid.UUID) (*entity.Bullet, error)
	MigrateBullet(ctx context.Context, id, userID uuid.UUID) (*entity.Bullet, error)
}
