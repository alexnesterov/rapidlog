package port

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	Resolve(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
}
