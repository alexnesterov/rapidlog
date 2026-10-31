package grpcapi_test

import (
	"log/slog"
	"testing"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/controller/grpcapi"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/port/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestJournalServer_CompleteBullet(t *testing.T) {
	bullet, err := entity.NewBullet(uuid.New(), entity.BulletTask, "test")
	require.NoError(t, err)

	err = bullet.Complete()
	require.NoError(t, err)

	protoBullet, err := grpcapi.ToProtoBullet(bullet)
	require.NoError(t, err)

	tests := []struct {
		name      string
		req       *journalv1.CompleteBulletRequest
		want      *journalv1.CompleteBulletResponse
		wantErr   error
		setupMock func(*mocks.MockBulletService)
	}{
		{
			name: "parse user id error",
			req: &journalv1.CompleteBulletRequest{
				Id:     bullet.ID.String(),
				UserId: "invalid-uuid",
			},
			want:      nil,
			wantErr:   status.Error(codes.InvalidArgument, "invalid user id"),
			setupMock: func(_ *mocks.MockBulletService) {},
		},
		{
			name: "parse bullet id error",
			req: &journalv1.CompleteBulletRequest{
				Id:     "invalid-uuid",
				UserId: bullet.UserID.String(),
			},
			want:      nil,
			wantErr:   status.Error(codes.InvalidArgument, "invalid bullet id"),
			setupMock: func(_ *mocks.MockBulletService) {},
		},
		{
			name: "usecase complete bullet error",
			req: &journalv1.CompleteBulletRequest{
				Id:     bullet.ID.String(),
				UserId: bullet.UserID.String(),
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal error"),
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					CompleteBullet(
						mock.Anything,
						bullet.ID,
						bullet.UserID,
					).
					Return(nil, assert.AnError).
					Once()
			},
		},
		{
			name: "usecase validation error",
			req: &journalv1.CompleteBulletRequest{
				Id:     bullet.ID.String(),
				UserId: bullet.UserID.String(),
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, assert.AnError.Error()),
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					CompleteBullet(
						mock.Anything,
						bullet.ID,
						bullet.UserID,
					).
					Return(nil, &entity.ValidationError{Err: assert.AnError}).
					Once()
			},
		},
		{
			name: "convert bullet to proto error",
			req: &journalv1.CompleteBulletRequest{
				Id:     bullet.ID.String(),
				UserId: bullet.UserID.String(),
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal error"),
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					CompleteBullet(
						mock.Anything,
						bullet.ID,
						bullet.UserID,
					).
					Return(&entity.Bullet{}, nil).
					Once()
			},
		},
		{
			name: "success",
			req: &journalv1.CompleteBulletRequest{
				Id:     bullet.ID.String(),
				UserId: bullet.UserID.String(),
			},
			want: &journalv1.CompleteBulletResponse{
				Data: protoBullet,
			},
			wantErr: nil,
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					CompleteBullet(
						mock.Anything,
						bullet.ID,
						bullet.UserID,
					).
					Return(bullet, nil).
					Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bulletServiceMock := mocks.NewMockBulletService(t)
			tt.setupMock(bulletServiceMock)

			s := grpcapi.NewJournalService(bulletServiceMock, slog.Default())

			got, gotErr := s.CompleteBullet(t.Context(), tt.req)
			if tt.wantErr != nil {
				require.Error(t, gotErr)
				assert.Equal(t, tt.wantErr, gotErr)
				return
			}

			require.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
