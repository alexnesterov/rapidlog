package usecase

import (
	"context"
	"errors"

	"github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/entity"
	"github.com/google/uuid"
)

func (s *Usecase) ResolveSession(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	user, err := s.repo.Get(ctx, id)
	if err == nil {
		return user.ID, nil
	}
	if !errors.Is(err, entity.ErrNotFound) {
		return uuid.Nil, err
	}

	newUser := entity.NewUser()
	err = s.repo.Create(ctx, newUser)
	if err != nil {
		return uuid.Nil, err
	}

	return newUser.ID, nil
}
