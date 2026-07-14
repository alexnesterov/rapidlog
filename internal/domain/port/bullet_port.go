// Package port contains domain ports
package port

import (
	"github.com/alexnesterov/dotline/internal/domain/entity"
	"github.com/google/uuid"
)

type BulletRepository interface {
	Create(bullet *entity.Bullet) error
	List() ([]*entity.Bullet, error)
	Read(id uuid.UUID) (*entity.Bullet, error)
	Update(bullet *entity.Bullet) error
	Delete(id uuid.UUID) error
}

type CreateBulletRequest struct {
	Title string
}

type UpdateBulletRequest struct {
	ID     uuid.UUID
	Title  string
	Status entity.Status
}

type BulletService interface {
	CreateBullet(req CreateBulletRequest) (*entity.Bullet, error)
	ListBullets() ([]*entity.Bullet, error)
	ReadBullet(id uuid.UUID) (*entity.Bullet, error)
	UpdateBullet(req UpdateBulletRequest) (*entity.Bullet, error)
	DeleteBullet(id uuid.UUID) error
}
