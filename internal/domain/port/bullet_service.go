package port

import (
	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
)

type CreateBulletRequest struct {
	Title string
}

type BulletService interface {
	ListBullets() ([]*entity.Bullet, error)
	CreateBullet(req CreateBulletRequest) (*entity.Bullet, error)
}
