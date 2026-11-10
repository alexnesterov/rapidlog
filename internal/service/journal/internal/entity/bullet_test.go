package entity

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBullet(t *testing.T) {
	userID := uuid.New()

	cases := []struct {
		name    string
		userID  uuid.UUID
		kind    Kind
		content string
		wantErr error
	}{
		{
			name:    "valid task",
			userID:  userID,
			kind:    KindTask,
			content: "Заголовок",
			wantErr: nil,
		},
		{
			name:    "valid event",
			userID:  userID,
			kind:    KindEvent,
			content: "Заголовок",
			wantErr: nil,
		},
		{
			name:    "valid note",
			userID:  userID,
			kind:    KindNote,
			content: "Заголовок",
			wantErr: nil,
		},
		{
			name:    "invalid type",
			userID:  userID,
			kind:    Kind("invalid"),
			content: "Заголовок",
			wantErr: ErrBulletKindInvalid,
		},
		{
			name:    "required content",
			userID:  userID,
			kind:    KindTask,
			content: "",
			wantErr: ErrBulletContentRequired,
		},
		{
			name:    "content too long",
			userID:  userID,
			kind:    KindTask,
			content: strings.Repeat("a", 201),
			wantErr: ErrBulletContentTooLong,
		},
		{
			name:    "user id required",
			userID:  uuid.Nil,
			kind:    KindTask,
			content: "Заголовок",
			wantErr: ErrBulletUserIDRequired,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewBullet(tc.userID, tc.kind, tc.content)

			if tc.wantErr != nil {
				require.Error(t, err)
				var validationErr *ValidationError
				assert.ErrorAs(t, err, &validationErr)
				assert.ErrorIs(t, validationErr.Err, tc.wantErr)
				assert.Equal(t, validationErr.Error(), tc.wantErr.Error())
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)

			assert.Equal(t, tc.userID, got.UserID)
			assert.Equal(t, tc.kind, got.Kind)
			assert.Equal(t, tc.content, got.Content)
			assert.Equal(t, SignifierOpen, got.Signifier)
			assert.NotEqual(t, uuid.Nil, got.ID)
			assert.WithinDuration(t, time.Now(), got.CreatedAt, 1*time.Second)
			assert.WithinDuration(t, time.Now(), got.UpdatedAt, 1*time.Second)
		})
	}
}

