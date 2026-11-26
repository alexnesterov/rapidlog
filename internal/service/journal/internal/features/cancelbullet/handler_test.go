package cancelbullet_test

import (
	"log/slog"
	"testing"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/features/cancelbullet"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/features/cancelbullet/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/protomap"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandler(t *testing.T) {
	bullet, err := entity.NewBullet(uuid.New(), entity.KindTask, "test")
	require.NoError(t, err)

	err = bullet.Cancel()
	require.NoError(t, err)

	protoBullet, err := protomap.ToProtoBullet(bullet)
	require.NoError(t, err)

	tests := []struct {
		name      string
		req       *journalv1.CancelBulletRequest
		want      *journalv1.CancelBulletResponse
		wantErr   error
		setupMock func(*mocks.MockExecutor)
	}{
		{
			name: "parse user id error",
			req: &journalv1.CancelBulletRequest{
				Id:     bullet.ID.String(),
				UserId: "invalid-uuid",
			},
			want:      nil,
			wantErr:   status.Error(codes.InvalidArgument, "invalid user id"),
			setupMock: func(_ *mocks.MockExecutor) {},
		},
		{
			name: "parse bullet id error",
			req: &journalv1.CancelBulletRequest{
				Id:     "invalid-uuid",
				UserId: bullet.UserID.String(),
			},
			want:      nil,
			wantErr:   status.Error(codes.InvalidArgument, "invalid bullet id"),
			setupMock: func(_ *mocks.MockExecutor) {},
		},
		{
			name: "usecase cancel bullet error",
			req: &journalv1.CancelBulletRequest{
				Id:     bullet.ID.String(),
				UserId: bullet.UserID.String(),
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal error"),
			setupMock: func(m *mocks.MockExecutor) {
				m.EXPECT().
					Execute(
						mock.Anything,
						cancelbullet.Input{
							ID:     bullet.ID,
							UserID: bullet.UserID,
						},
					).
					Return(nil, assert.AnError).
					Once()
			},
		},
		{
			name: "usecase validation error",
			req: &journalv1.CancelBulletRequest{
				Id:     bullet.ID.String(),
				UserId: bullet.UserID.String(),
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, assert.AnError.Error()),
			setupMock: func(m *mocks.MockExecutor) {
				m.EXPECT().
					Execute(
						mock.Anything,
						cancelbullet.Input{
							ID:     bullet.ID,
							UserID: bullet.UserID,
						},
					).
					Return(nil, &entity.ValidationError{Err: assert.AnError}).
					Once()
			},
		},
		{
			name: "convert bullet to proto error",
			req: &journalv1.CancelBulletRequest{
				Id:     bullet.ID.String(),
				UserId: bullet.UserID.String(),
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal error"),
			setupMock: func(m *mocks.MockExecutor) {
				m.EXPECT().
					Execute(
						mock.Anything,
						cancelbullet.Input{
							ID:     bullet.ID,
							UserID: bullet.UserID,
						},
					).
					Return(&entity.Bullet{}, nil).
					Once()
			},
		},
		{
			name: "success",
			req: &journalv1.CancelBulletRequest{
				Id:     bullet.ID.String(),
				UserId: bullet.UserID.String(),
			},
			want: &journalv1.CancelBulletResponse{
				Data: protoBullet,
			},
			wantErr: nil,
			setupMock: func(m *mocks.MockExecutor) {
				m.EXPECT().
					Execute(
						mock.Anything,
						cancelbullet.Input{
							ID:     bullet.ID,
							UserID: bullet.UserID,
						},
					).
					Return(bullet, nil).
					Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cancelBulletUsecaseMock := mocks.NewMockExecutor(t)
			tt.setupMock(cancelBulletUsecaseMock)

			h := cancelbullet.NewHandler(cancelBulletUsecaseMock, slog.Default())

			got, gotErr := h.Handle(t.Context(), tt.req)
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
