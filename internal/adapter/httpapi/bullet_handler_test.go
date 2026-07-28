package httpapi

import (
	"bytes"
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
	"github.com/stretchr/testify/require"
)

func TestCreateBullet(t *testing.T) {
	fixedID := uuid.New()

	cases := []struct {
		name       string
		body       string
		setupMock  func(m *mocks.MockBulletService)
		wantStatus int
		wantData   *entity.Bullet
		wantError  *errorBody
	}{
		{
			name: "success",
			body: `{"title": "Заголовок"}`,
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					CreateBullet(port.CreateBulletRequest{Title: "Заголовок"}).
					Return(&entity.Bullet{
						ID:    fixedID,
						Title: "Заголовок",
					}, nil).
					Once()
			},
			wantStatus: http.StatusCreated,
			wantData: &entity.Bullet{
				ID:    fixedID,
				Title: "Заголовок",
			},
		},
		{
			name:       "invalid request body",
			body:       "не json",
			setupMock:  func(m *mocks.MockBulletService) {},
			wantStatus: http.StatusBadRequest,
			wantError: &errorBody{
				Code:    http.StatusBadRequest,
				Message: "invalid request body",
			},
		},
		{
			name: "service error",
			body: `{"title": "Заголовок"}`,
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					CreateBullet(port.CreateBulletRequest{Title: "Заголовок"}).
					Return(nil, errors.New("service error")).
					Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantError: &errorBody{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			},
		},
		{
			name: "validation error",
			body: `{"title": ""}`,
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					CreateBullet(port.CreateBulletRequest{Title: ""}).
					Return(nil, &entity.ValidationError{Err: errors.New("some validation error")}).
					Once()
			},
			wantStatus: http.StatusBadRequest,
			wantError: &errorBody{
				Code:    http.StatusBadRequest,
				Message: "some validation error",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockBulletService := mocks.NewMockBulletService(t)
			tc.setupMock(mockBulletService)
			bulletHandler := NewBulletHandler(mockBulletService)

			req := httptest.NewRequest(http.MethodPost, "/api/bullets", bytes.NewBufferString(tc.body))
			res := httptest.NewRecorder()
			bulletHandler.CreateBullet(res, req)

			require.Equal(t, tc.wantStatus, res.Code)
			assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

			if tc.wantData != nil {
				var got response[entity.Bullet]
				require.NoError(t, json.NewDecoder(res.Body).Decode(&got))
				assert.Equal(t, responseStatus("success"), got.Status)
				assert.Equal(t, *tc.wantData, got.Data)
			}

			if tc.wantError != nil {
				var got response[any]
				require.NoError(t, json.NewDecoder(res.Body).Decode(&got))
				assert.Equal(t, responseStatus("error"), got.Status)
				assert.Equal(t, tc.wantError.Code, got.Error.Code)
				assert.Equal(t, tc.wantError.Message, got.Error.Message)
			}
		})
	}
}

func TestListBullets(t *testing.T) {
	cases := []struct {
		name       string
		setupMock  func(m *mocks.MockBulletService)
		wantStatus int
		wantBody   []*entity.Bullet
	}{
		{
			name: "success",
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().ListBullets().
					Return([]*entity.Bullet{{Title: "Заголовок"}}, nil).
					Once()
			},
			wantStatus: http.StatusOK,
			wantBody:   []*entity.Bullet{{Title: "Заголовок"}},
		},
		{
			name: "empty list",
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().ListBullets().
					Return(nil, nil).
					Once()
			},
			wantStatus: http.StatusOK,
			wantBody:   []*entity.Bullet{},
		},
		{
			name: "service error",
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().ListBullets().
					Return(nil, errors.New("service error")).
					Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockBulletService := mocks.NewMockBulletService(t)
			tc.setupMock(mockBulletService)
			bulletHandler := NewBulletHandler(mockBulletService)

			req := httptest.NewRequest(http.MethodGet, "/api/bullets", nil)
			res := httptest.NewRecorder()
			bulletHandler.ListBullets(res, req)

			require.Equal(t, tc.wantStatus, res.Code)
			assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

			if tc.wantBody != nil {
				var got []*entity.Bullet
				require.NoError(t, json.NewDecoder(res.Body).Decode(&got))
				assert.Equal(t, tc.wantBody, got)
			}
		})
	}
}
