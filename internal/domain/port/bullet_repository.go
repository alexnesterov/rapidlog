package port

import (
	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/google/uuid"
)

type BulletRepository interface {
	Create(bullet *entity.Bullet) error
	List() ([]*entity.Bullet, error)
	Read(id uuid.UUID) (*entity.Bullet, error)
}
