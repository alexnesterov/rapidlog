package port

import (
	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/google/uuid"
)

type BulletRepository interface {
	List() ([]*entity.Bullet, error)
	Create(bullet *entity.Bullet) error
	Read(id uuid.UUID) (*entity.Bullet, error)
}
