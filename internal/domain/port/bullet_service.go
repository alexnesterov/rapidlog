package port

import (
	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/google/uuid"
)

type CreateBulletRequest struct {
	Title string
}

type UpdateBulletRequest struct {
	ID     uuid.UUID
	Title  string
	Status entity.BulletStatus
}

type BulletService interface {
	CreateBullet(req CreateBulletRequest) (*entity.Bullet, error)
	ListBullets() ([]*entity.Bullet, error)
	ReadBullet(id uuid.UUID) (*entity.Bullet, error)
	UpdateBullet(req UpdateBulletRequest) (*entity.Bullet, error)
	DeleteBullet(id uuid.UUID) error
}
