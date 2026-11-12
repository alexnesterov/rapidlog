package grpcapi_test

import (
	"log/slog"
	"testing"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/adapter/grpcapi"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/adapter/grpcapi/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestJournalServer_ListBullets(t *testing.T) {
	userID := uuid.New()
	task, err := entity.NewBullet(userID, entity.KindTask, "Test bullet")
	require.NoError(t, err)

	event, err := entity.NewBullet(userID, entity.KindEvent, "Test event")
	require.NoError(t, err)

	note, err := entity.NewBullet(userID, entity.KindNote, "Test note")
	require.NoError(t, err)

	protoTask, err := grpcapi.ToProtoBullet(task)
	require.NoError(t, err)

	protoEvent, err := grpcapi.ToProtoBullet(event)
	require.NoError(t, err)

	protoNote, err := grpcapi.ToProtoBullet(note)
	require.NoError(t, err)

	tests := []struct {
		name      string
		req       *journalv1.ListBulletsRequest
		want      []*journalv1.Bullet
		wantErr   error
		setupMock func(m *mocks.MockUsecase)
	}{
		{
			name: "parse user id error",
			req: &journalv1.ListBulletsRequest{
				UserId: "invalid-uuid",
			},
			want:      nil,
			wantErr:   status.Error(codes.InvalidArgument, "invalid user id"),
			setupMock: func(m *mocks.MockUsecase) {},
		},
		{
			name: "usecase list bullets error",
			req: &journalv1.ListBulletsRequest{
				UserId: userID.String(),
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal error"),
			setupMock: func(m *mocks.MockUsecase) {
				m.EXPECT().
					ListBullets(
						mock.Anything,
						userID,
					).
					Return(nil, assert.AnError).
					Once()
			},
		},
		{
			name: "convert bullets to proto error",
			req: &journalv1.ListBulletsRequest{
				UserId: userID.String(),
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal error"),
			setupMock: func(m *mocks.MockUsecase) {
				m.EXPECT().
					ListBullets(
						mock.Anything,
						userID,
					).
					Return([]*entity.Bullet{{}}, nil).
					Once()
			},
		},
		{
			name: "success with single bullet",
			req: &journalv1.ListBulletsRequest{
				UserId: userID.String(),
			},
			want:    []*journalv1.Bullet{protoTask},
			wantErr: nil,
			setupMock: func(m *mocks.MockUsecase) {
				m.EXPECT().
					ListBullets(
						mock.Anything,
						userID,
					).
					Return([]*entity.Bullet{task}, nil).
					Once()
			},
		},
		{
			name: "success with multiple bullets",
			req: &journalv1.ListBulletsRequest{
				UserId: userID.String(),
			},
			want: []*journalv1.Bullet{
				protoTask,
				protoEvent,
				protoNote,
			},
			wantErr: nil,
			setupMock: func(m *mocks.MockUsecase) {
				m.EXPECT().
					ListBullets(
						mock.Anything,
						userID,
					).
					Return([]*entity.Bullet{task, event, note}, nil).
					Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecaseMock := mocks.NewMockUsecase(t)
			tt.setupMock(usecaseMock)

			s := grpcapi.NewJournalService(usecaseMock, slog.Default())

			got, gotErr := s.ListBullets(t.Context(), tt.req)
			if tt.wantErr != nil {
				require.Nil(t, got)
				assert.ErrorIs(t, gotErr, tt.wantErr)
				return
			}

			require.NoError(t, gotErr)
			assert.Len(t, got.GetData(), len(tt.want))
			for i, want := range tt.want {
				if !proto.Equal(got.GetData()[i], want) {
					t.Errorf("Bullet %d does not match\nwant: %v\n got: %v", i, want, got.GetData()[i])
				}
			}
		})
	}
}
