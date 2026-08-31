package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	before := time.Now()

	got := NewUser()

	assert.NotEqual(t, uuid.Nil, got.ID)
	assert.WithinDuration(t, before, got.CreatedAt, time.Second)
}
