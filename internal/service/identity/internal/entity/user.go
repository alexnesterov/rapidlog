package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("not found")

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
}

func NewUser() *User {
	return &User{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
	}
}
