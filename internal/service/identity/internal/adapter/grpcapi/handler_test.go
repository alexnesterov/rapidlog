package grpcapi_test

import (
	"context"
	"errors"
	"testing"

	identityv1 "github.com/alexnesterov/rapidlog-api/gen/identity/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/adapter/grpcapi"
	"github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/domain/port/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResolveSession(t *testing.T) {
	fixedSessionID := uuid.New()
	errUsecase := errors.New("usecase error")

	cases := []struct {
		name      string
		setupMock func(m *mocks.MockIdentityService)
		wantErr   error
	}{
		{
			name: "usecase error",
			setupMock: func(m *mocks.MockIdentityService) {
				m.EXPECT().
					ResolveSession(mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(uuid.Nil, errUsecase).
					Once()
			},
			wantErr: errUsecase,
		},
		{
			name: "success",
			setupMock: func(m *mocks.MockIdentityService) {
				m.EXPECT().
					ResolveSession(mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(fixedSessionID, nil).
					Once()
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockIdentityService := mocks.NewMockIdentityService(t)
			tc.setupMock(mockIdentityService)

			handler := grpcapi.NewHandler(mockIdentityService)
			resp, err := handler.ResolveSession(context.Background(), &identityv1.ResolveSessionRequest{SessionId: fixedSessionID.String()})

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, fixedSessionID.String(), resp.GetUserId())
		})
	}
}
