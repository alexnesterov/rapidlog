package entity

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

var ErrBulletContentRequired = errors.New("bullet content is required")
var ErrBulletContentTooLong = errors.New("bullet content is too long")
var ErrBulletKindInvalid = errors.New("bullet kind must be task, event or note")
var ErrBulletCreatedToday = errors.New("bullet was created today")
var ErrBulletUserIDRequired = errors.New("bullet user id is required")

var ErrBulletNotOpen = errors.New("bullet is not open")
var ErrBulletNotTask = errors.New("bullet is not a task")
var ErrBulletNotOpenTask = errors.New("bullet is not an open task")

var ErrBulletAlreadyCompleted = errors.New("bullet already completed")
var ErrBulletAlreadyCancelled = errors.New("bullet already cancelled")

const MaxBulletContentLength = 200

type Kind string

const (
	KindTask  Kind = "task"
	KindEvent Kind = "event"
	KindNote  Kind = "note"
)

type Signifier string

const (
	SignifierOpen      Signifier = "open"
	SignifierCompleted Signifier = "completed"
	SignifierCancelled Signifier = "cancelled"
)

type Bullet struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Kind      Kind      `json:"kind"`
	Signifier Signifier `json:"signifier"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewBullet(userID uuid.UUID, kind Kind, content string) (*Bullet, error) {
	now := time.Now()

	bullet := &Bullet{
		ID:        uuid.New(),
		Kind:      kind,
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
	switch b.Kind {
	case KindTask, KindEvent, KindNote:
	default:
		return &ValidationError{Err: ErrBulletKindInvalid}
	}

	if b.Content == "" {
		return &ValidationError{Err: ErrBulletContentRequired}
	}

	if utf8.RuneCountInString(b.Content) > MaxBulletContentLength {
		return &ValidationError{Err: ErrBulletContentTooLong}
	}

	if b.UserID == uuid.Nil {
		return &ValidationError{Err: ErrBulletUserIDRequired}
	}

	return nil
}

func (b *Bullet) Complete() error {
	if b.Signifier == SignifierCompleted {
		return &ValidationError{Err: ErrBulletAlreadyCompleted}
	}
	if b.Signifier != SignifierOpen {
		return &ValidationError{Err: ErrBulletNotOpen}
	}

	if b.Kind != KindTask {
		return &ValidationError{Err: ErrBulletNotTask}
	}

	b.Signifier = SignifierCompleted
	b.UpdatedAt = time.Now()

	return nil
}

func (b *Bullet) Migrate() error {
	if b.Kind != KindTask || b.Signifier != SignifierOpen {
		return &ValidationError{Err: ErrBulletNotOpenTask}
	}

	if b.CreatedAt.Format("2006-01-02") == time.Now().Format("2006-01-02") {
		return &ValidationError{Err: ErrBulletCreatedToday}
	}

	now := time.Now()

	b.CreatedAt = now
	b.UpdatedAt = now

	return nil
}

func (b *Bullet) Cancel() error {
	if b.Signifier == SignifierCancelled {
		return &ValidationError{Err: ErrBulletAlreadyCancelled}
	}
	if b.Signifier != SignifierOpen {
		return &ValidationError{Err: ErrBulletNotOpen}
	}

	b.Signifier = SignifierCancelled
	b.UpdatedAt = time.Now()

	return nil
}
