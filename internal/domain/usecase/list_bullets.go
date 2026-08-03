package usecase

import "github.com/alexnesterov/rapidlog-api/internal/domain/entity"

func (s *bulletService) ListBullets() ([]*entity.Bullet, error) {
	return s.repo.List()
}
