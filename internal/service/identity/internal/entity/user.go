package entity

import (
	"time"

	"github.com/google/uuid"
)

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
