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
	StatusOpen BulletStatus = "OPEN"
	StatusDone BulletStatus = "DONE"
)

type Bullet struct {
	ID        uuid.UUID
	Title     string
	Date      time.Time
	Status    BulletStatus
	CreatedAt time.Time
	UpdatedAt time.Time
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
