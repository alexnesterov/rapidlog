package handler_test

import (
	"log/slog"
	"testing"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain"
	handler "github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/handler"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/handler/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestJournalServer_MigrateBullet(t *testing.T) {
	bulletID := uuid.New()

	// Usecase возвращает новый open-bullet, созданный при миграции.
	migrated, err := domain.NewBullet(uuid.New(), domain.KindTask, "test")
	require.NoError(t, err)

	protoBullet, err := handler.ToProtoBullet(migrated)
	require.NoError(t, err)

	tests := []struct {
		name      string
		req       *journalv1.MigrateBulletRequest
		want      *journalv1.MigrateBulletResponse
		wantErr   error
		setupMock func(*mocks.MockMigrateBulletUsecase)
	}{
		{
			name: "parse user id error",
			req: &journalv1.MigrateBulletRequest{
				Id:     bulletID.String(),
				UserId: "invalid-uuid",
			},
			want:      nil,
			wantErr:   status.Error(codes.InvalidArgument, "invalid user id"),
			setupMock: func(_ *mocks.MockMigrateBulletUsecase) {},
		},
		{
			name: "parse bullet id error",
			req: &journalv1.MigrateBulletRequest{
				Id:     "invalid-uuid",
				UserId: migrated.UserID.String(),
			},
			want:      nil,
			wantErr:   status.Error(codes.InvalidArgument, "invalid bullet id"),
			setupMock: func(_ *mocks.MockMigrateBulletUsecase) {},
		},
		{
			name: "usecase migrate bullet error",
			req: &journalv1.MigrateBulletRequest{
				Id:     bulletID.String(),
				UserId: migrated.UserID.String(),
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal error"),
			setupMock: func(m *mocks.MockMigrateBulletUsecase) {
				m.EXPECT().
					Execute(
						mock.Anything,
						usecase.MigrateBulletInput{
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
			setupMock: func(m *mocks.MockMigrateBulletUsecase) {
				m.EXPECT().
					Execute(
						mock.Anything,
						usecase.MigrateBulletInput{
							ID:     bulletID,
							UserID: migrated.UserID,
						},
					).
					Return(nil, &domain.ValidationError{Err: assert.AnError}).
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
			setupMock: func(m *mocks.MockMigrateBulletUsecase) {
				m.EXPECT().
					Execute(
						mock.Anything,
						usecase.MigrateBulletInput{
							ID:     bulletID,
							UserID: migrated.UserID,
						},
					).
					Return(&domain.Bullet{}, nil).
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
			setupMock: func(m *mocks.MockMigrateBulletUsecase) {
				m.EXPECT().
					Execute(
						mock.Anything,
						usecase.MigrateBulletInput{
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
			migrateBulletUsecaseMock := mocks.NewMockMigrateBulletUsecase(t)
			tt.setupMock(migrateBulletUsecaseMock)

			s := handler.NewJournalService(handler.Usecases{
				MigrateBullet: migrateBulletUsecaseMock,
			}, slog.Default())

			got, gotErr := s.MigrateBullet(t.Context(), tt.req)
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
