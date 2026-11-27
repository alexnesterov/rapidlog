package handler_test

import (
	"log/slog"
	"testing"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain"
	handler "github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/handler"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/handler/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestJournalServer_CreateBullet(t *testing.T) {
	bullet, err := domain.NewBullet(uuid.New(), domain.KindTask, "test")
	require.NoError(t, err)

	protoBullet, err := handler.ToProtoBullet(bullet)
	require.NoError(t, err)

	tests := []struct {
		name      string
		req       *journalv1.CreateBulletRequest
		want      *journalv1.CreateBulletResponse
		wantErr   error
		setupMock func(*mocks.MockCreateBulletUsecase)
	}{
		{
			name: "parse user id error",
			req: &journalv1.CreateBulletRequest{
				UserId:  "invalid-uuid",
				Kind:    journalv1.Kind_KIND_TASK,
				Content: "test",
			},
			want:      nil,
			wantErr:   status.Error(codes.InvalidArgument, "invalid user id"),
			setupMock: func(_ *mocks.MockCreateBulletUsecase) {},
		},
		{
			name: "parse kind error",
			req: &journalv1.CreateBulletRequest{
				UserId:  protoBullet.UserId,
				Kind:    journalv1.Kind_KIND_UNSPECIFIED,
				Content: "test",
			},
			want:      nil,
			wantErr:   status.Error(codes.InvalidArgument, "invalid kind"),
			setupMock: func(_ *mocks.MockCreateBulletUsecase) {},
		},
		{
			name: "usecase create bullet error",
			req: &journalv1.CreateBulletRequest{
				UserId:  protoBullet.UserId,
				Kind:    protoBullet.Kind,
				Content: protoBullet.Content,
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal error"),
			setupMock: func(m *mocks.MockCreateBulletUsecase) {
				m.EXPECT().
					Execute(
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
				Kind:    protoBullet.Kind,
				Content: protoBullet.Content,
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, assert.AnError.Error()),
			setupMock: func(m *mocks.MockCreateBulletUsecase) {
				m.EXPECT().
					Execute(
						mock.Anything,
						mock.Anything,
					).
					Return(nil, &domain.ValidationError{Err: assert.AnError}).
					Once()
			},
		},
		{
			name: "convert bullet to proto error",
			req: &journalv1.CreateBulletRequest{
				UserId:  protoBullet.UserId,
				Kind:    protoBullet.Kind,
				Content: protoBullet.Content,
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal error"),
			setupMock: func(m *mocks.MockCreateBulletUsecase) {
				m.EXPECT().
					Execute(
						mock.Anything,
						mock.Anything,
					).
					Return(&domain.Bullet{}, nil).
					Once()
			},
		},
		{
			name: "success",
			req: &journalv1.CreateBulletRequest{
				UserId:  protoBullet.UserId,
				Kind:    protoBullet.Kind,
				Content: protoBullet.Content,
			},
			want: &journalv1.CreateBulletResponse{
				Data: protoBullet,
			},
			wantErr: nil,
			setupMock: func(m *mocks.MockCreateBulletUsecase) {
				m.EXPECT().
					Execute(
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
			createBulletUsecaseMock := mocks.NewMockCreateBulletUsecase(t)
			tt.setupMock(createBulletUsecaseMock)

			s := handler.NewJournalService(handler.Usecases{
				CreateBullet: createBulletUsecaseMock,
			}, slog.Default())

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
