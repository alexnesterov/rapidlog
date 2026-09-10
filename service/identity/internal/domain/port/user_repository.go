package port

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/service/identity/internal/domain/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	Get(ctx context.Context, id uuid.UUID) (*entity.User, error)
}
