package usecase

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
}
