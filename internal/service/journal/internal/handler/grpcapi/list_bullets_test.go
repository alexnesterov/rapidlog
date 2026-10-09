package grpcapi_test

import (
	"context"
	"testing"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain/port/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/handler/grpcapi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestListBullets(t *testing.T) {
	userID := uuid.New()
	task, err := entity.NewBullet(userID, entity.BulletTask, "Test bullet")
	require.NoError(t, err)

	event, err := entity.NewBullet(userID, entity.BulletEvent, "Test event")
	require.NoError(t, err)

	note, err := entity.NewBullet(userID, entity.BulletNote, "Test note")
	require.NoError(t, err)

	assert.Equal(t, entity.SignifierOpen, task.Signifier)

	cases := []struct {
		name      string
		req       *journalv1.ListBulletsRequest
		setupMock func(m *mocks.MockBulletService)
		wantErr   error
		wantData  []*journalv1.Bullet
	}{
		{
			name: "success with single bullet",
			req: &journalv1.ListBulletsRequest{
				UserId: userID.String(),
			},
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					ListBullets(
						mock.Anything,
						userID,
					).
					Return([]*entity.Bullet{task}, nil).
					Once()
			},
			wantErr: nil,
			wantData: []*journalv1.Bullet{{
				Id:        task.ID.String(),
				UserId:    task.UserID.String(),
				Type:      journalv1.BulletType_BULLET_TYPE_TASK,
				Signifier: journalv1.Signifier_SIGNIFIER_OPEN,
				Content:   "Test bullet",
				CreatedAt: timestamppb.New(task.CreatedAt),
				UpdatedAt: timestamppb.New(task.UpdatedAt),
			}},
		},
		{
			name: "success with multiple bullets",
			req: &journalv1.ListBulletsRequest{
				UserId: userID.String(),
			},
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					ListBullets(
						mock.Anything,
						userID,
					).
					Return([]*entity.Bullet{task, event, note}, nil).
					Once()
			},
			wantErr: nil,
			wantData: []*journalv1.Bullet{
				{
					Id:        task.ID.String(),
					UserId:    task.UserID.String(),
					Type:      journalv1.BulletType_BULLET_TYPE_TASK,
					Signifier: journalv1.Signifier_SIGNIFIER_OPEN,
					Content:   "Test bullet",
					CreatedAt: timestamppb.New(task.CreatedAt),
					UpdatedAt: timestamppb.New(task.UpdatedAt),
				},
				{
					Id:        event.ID.String(),
					UserId:    event.UserID.String(),
					Type:      journalv1.BulletType_BULLET_TYPE_EVENT,
					Signifier: journalv1.Signifier_SIGNIFIER_OPEN,
					Content:   "Test event",
					CreatedAt: timestamppb.New(event.CreatedAt),
					UpdatedAt: timestamppb.New(event.UpdatedAt),
				},
				{
					Id:        note.ID.String(),
					UserId:    note.UserID.String(),
					Type:      journalv1.BulletType_BULLET_TYPE_NOTE,
					Signifier: journalv1.Signifier_SIGNIFIER_OPEN,
					Content:   "Test note",
					CreatedAt: timestamppb.New(note.CreatedAt),
					UpdatedAt: timestamppb.New(note.UpdatedAt),
				},
			},
		},
		{
			name: "parse user id error",
			req: &journalv1.ListBulletsRequest{
				UserId: "invalid-uuid",
			},
			setupMock: func(m *mocks.MockBulletService) {},
			wantErr:   status.Error(codes.InvalidArgument, codes.InvalidArgument.String()),
			wantData:  nil,
		},
		{
			name: "usecase list bullets error",
			req: &journalv1.ListBulletsRequest{
				UserId: userID.String(),
			},
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					ListBullets(
						mock.Anything,
						userID,
					).
					Return(nil, assert.AnError).
					Once()
			},
			wantErr:  status.Error(codes.Internal, codes.Internal.String()),
			wantData: nil,
		},
		{
			name: "map bullets to proto error",
			req: &journalv1.ListBulletsRequest{
				UserId: userID.String(),
			},
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					ListBullets(
						mock.Anything,
						userID,
					).
					Return([]*entity.Bullet{{Content: "test"}}, nil).
					Once()
			},
			wantErr:  status.Error(codes.Internal, codes.Internal.String()),
			wantData: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockBulletService := mocks.NewMockBulletService(t)
			tc.setupMock(mockBulletService)

			journalServer := &grpcapi.JournalServer{Usecase: mockBulletService}
			resp, err := journalServer.ListBullets(context.Background(), tc.req)

			if tc.wantErr != nil {
				require.Nil(t, resp)
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Len(t, resp.GetBullets(), len(tc.wantData))
			for i, want := range tc.wantData {
				if !proto.Equal(resp.GetBullets()[i], want) {
					t.Errorf("Bullet %d does not match\nwant: %v\n got: %v", i, want, resp.GetBullets()[i])
				}
			}
		})
	}
}
