package usecase

import (
	"context"
	"errors"

	"github.com/alexnesterov/rapidlog-api/service/identity/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/service/identity/internal/domain/port"
	"github.com/google/uuid"
)

func (s *identityService) ResolveSession(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	user, err := s.userRepo.Get(ctx, id)
	if err == nil {
		return user.ID, nil
	}
	if !errors.Is(err, port.ErrNotFound) {
		return uuid.Nil, err
	}

	newUser := entity.NewUser()
	err = s.userRepo.Create(ctx, newUser)
	if err != nil {
		return uuid.Nil, err
	}

	return newUser.ID, nil
}
