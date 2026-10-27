package grpcapi_test

import (
	"log/slog"
	"testing"
	"time"

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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestJournalServer_CreateBullet(t *testing.T) {
	fixedID := uuid.New()
	fixedUserID := uuid.New()
	fixedTime := time.Now()

	tests := []struct {
		name      string
		setupMock func(*mocks.MockBulletService)
		req       *journalv1.CreateBulletRequest
		want      *journalv1.CreateBulletResponse
		wantErr   error
	}{
		{
			name:      "parse user id error",
			setupMock: func(_ *mocks.MockBulletService) {},
			req: &journalv1.CreateBulletRequest{
				UserId:  "invalid-uuid",
				Type:    journalv1.BulletType_BULLET_TYPE_TASK,
				Content: "test",
			},
			wantErr: status.Error(codes.InvalidArgument, "invalid user id"),
		},
		{
			name: "usecase create bullet error",
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					CreateBullet(
						mock.Anything,
						mock.Anything,
					).
					Return(nil, assert.AnError).
					Once()
			},
			req: &journalv1.CreateBulletRequest{
				UserId:  fixedUserID.String(),
				Type:    journalv1.BulletType_BULLET_TYPE_TASK,
				Content: "test",
			},
			wantErr: status.Error(codes.Internal, "internal error"),
		},
		{
			name: "usecase validation error",
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					CreateBullet(
						mock.Anything,
						mock.Anything,
					).
					Return(nil, &entity.ValidationError{Err: assert.AnError}).
					Once()
			},
			req: &journalv1.CreateBulletRequest{
				UserId:  fixedUserID.String(),
				Type:    journalv1.BulletType_BULLET_TYPE_TASK,
				Content: "test",
			},
			wantErr: status.Error(codes.InvalidArgument, assert.AnError.Error()),
		},
		{
			name: "convert bullet to proto error",
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					CreateBullet(
						mock.Anything,
						mock.Anything,
					).
					Return(&entity.Bullet{
						ID:        fixedID,
						UserID:    fixedUserID,
						Type:      entity.BulletType("test"),
						Signifier: entity.SignifierOpen,
						Content:   "test",
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
					}, nil).
					Once()
			},
			req: &journalv1.CreateBulletRequest{
				UserId:  fixedUserID.String(),
				Type:    journalv1.BulletType_BULLET_TYPE_TASK,
				Content: "test",
			},
			wantErr: status.Error(codes.Internal, "internal error"),
		},
		{
			name: "success",
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					CreateBullet(
						mock.Anything,
						mock.Anything,
					).
					Return(&entity.Bullet{
						ID:        fixedID,
						UserID:    fixedUserID,
						Type:      entity.BulletTask,
						Signifier: entity.SignifierOpen,
						Content:   "test",
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
					}, nil).
					Once()
			},
			req: &journalv1.CreateBulletRequest{
				UserId:  fixedUserID.String(),
				Type:    journalv1.BulletType_BULLET_TYPE_TASK,
				Content: "test",
			},
			want: &journalv1.CreateBulletResponse{
				Data: &journalv1.Bullet{
					Id:        fixedID.String(),
					UserId:    fixedUserID.String(),
					Type:      journalv1.BulletType_BULLET_TYPE_TASK,
					Signifier: journalv1.Signifier_SIGNIFIER_OPEN,
					Content:   "test",
					CreatedAt: timestamppb.New(fixedTime),
					UpdatedAt: timestamppb.New(fixedTime),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBulletService := mocks.NewMockBulletService(t)
			tt.setupMock(mockBulletService)

			var s grpcapi.JournalServer
			s.Usecase = mockBulletService
			s.Logger = slog.Default()

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
