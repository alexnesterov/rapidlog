package domain_test

import (
	"strings"
	"testing"
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBullet(t *testing.T) {
	userID := uuid.New()

	cases := []struct {
		name    string
		userID  uuid.UUID
		kind    domain.Kind
		content string
		wantErr error
	}{
		{
			name:    "valid task",
			userID:  userID,
			kind:    domain.KindTask,
			content: "Заголовок",
			wantErr: nil,
		},
		{
			name:    "valid event",
			userID:  userID,
			kind:    domain.KindEvent,
			content: "Заголовок",
			wantErr: nil,
		},
		{
			name:    "valid note",
			userID:  userID,
			kind:    domain.KindNote,
			content: "Заголовок",
			wantErr: nil,
		},
		{
			name:    "invalid type",
			userID:  userID,
			kind:    domain.Kind("invalid"),
			content: "Заголовок",
			wantErr: domain.ErrBulletKindInvalid,
		},
		{
			name:    "required content",
			userID:  userID,
			kind:    domain.KindTask,
			content: "",
			wantErr: domain.ErrBulletContentRequired,
		},
		{
			name:    "content too long",
			userID:  userID,
			kind:    domain.KindTask,
			content: strings.Repeat("a", 201),
			wantErr: domain.ErrBulletContentTooLong,
		},
		{
			name:    "user id required",
			userID:  uuid.Nil,
			kind:    domain.KindTask,
			content: "Заголовок",
			wantErr: domain.ErrBulletUserIDRequired,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := domain.NewBullet(tc.userID, tc.kind, tc.content)

			if tc.wantErr != nil {
				require.Error(t, err)
				var validationErr *domain.ValidationError
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
			assert.Equal(t, domain.SignifierOpen, got.Signifier)
			assert.NotEqual(t, uuid.Nil, got.ID)
			assert.WithinDuration(t, time.Now(), got.CreatedAt, 1*time.Second)
			assert.WithinDuration(t, time.Now(), got.UpdatedAt, 1*time.Second)
		})
	}
}

func TestBullet_Validate(t *testing.T) {
	cases := []struct {
		name    string
		bullet  domain.Bullet
		wantErr error
	}{
		{
			name: "kind invalid",
			bullet: domain.Bullet{
				UserID:    uuid.New(),
				Content:   "Заголовок",
				Kind:      "invalid",
				Signifier: domain.SignifierOpen,
			},
			wantErr: domain.ErrBulletKindInvalid,
		},
		{
			name: "kind empty",
			bullet: domain.Bullet{
				UserID:    uuid.New(),
				Content:   "Заголовок",
				Kind:      "",
				Signifier: domain.SignifierOpen,
			},
			wantErr: domain.ErrBulletKindInvalid,
		},
		{
			name: "signifier invalid",
			bullet: domain.Bullet{
				UserID:    uuid.New(),
				Content:   "Заголовок",
				Kind:      domain.KindTask,
				Signifier: "close",
			},
			wantErr: domain.ErrBulletSignifierInvalid,
		},
		{
			name: "signifier empty",
			bullet: domain.Bullet{
				UserID:    uuid.New(),
				Content:   "Заголовок",
				Kind:      domain.KindTask,
				Signifier: "",
			},
			wantErr: domain.ErrBulletSignifierInvalid,
		},
		{
			name: "valid content",
			bullet: domain.Bullet{
				UserID:    uuid.New(),
				Content:   "Заголовок",
				Kind:      domain.KindTask,
				Signifier: domain.SignifierOpen,
			},
			wantErr: nil,
		},
		{
			name: "empty content",
			bullet: domain.Bullet{
				UserID:    uuid.New(),
				Content:   "",
				Kind:      domain.KindTask,
				Signifier: domain.SignifierOpen,
			},
			wantErr: domain.ErrBulletContentRequired,
		},
		{
			name: "content too long",
			bullet: domain.Bullet{
				UserID:    uuid.New(),
				Content:   strings.Repeat("а", 201),
				Kind:      domain.KindTask,
				Signifier: domain.SignifierOpen,
			},
			wantErr: domain.ErrBulletContentTooLong,
		},
		{
			name: "200 cyrillic chars is valid",
			bullet: domain.Bullet{
				UserID:    uuid.New(),
				Content:   strings.Repeat("ф", 200),
				Kind:      domain.KindTask,
				Signifier: domain.SignifierOpen,
			},
			wantErr: nil,
		},
		{
			name: "user id required",
			bullet: domain.Bullet{
				UserID:    uuid.Nil,
				Content:   "Заголовок",
				Kind:      domain.KindTask,
				Signifier: domain.SignifierOpen,
			},
			wantErr: domain.ErrBulletUserIDRequired,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.bullet.Validate()

			if tc.wantErr != nil {
				require.Error(t, err)
				var validationErr *domain.ValidationError
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
		bullet  domain.Bullet
		wantErr error
	}{
		{
			name: "task open",
			bullet: domain.Bullet{
				Kind:      domain.KindTask,
				Signifier: domain.SignifierOpen,
			},
			wantErr: nil,
		},
		{
			name: "task completed",
			bullet: domain.Bullet{
				Kind:      domain.KindTask,
				Signifier: domain.SignifierCompleted,
			},
			wantErr: domain.ErrBulletAlreadyCompleted,
		},
		{
			name: "task canceled",
			bullet: domain.Bullet{
				Kind:      domain.KindTask,
				Signifier: domain.SignifierCancelled,
			},
			wantErr: domain.ErrBulletNotOpen,
		},
		{
			name: "event open",
			bullet: domain.Bullet{
				Kind:      domain.KindEvent,
				Signifier: domain.SignifierOpen,
			},
			wantErr: domain.ErrBulletNotTask,
		},
		{
			name: "note open",
			bullet: domain.Bullet{
				Kind:      domain.KindNote,
				Signifier: domain.SignifierOpen,
			},
			wantErr: domain.ErrBulletNotTask,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			initialUpdatedAt := tc.bullet.UpdatedAt

			err := tc.bullet.Complete()
			if tc.wantErr != nil {
				require.Error(t, err)
				var validationErr *domain.ValidationError
				assert.ErrorAs(t, err, &validationErr)
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, domain.SignifierCompleted, tc.bullet.Signifier)
			assert.True(t, tc.bullet.UpdatedAt.After(initialUpdatedAt))
		})
	}
}

func TestBullet_Migrate(t *testing.T) {
	cases := []struct {
		name      string
		kind      domain.Kind
		signifier domain.Signifier
		createdAt time.Time
		wantErr   error
	}{
		{
			name:      "task open",
			kind:      domain.KindTask,
			signifier: domain.SignifierOpen,
			createdAt: time.Now().Add(-24 * time.Hour),
			wantErr:   nil,
		},
		{
			name:      "task completed",
			kind:      domain.KindTask,
			signifier: domain.SignifierCompleted,
			createdAt: time.Now().Add(-24 * time.Hour),
			wantErr:   domain.ErrBulletNotOpenTask,
		},
		{
			name:      "task cancelled",
			kind:      domain.KindTask,
			signifier: domain.SignifierCancelled,
			createdAt: time.Now().Add(-24 * time.Hour),
			wantErr:   domain.ErrBulletNotOpenTask,
		},
		{
			name:      "event open",
			kind:      domain.KindEvent,
			signifier: domain.SignifierOpen,
			createdAt: time.Now().Add(-24 * time.Hour),
			wantErr:   domain.ErrBulletNotOpenTask,
		},
		{
			name:      "event canceled",
			kind:      domain.KindEvent,
			signifier: domain.SignifierCancelled,
			createdAt: time.Now().Add(-24 * time.Hour),
			wantErr:   domain.ErrBulletNotOpenTask,
		},
		{
			name:      "note open",
			kind:      domain.KindNote,
			signifier: domain.SignifierOpen,
			createdAt: time.Now().Add(-24 * time.Hour),
			wantErr:   domain.ErrBulletNotOpenTask,
		},
		{
			name:      "note canceled",
			kind:      domain.KindNote,
			signifier: domain.SignifierCancelled,
			createdAt: time.Now().Add(-24 * time.Hour),
			wantErr:   domain.ErrBulletNotOpenTask,
		},
		{
			name:      "created at today",
			kind:      domain.KindTask,
			signifier: domain.SignifierOpen,
			createdAt: time.Now(),
			wantErr:   domain.ErrBulletCreatedToday,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bullet := domain.Bullet{
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
				var validationErr *domain.ValidationError
				assert.ErrorAs(t, gotErr, &validationErr)
				assert.ErrorIs(t, validationErr.Err, tc.wantErr)
				return
			}

			require.NoError(t, gotErr)
			assert.Equal(t, initial.ID, bullet.ID)
			assert.Equal(t, domain.SignifierOpen, bullet.Signifier)
			assert.True(t, bullet.CreatedAt.After(initial.CreatedAt))
			assert.True(t, bullet.UpdatedAt.After(initial.UpdatedAt))
		})
	}
}

func TestBullet_Cancel(t *testing.T) {
	cases := []struct {
		name    string
		bullet  domain.Bullet
		wantErr error
	}{
		{
			name:    "task open",
			bullet:  domain.Bullet{Kind: domain.KindTask, Signifier: domain.SignifierOpen},
			wantErr: nil,
		},
		{
			name:    "event open",
			bullet:  domain.Bullet{Kind: domain.KindEvent, Signifier: domain.SignifierOpen},
			wantErr: nil,
		},
		{
			name:    "note open",
			bullet:  domain.Bullet{Kind: domain.KindNote, Signifier: domain.SignifierOpen},
			wantErr: nil,
		},
		{
			name:    "task cancelled",
			bullet:  domain.Bullet{Kind: domain.KindTask, Signifier: domain.SignifierCancelled},
			wantErr: domain.ErrBulletAlreadyCancelled,
		},
		{
			name:    "task completed",
			bullet:  domain.Bullet{Kind: domain.KindTask, Signifier: domain.SignifierCompleted},
			wantErr: domain.ErrBulletNotOpen,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			initialUpdatedAt := tc.bullet.UpdatedAt

			err := tc.bullet.Cancel()
			if tc.wantErr != nil {
				require.Error(t, err)
				var validationErr *domain.ValidationError
				assert.ErrorAs(t, err, &validationErr)
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, domain.SignifierCancelled, tc.bullet.Signifier)
			assert.True(t, tc.bullet.UpdatedAt.After(initialUpdatedAt))
		})
	}
}
