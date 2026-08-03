package port

import (
	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/google/uuid"
)

type CreateBulletInput struct {
	Type    entity.BulletType
	Content string
}

type BulletService interface {
	CreateBullet(input CreateBulletInput) (*entity.Bullet, error)
	ListBullets() ([]*entity.Bullet, error)
	CompleteBullet(id uuid.UUID) (*entity.Bullet, error)
}
