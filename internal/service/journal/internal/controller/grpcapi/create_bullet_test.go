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

func TestJournalServer_CreateBullet(t *testing.T) {
	bullet, err := entity.NewBullet(uuid.New(), entity.BulletTask, "test")
	require.NoError(t, err)

	protoBullet, err := grpcapi.ToProtoBullet(bullet)
	require.NoError(t, err)

	tests := []struct {
		name      string
		req       *journalv1.CreateBulletRequest
		want      *journalv1.CreateBulletResponse
		wantErr   error
		setupMock func(*mocks.MockUsecase)
	}{
		{
			name: "parse user id error",
			req: &journalv1.CreateBulletRequest{
				UserId:  "invalid-uuid",
				Type:    journalv1.BulletType_BULLET_TYPE_TASK,
				Content: "test",
			},
			want:      nil,
			wantErr:   status.Error(codes.InvalidArgument, "invalid user id"),
			setupMock: func(_ *mocks.MockUsecase) {},
		},
		{
			name: "usecase create bullet error",
			req: &journalv1.CreateBulletRequest{
				UserId:  protoBullet.UserId,
				Type:    protoBullet.Type,
				Content: protoBullet.Content,
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal error"),
			setupMock: func(m *mocks.MockUsecase) {
				m.EXPECT().
					CreateBullet(
						mock.Anything,
						mock.Anything,
					).
					Return(nil, assert.AnError).
					Once()
			},
		},
		{
			name: "usecase validation error",
			req: &journalv1.CreateBulletRequest{
				UserId:  protoBullet.UserId,
				Type:    protoBullet.Type,
				Content: protoBullet.Content,
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, assert.AnError.Error()),
			setupMock: func(m *mocks.MockUsecase) {
				m.EXPECT().
					CreateBullet(
						mock.Anything,
						mock.Anything,
					).
					Return(nil, &entity.ValidationError{Err: assert.AnError}).
					Once()
			},
		},
		{
			name: "convert bullet to proto error",
			req: &journalv1.CreateBulletRequest{
				UserId:  protoBullet.UserId,
				Type:    protoBullet.Type,
				Content: protoBullet.Content,
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal error"),
			setupMock: func(m *mocks.MockUsecase) {
				m.EXPECT().
					CreateBullet(
						mock.Anything,
						mock.Anything,
					).
					Return(&entity.Bullet{}, nil).
					Once()
			},
		},
		{
			name: "success",
			req: &journalv1.CreateBulletRequest{
				UserId:  protoBullet.UserId,
				Type:    protoBullet.Type,
				Content: protoBullet.Content,
			},
			want: &journalv1.CreateBulletResponse{
				Data: protoBullet,
			},
			wantErr: nil,
			setupMock: func(m *mocks.MockUsecase) {
				m.EXPECT().
					CreateBullet(
						mock.Anything,
						mock.Anything,
					).
					Return(bullet, nil).
					Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecaseMock := mocks.NewMockUsecase(t)
			tt.setupMock(usecaseMock)

			s := grpcapi.NewJournalService(usecaseMock, slog.Default())

			got, gotErr := s.CreateBullet(t.Context(), tt.req)
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