func TestBullet_Validate(t *testing.T) {
	cases := []struct {
		name    string
		bullet  Bullet
		wantErr error
	}{
		{
			name: "type invalid",
			bullet: Bullet{
				UserID:  uuid.New(),
				Content: "Заголовок",
				Kind:    "invalid",
			},
			wantErr: ErrBulletKindInvalid,
		},
		{
			name: "type empty",
			bullet: Bullet{
				UserID:  uuid.New(),
				Content: "Заголовок",
				Kind:    "",
			},
			wantErr: ErrBulletKindInvalid,
		},
		{
			name: "valid content",
			bullet: Bullet{
				UserID:  uuid.New(),
				Content: "Заголовок",
				Kind:    KindTask,
			},
			wantErr: nil,
		},
		{
			name: "empty content",
			bullet: Bullet{
				UserID:  uuid.New(),
				Content: "",
				Kind:    KindTask,
			},
			wantErr: ErrBulletContentRequired,
		},
		{
			name: "content too long",
			bullet: Bullet{
				UserID:  uuid.New(),
				Content: strings.Repeat("а", 201),
				Kind:    KindTask,
			},
			wantErr: ErrBulletContentTooLong,
		},
		{
			name: "200 cyrillic chars is valid",
			bullet: Bullet{
				UserID:  uuid.New(),
				Content: strings.Repeat("ф", 200),
				Kind:    KindTask,
			},
			wantErr: nil,
		},
		{
			name: "user id required",
			bullet: Bullet{
				UserID:  uuid.Nil,
				Content: "Заголовок",
				Kind:    KindTask,
			},
			wantErr: ErrBulletUserIDRequired,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.bullet.Validate()

			if tc.wantErr != nil {
				require.Error(t, err)
				var validationErr *ValidationError
				assert.ErrorAs(t, err, &validationErr)
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestBullet_Complete(t *testing.T) {
	cases := []struct {
		name    string
		bullet  Bullet
		wantErr error
	}{
		{
			name: "task open",
			bullet: Bullet{
				Kind:      KindTask,
				Signifier: SignifierOpen,
			},
			wantErr: nil,
		},
		{
			name: "task completed",
			bullet: Bullet{
				Kind:      KindTask,
				Signifier: SignifierCompleted,
			},
			wantErr: ErrBulletAlreadyCompleted,
		},
		{
			name: "task canceled",
			bullet: Bullet{
				Kind:      KindTask,
				Signifier: SignifierCancelled,
			},
			wantErr: ErrBulletNotOpen,
		},
		{
			name: "event open",
			bullet: Bullet{
				Kind:      KindEvent,
				Signifier: SignifierOpen,
			},
			wantErr: ErrBulletNotTask,
		},
		{
			name: "note open",
			bullet: Bullet{
				Kind:      KindNote,
				Signifier: SignifierOpen,
			},
			wantErr: ErrBulletNotTask,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			initialUpdatedAt := tc.bullet.UpdatedAt

			err := tc.bullet.Complete()
			if tc.wantErr != nil {
				require.Error(t, err)
				var validationErr *ValidationError
				assert.ErrorAs(t, err, &validationErr)
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, SignifierCompleted, tc.bullet.Signifier)
			assert.True(t, tc.bullet.UpdatedAt.After(initialUpdatedAt))
		})
	}
}

func TestBullet_Migrate(t *testing.T) {
	cases := []struct {
		name      string
		kind      Kind
		signifier Signifier
		createdAt time.Time
		wantErr   error
	}{
		{
			name:      "task open",
			kind:      KindTask,
			signifier: SignifierOpen,
			createdAt: time.Now().Add(-24 * time.Hour),
			wantErr:   nil,
		},
		{
			name:      "task completed",
			kind:      KindTask,
			signifier: SignifierCompleted,
			createdAt: time.Now().Add(-24 * time.Hour),
			wantErr:   ErrBulletNotOpenTask,
		},
		{
			name:      "task cancelled",
			kind:      KindTask,
			signifier: SignifierCancelled,
			createdAt: time.Now().Add(-24 * time.Hour),
			wantErr:   ErrBulletNotOpenTask,
		},
		{
			name:      "event open",
			kind:      KindEvent,
			signifier: SignifierOpen,
			createdAt: time.Now().Add(-24 * time.Hour),
			wantErr:   ErrBulletNotOpenTask,
		},
		{
			name:      "event canceled",
			kind:      KindEvent,
			signifier: SignifierCancelled,
			createdAt: time.Now().Add(-24 * time.Hour),
			wantErr:   ErrBulletNotOpenTask,
		},
		{
			name:      "note open",
			kind:      KindNote,
			signifier: SignifierOpen,
			createdAt: time.Now().Add(-24 * time.Hour),
			wantErr:   ErrBulletNotOpenTask,
		},
		{
			name:      "note canceled",
			kind:      KindNote,
			signifier: SignifierCancelled,
			createdAt: time.Now().Add(-24 * time.Hour),
			wantErr:   ErrBulletNotOpenTask,
		},
		{
			name:      "created at today",
			kind:      KindTask,
			signifier: SignifierOpen,
			createdAt: time.Now(),
			wantErr:   ErrBulletCreatedToday,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bullet := Bullet{
				ID:        uuid.New(),
				Kind:      tc.kind,
				Signifier: tc.signifier,
				Content:   "Test task",
				UserID:    uuid.New(),
				CreatedAt: tc.createdAt,
				UpdatedAt: tc.createdAt,
			}

			initial := bullet

			gotErr := bullet.Migrate()
			if tc.wantErr != nil {
				require.Error(t, gotErr)
				var validationErr *ValidationError
				assert.ErrorAs(t, gotErr, &validationErr)
				assert.ErrorIs(t, validationErr.Err, tc.wantErr)
				return
			}

			require.NoError(t, gotErr)
			assert.Equal(t, initial.ID, bullet.ID)
			assert.Equal(t, SignifierOpen, bullet.Signifier)
			assert.True(t, bullet.CreatedAt.After(initial.CreatedAt))
			assert.True(t, bullet.UpdatedAt.After(initial.UpdatedAt))
		})
	}
}

func TestBullet_Cancel(t *testing.T) {
	cases := []struct {
		name    string
		bullet  Bullet
		wantErr error
	}{
		{
			name:    "task open",
			bullet:  Bullet{Kind: KindTask, Signifier: SignifierOpen},
			wantErr: nil,
		},
		{
			name:    "event open",
			bullet:  Bullet{Kind: KindEvent, Signifier: SignifierOpen},
			wantErr: nil,
		},
		{
			name:    "note open",
			bullet:  Bullet{Kind: KindNote, Signifier: SignifierOpen},
			wantErr: nil,
		},
		{
			name:    "task cancelled",
			bullet:  Bullet{Kind: KindTask, Signifier: SignifierCancelled},
			wantErr: ErrBulletAlreadyCancelled,
		},
		{
			name:    "task completed",
			bullet:  Bullet{Kind: KindTask, Signifier: SignifierCompleted},
			wantErr: ErrBulletNotOpen,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			initialUpdatedAt := tc.bullet.UpdatedAt

			err := tc.bullet.Cancel()
			if tc.wantErr != nil {
				require.Error(t, err)
				var validationErr *ValidationError
				assert.ErrorAs(t, err, &validationErr)
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, SignifierCancelled, tc.bullet.Signifier)
			assert.True(t, tc.bullet.UpdatedAt.After(initialUpdatedAt))
		})
	}
}
