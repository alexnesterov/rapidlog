package entity

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

var ErrContentRequired = errors.New("content is required")
var ErrContentTooLong = errors.New("content is too long")
var ErrTypeInvalid = errors.New("type must be task, event or note")
var ErrNotOpenTask = errors.New("bullet is not an open task")
var ErrTodayTask = errors.New("bullet is already scheduled for today")
var ErrUserIDRequired = errors.New("user id is required")

var ErrBulletAlreadyCancelled = errors.New("bullet already cancelled")
var ErrBulletNotOpen = errors.New("bullet must be open to be cancelled")

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
	Content   string     `json:"content"`
	UserID    uuid.UUID  `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewBullet(userID uuid.UUID, bulletType BulletType, content string) (*Bullet, error) {
	now := time.Now()

	bullet := &Bullet{
		ID:        uuid.New(),
		Type:      bulletType,
		Signifier: SignifierOpen,
		Content:   content,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := bullet.Validate(); err != nil {
		return nil, err
	}

	return bullet, nil
}

func (b *Bullet) Validate() error {
	switch b.Type {
	case BulletTask, BulletEvent, BulletNote:
	default:
		return &ValidationError{Err: ErrTypeInvalid}
	}

	if b.Content == "" {
		return &ValidationError{Err: ErrContentRequired}
	}

	if utf8.RuneCountInString(b.Content) > 200 {
		return &ValidationError{Err: ErrContentTooLong}
	}

	if b.UserID == uuid.Nil {
		return &ValidationError{Err: ErrUserIDRequired}
	}

	return nil
}

func (b *Bullet) Migrate() (*Bullet, error) {
	if b.Type != BulletTask || b.Signifier != SignifierOpen {
		return nil, &ValidationError{Err: ErrNotOpenTask}
	}

	if b.CreatedAt.Format("2006-01-02") == time.Now().Format("2006-01-02") {
		return nil, &ValidationError{Err: ErrTodayTask}
	}

	now := time.Now()
	b.Signifier = SignifierMigrated
	b.UpdatedAt = now

	migrated := &Bullet{
		ID:        uuid.New(),
		Type:      b.Type,
		Signifier: SignifierOpen,
		Content:   b.Content,
		UserID:    b.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return migrated, nil
}

func (b *Bullet) Cancel() error {
	switch b.Signifier {
	case SignifierCancelled:
		return &ValidationError{Err: ErrBulletAlreadyCancelled}
	case SignifierOpen:
	default:
		return &ValidationError{Err: ErrBulletNotOpen}
	}

	b.Signifier = SignifierCancelled
	b.UpdatedAt = time.Now()

	return nil
}
