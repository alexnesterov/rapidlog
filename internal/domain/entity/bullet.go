package entity

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

var ErrTitleRequired = errors.New("title is required")
var ErrTitleTooLong = errors.New("title is too long")

type BulletStatus string

const (
	StatusOpened    BulletStatus = "OPENED"
	StatusCompleted BulletStatus = "COMPLETED"
)

type Bullet struct {
	ID        uuid.UUID    `json:"id"`
	Title     string       `json:"title"`
	Status    BulletStatus `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (b *Bullet) Validate() error {
	if b.Title == "" {
		return &ValidationError{Err: ErrTitleRequired}
	}

	if utf8.RuneCountInString(b.Title) > 200 {
		return &ValidationError{Err: ErrTitleTooLong}
	}

	return nil
}
