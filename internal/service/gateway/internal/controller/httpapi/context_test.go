package httpapi

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserIDContext(t *testing.T) {
	id := uuid.New()
	ctx := WithUserID(context.Background(), id)

	got, ok := UserIDFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, id, got)
}

func TestUserIDFromContext_Missing(t *testing.T) {
	_, ok := UserIDFromContext(context.Background())
	assert.False(t, ok)
}
