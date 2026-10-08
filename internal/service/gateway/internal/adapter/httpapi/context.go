package httpapi

import (
	"context"

	"github.com/google/uuid"
)

type ctxKeyUserID struct{}

func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxKeyUserID{}, id)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(ctxKeyUserID{}).(uuid.UUID)
	return id, ok
}
