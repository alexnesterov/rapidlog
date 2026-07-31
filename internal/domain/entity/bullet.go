package entity

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

var ErrTitleRequired = errors.New("title is required")
var ErrTitleTooLong = errors.New("title is too long")

type BulletType string

const (
	BulletTask  BulletType = "task"
	BulletEvent BulletType = "event"
	BulletNote  BulletType = "note"
)

type Signifier string

const (
	SignifierOpen      Signifier = "open"
	SignifierCompleted Signifier = "completed"
	SignifierMigrated  Signifier = "migrated"
	SignifierScheduled Signifier = "scheduled"
	SignifierCancelled Signifier = "cancelled"
)

type Bullet struct {
	ID        uuid.UUID  `json:"id"`
	Type      BulletType `json:"type"`
	Signifier Signifier  `json:"signifier"`
	Title     string     `json:"title"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
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
