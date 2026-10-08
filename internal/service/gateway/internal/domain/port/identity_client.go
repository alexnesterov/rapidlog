package port

import (
	"context"

	"github.com/google/uuid"
)

type IdentityClient interface {
	ResolveSession(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
}
