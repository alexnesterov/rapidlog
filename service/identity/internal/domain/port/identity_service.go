package port

import (
	"context"

	"github.com/google/uuid"
)

type IdentityService interface {
	ResolveSession(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
}
