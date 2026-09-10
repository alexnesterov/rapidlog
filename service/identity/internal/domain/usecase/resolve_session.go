package usecase

import (
	"context"

	"github.com/google/uuid"
)

func (s *identityService) ResolveSession(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	user, err := s.userRepo.Get(ctx, id)
	if err == nil {
		return user.ID, nil
	}

	return uuid.Nil, err
}
