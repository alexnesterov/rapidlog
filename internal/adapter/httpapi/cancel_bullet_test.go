package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCancelBulletHandler(t *testing.T) {
	cases := []struct {
		name       string
		id         string
		setupMock  func(*mocks.MockBulletService)
		wantStatus int
		wantData   *entity.Bullet
		wantErr    *errorResponse
	}{
		{
			name: "success",
			id:   uuid.New().String(),
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().CancelBullet(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
					Return(&entity.Bullet{Signifier: entity.SignifierCancelled}, nil).
					Once()
			},
			wantStatus: http.StatusOK,
			wantData: &entity.Bullet{
				Signifier: entity.SignifierCancelled,
			},
		},
		{
			name:       "parse id error",
			id:         "123",
			setupMock:  func(m *mocks.MockBulletService) {},
			wantStatus: http.StatusBadRequest,
			wantErr: &errorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid id",
			},
		},
		{
			name: "usecase error",
			id:   uuid.New().String(),
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().CancelBullet(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
					Return(nil, errors.New("usecase error")).
					Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantErr: &errorResponse{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			},
		},
		{
			name: "not found error",
			id:   uuid.New().String(),
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().CancelBullet(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
					Return(nil, port.ErrNotFound).
					Once()
			},
			wantStatus: http.StatusNotFound,
			wantErr: &errorResponse{
				Code:    http.StatusNotFound,
				Message: "bullet not found",
			},
		},
		{
			name: "validation error",
			id:   uuid.New().String(),
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().CancelBullet(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
					Return(nil, &entity.ValidationError{Err: errors.New("validation error")}).
					Once()
			},
			wantStatus: http.StatusBadRequest,
			wantErr: &errorResponse{
				Code:    http.StatusBadRequest,
				Message: "validation error",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockBulletService := mocks.NewMockBulletService(t)
			tc.setupMock(mockBulletService)
			handler := NewBulletHandler(mockBulletService)

			req := httptest.NewRequest(http.MethodPost, "/api/bullets/"+tc.id+"/cancel", nil)
			req.SetPathValue("id", tc.id)
			resp := httptest.NewRecorder()
			handler.CancelBullet(resp, req)

			require.Equal(t, tc.wantStatus, resp.Code)
			assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))

			if tc.wantData != nil {
				var got response[entity.Bullet]
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
				assert.Equal(t, *tc.wantData, got.Data)
			}

			if tc.wantErr != nil {
				var got response[any]
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
				assert.Equal(t, tc.wantErr, got.Error)
			}
		})
	}
}
