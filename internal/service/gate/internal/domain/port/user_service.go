package port

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	ResolveUser(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
}
