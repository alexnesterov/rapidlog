package migratebullet_test

import (
	"log/slog"
	"testing"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/features/migratebullet"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/features/migratebullet/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/protomap"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandler(t *testing.T) {
	bulletID := uuid.New()

	// Usecase возвращает новый open-bullet, созданный при миграции.
	migrated, err := entity.NewBullet(uuid.New(), entity.KindTask, "test")
	require.NoError(t, err)

	protoBullet, err := protomap.ToProtoBullet(migrated)
	require.NoError(t, err)

	tests := []struct {
		name      string
		req       *journalv1.MigrateBulletRequest
		want      *journalv1.MigrateBulletResponse
		wantErr   error
		setupMock func(*mocks.MockExecutor)
	}{
		{
			name: "parse user id error",
			req: &journalv1.MigrateBulletRequest{
				Id:     bulletID.String(),
				UserId: "invalid-uuid",
			},
			want:      nil,
			wantErr:   status.Error(codes.InvalidArgument, "invalid user id"),
			setupMock: func(_ *mocks.MockExecutor) {},
		},
		{
			name: "parse bullet id error",
			req: &journalv1.MigrateBulletRequest{
				Id:     "invalid-uuid",
				UserId: migrated.UserID.String(),
			},
			want:      nil,
			wantErr:   status.Error(codes.InvalidArgument, "invalid bullet id"),
			setupMock: func(_ *mocks.MockExecutor) {},
		},
		{
			name: "usecase migrate bullet error",
			req: &journalv1.MigrateBulletRequest{
				Id:     bulletID.String(),
				UserId: migrated.UserID.String(),
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal error"),
			setupMock: func(m *mocks.MockExecutor) {
				m.EXPECT().
					Execute(
						mock.Anything,
						migratebullet.Input{
							ID:     bulletID,
							UserID: migrated.UserID,
						},
					).
					Return(nil, assert.AnError).
					Once()
			},
		},
		{
			name: "usecase validation error",
			req: &journalv1.MigrateBulletRequest{
				Id:     bulletID.String(),
				UserId: migrated.UserID.String(),
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, assert.AnError.Error()),
			setupMock: func(m *mocks.MockExecutor) {
				m.EXPECT().
					Execute(
						mock.Anything,
						migratebullet.Input{
							ID:     bulletID,
							UserID: migrated.UserID,
						},
					).
					Return(nil, &entity.ValidationError{Err: assert.AnError}).
					Once()
			},
		},
		{
			name: "convert bullet to proto error",
			req: &journalv1.MigrateBulletRequest{
				Id:     bulletID.String(),
				UserId: migrated.UserID.String(),
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal error"),
			setupMock: func(m *mocks.MockExecutor) {
				m.EXPECT().
					Execute(
						mock.Anything,
						migratebullet.Input{
							ID:     bulletID,
							UserID: migrated.UserID,
						},
					).
					Return(&entity.Bullet{}, nil).
					Once()
			},
		},
		{
			name: "success",
			req: &journalv1.MigrateBulletRequest{
				Id:     bulletID.String(),
				UserId: migrated.UserID.String(),
			},
			want: &journalv1.MigrateBulletResponse{
				Data: protoBullet,
			},
			wantErr: nil,
			setupMock: func(m *mocks.MockExecutor) {
				m.EXPECT().
					Execute(
						mock.Anything,
						migratebullet.Input{
							ID:     bulletID,
							UserID: migrated.UserID,
						},
					).
					Return(migrated, nil).
					Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			migrateBulletUsecaseMock := mocks.NewMockExecutor(t)
			tt.setupMock(migrateBulletUsecaseMock)

			h := migratebullet.NewHandler(migrateBulletUsecaseMock, slog.Default())

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
