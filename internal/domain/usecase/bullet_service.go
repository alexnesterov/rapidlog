package usecase

import (
	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/google/uuid"
)

type BulletService struct {
	BulletRepo port.BulletRepository
}

func (s *BulletService) CreateBullet(req port.CreateBulletRequest) (*entity.Bullet, error) {
	bullet := &entity.Bullet{
		ID:    uuid.New(),
		Title: req.Title,
	}

	if err := bullet.Validate(); err != nil {
		return nil, err
	}

	if err := s.BulletRepo.Create(bullet); err != nil {
		return nil, err
	}

	return bullet, nil
}

func (s *BulletService) ListBullets() ([]*entity.Bullet, error) {
	return nil, nil
}

func (s *BulletService) ReadBullet(id uuid.UUID) (*entity.Bullet, error) {
	return nil, nil
}

func (s *BulletService) UpdateBullet(req port.UpdateBulletRequest) (*entity.Bullet, error) {
	return nil, nil
}

func (s *BulletService) DeleteBullet(id uuid.UUID) error {
	return nil
}

var _ port.BulletService = &BulletService{}
