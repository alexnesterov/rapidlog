package usecase

import (
	"context"
	"errors"

	"github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/entity"
	"github.com/google/uuid"
)

type ResolveSession struct {
	userRepo UserRepository
}

func NewResolveSession(userRepo UserRepository) *ResolveSession {
	return &ResolveSession{
		userRepo: userRepo,
	}
}

func (uc *ResolveSession) Execute(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	if id != uuid.Nil {
		user, err := uc.userRepo.Get(ctx, id)
		if err == nil {
			return user.ID, nil
		}
		if !errors.Is(err, entity.ErrNotFound) {
			return uuid.Nil, err
		}
	}

	newUser := entity.NewUser()
	if err := uc.userRepo.Save(ctx, newUser); err != nil {
		return uuid.Nil, err
	}

	return newUser.ID, nil
}
