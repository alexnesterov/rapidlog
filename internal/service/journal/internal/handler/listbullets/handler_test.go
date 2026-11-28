package listbullets_test

import (
	"log/slog"
	"testing"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/handler/listbullets"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/handler/listbullets/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/platform/server"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase"
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
	task, err := domain.NewBullet(userID, domain.KindTask, "Test bullet")
	require.NoError(t, err)

	event, err := domain.NewBullet(userID, domain.KindEvent, "Test event")
	require.NoError(t, err)

	note, err := domain.NewBullet(userID, domain.KindNote, "Test note")
	require.NoError(t, err)

	protoTask, err := server.ToProtoBullet(task)
	require.NoError(t, err)

	protoEvent, err := server.ToProtoBullet(event)
	require.NoError(t, err)

	protoNote, err := server.ToProtoBullet(note)
	require.NoError(t, err)

	tests := []struct {
		name      string
		req       *journalv1.ListBulletsRequest
		want      []*journalv1.Bullet
		wantErr   error
		setupMock func(m *mocks.MockListBulletsUsecase)
	}{
		{
			name: "parse user id error",
			req: &journalv1.ListBulletsRequest{
				UserId: "invalid-uuid",
			},
			want:      nil,
			wantErr:   status.Error(codes.InvalidArgument, "invalid user id"),
			setupMock: func(m *mocks.MockListBulletsUsecase) {},
		},
		{
			name: "usecase list bullets error",
			req: &journalv1.ListBulletsRequest{
				UserId: userID.String(),
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal error"),
			setupMock: func(m *mocks.MockListBulletsUsecase) {
				m.EXPECT().
					Execute(
						mock.Anything,
						usecase.ListBulletsInput{UserID: userID},
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
			setupMock: func(m *mocks.MockListBulletsUsecase) {
				m.EXPECT().
					Execute(
						mock.Anything,
						usecase.ListBulletsInput{UserID: userID},
					).
					Return([]*domain.Bullet{{}}, nil).
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
			setupMock: func(m *mocks.MockListBulletsUsecase) {
				m.EXPECT().
					Execute(
						mock.Anything,
						usecase.ListBulletsInput{UserID: userID},
					).
					Return([]*domain.Bullet{task}, nil).
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
			setupMock: func(m *mocks.MockListBulletsUsecase) {
				m.EXPECT().
					Execute(
						mock.Anything,
						usecase.ListBulletsInput{UserID: userID},
					).
					Return([]*domain.Bullet{task, event, note}, nil).
					Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listBulletsUsecaseMock := mocks.NewMockListBulletsUsecase(t)
			tt.setupMock(listBulletsUsecaseMock)

			h := listbullets.NewHandler(listBulletsUsecaseMock, slog.Default())
			s := server.NewJournalService(server.Handlers{
				ListBullets: h,
			})

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
