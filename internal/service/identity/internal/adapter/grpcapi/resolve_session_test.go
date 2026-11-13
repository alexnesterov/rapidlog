package grpcapi_test

import (
	"errors"
	"log/slog"
	"testing"

	identityv1 "github.com/alexnesterov/rapidlog-api/api/gen/identity/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/adapter/grpcapi"
	"github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/adapter/grpcapi/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestResolveSession(t *testing.T) {
	fixedSessionID := uuid.New()

	tests := []struct {
		name      string
		req       *identityv1.ResolveSessionRequest
		want      *identityv1.ResolveSessionResponse
		wantErr   error
		setupMock func(m *mocks.MockResolveSessionUsecase)
	}{
		{
			name: "parse session id error",
			req: &identityv1.ResolveSessionRequest{
				SessionId: "invalid-uuid",
			},
			want:      nil,
			wantErr:   status.Error(codes.InvalidArgument, "invalid session id"),
			setupMock: func(_ *mocks.MockResolveSessionUsecase) {},
		},
		{
			name: "usecase error",
			req: &identityv1.ResolveSessionRequest{
				SessionId: fixedSessionID.String(),
			},
			want:    nil,
			wantErr: status.Error(codes.Internal, "internal error"),
			setupMock: func(m *mocks.MockResolveSessionUsecase) {
				m.EXPECT().
					Execute(mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(uuid.Nil, errors.New("usecase error")).
					Once()
			},
		},
		{
			name: "success",
			req: &identityv1.ResolveSessionRequest{
				SessionId: fixedSessionID.String(),
			},
			want: &identityv1.ResolveSessionResponse{
				UserId: fixedSessionID.String(),
			},
			wantErr: nil,
			setupMock: func(m *mocks.MockResolveSessionUsecase) {
				m.EXPECT().
					Execute(mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(fixedSessionID, nil).
					Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolveSessionUsecasMock := mocks.NewMockResolveSessionUsecase(t)
			tt.setupMock(resolveSessionUsecasMock)

			s := grpcapi.NewIdentityService(grpcapi.Usecase{
				ResolveSession: resolveSessionUsecasMock,
			}, slog.Default())

			got, gotErr := s.ResolveSession(t.Context(), tt.req)
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
