package server_test

import (
	"context"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/platform/server"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserIDContext(t *testing.T) {
	id := uuid.New()
	ctx := server.WithUserID(context.Background(), id)

	got, ok := server.UserIDFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, id, got)
}

func TestUserIDFromContext_Missing(t *testing.T) {
	_, ok := server.UserIDFromContext(context.Background())
	assert.False(t, ok)
}
