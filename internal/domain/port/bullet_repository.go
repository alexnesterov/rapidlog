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
