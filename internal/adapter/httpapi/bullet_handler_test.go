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
	"github.com/stretchr/testify/mock"
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
		wantError  *errorResponse
	}{
		{
			name: "success",
			body: `{"content": "Заголовок"}`,
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					CreateBullet(port.CreateBulletInput{Content: "Заголовок"}).
					Return(&entity.Bullet{
						ID:      fixedID,
						Content: "Заголовок",
					}, nil).
					Once()
			},
			wantStatus: http.StatusCreated,
			wantData: &entity.Bullet{
				ID:      fixedID,
				Content: "Заголовок",
			},
		},
		{
			name:       "invalid request body",
			body:       "не json",
			setupMock:  func(m *mocks.MockBulletService) {},
			wantStatus: http.StatusBadRequest,
			wantError: &errorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid request body",
			},
		},
		{
			name: "service error",
			body: `{"content": "Заголовок"}`,
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					CreateBullet(port.CreateBulletInput{Content: "Заголовок"}).
					Return(nil, errors.New("service error")).
					Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantError: &errorResponse{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			},
		},
		{
			name: "validation error",
			body: `{"content": ""}`,
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					CreateBullet(port.CreateBulletInput{Content: ""}).
					Return(nil, &entity.ValidationError{Err: errors.New("some validation error")}).
					Once()
			},
			wantStatus: http.StatusBadRequest,
			wantError: &errorResponse{
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
				assert.Equal(t, *tc.wantData, got.Data)
			}

			if tc.wantError != nil {
				var got response[any]
				require.NoError(t, json.NewDecoder(res.Body).Decode(&got))
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
		wantData   []bulletDayGroup
		wantError  *errorResponse
	}{
		{
			name: "success",
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().ListBullets().
					Return([]*entity.Bullet{{Content: "Заголовок"}}, nil).
					Once()
			},
			wantStatus: http.StatusOK,
			wantData: []bulletDayGroup{
				{Day: "0001-01-01", Bullets: []*entity.Bullet{{Content: "Заголовок"}}},
			},
		},
		{
			name: "empty list",
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().ListBullets().
					Return(nil, nil).
					Once()
			},
			wantStatus: http.StatusOK,
			wantData:   []bulletDayGroup{},
		},
		{
			name: "service error",
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().ListBullets().
					Return(nil, errors.New("service error")).
					Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantError: &errorResponse{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			},
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

			if tc.wantData != nil {
				var got response[[]bulletDayGroup]
				require.NoError(t, json.NewDecoder(res.Body).Decode(&got))
				assert.Equal(t, tc.wantData, got.Data)
			}

			if tc.wantError != nil {
				var got response[any]
				require.NoError(t, json.NewDecoder(res.Body).Decode(&got))
				assert.Equal(t, tc.wantError.Code, got.Error.Code)
				assert.Equal(t, tc.wantError.Message, got.Error.Message)
			}
		})
	}
}

func TestCompleteBullet(t *testing.T) {
	cases := []struct {
		name       string
		params     struct{ id string }
		setupMock  func(m *mocks.MockBulletService)
		wantStatus int
		wantData   *entity.Bullet
		wantError  *errorResponse
	}{
		{
			name:   "success",
			params: struct{ id string }{id: uuid.New().String()},
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().CompleteBullet(mock.AnythingOfType("uuid.UUID")).
					Return(&entity.Bullet{Signifier: entity.SignifierCompleted}, nil).
					Once()
			},
			wantStatus: http.StatusOK,
			wantData: &entity.Bullet{
				Signifier: entity.SignifierCompleted,
			},
		},
		{
			name:       "parse id error",
			params:     struct{ id string }{id: "123"},
			setupMock:  func(m *mocks.MockBulletService) {},
			wantStatus: http.StatusBadRequest,
			wantError: &errorResponse{
				Code:    http.StatusBadRequest,
				Message: "parse id error",
			},
		},
		{
			name:   "service error",
			params: struct{ id string }{id: uuid.New().String()},
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().CompleteBullet(mock.AnythingOfType("uuid.UUID")).
					Return(nil, errors.New("service error")).
					Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantError: &errorResponse{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			},
		},
		{
			name:   "not found",
			params: struct{ id string }{id: uuid.New().String()},
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().CompleteBullet(mock.AnythingOfType("uuid.UUID")).
					Return(nil, port.ErrNotFound).
					Once()
			},
			wantStatus: http.StatusNotFound,
			wantError: &errorResponse{
				Code:    http.StatusNotFound,
				Message: "bullet not found",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockBulletService := mocks.NewMockBulletService(t)
			tc.setupMock(mockBulletService)
			bulletHandler := NewBulletHandler(mockBulletService)

			req := httptest.NewRequest(http.MethodPost, "/api/bullets/"+tc.params.id+"/complete", nil)
			res := httptest.NewRecorder()

			req.SetPathValue("id", tc.params.id)

			bulletHandler.CompleteBullet(res, req)

			require.Equal(t, tc.wantStatus, res.Code)
			assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

			if tc.wantData != nil {
				var got response[entity.Bullet]
				require.NoError(t, json.NewDecoder(res.Body).Decode(&got))
				assert.Equal(t, *tc.wantData, got.Data)
			}

			if tc.wantError != nil {
				var got response[any]
				require.NoError(t, json.NewDecoder(res.Body).Decode(&got))
				assert.Equal(t, tc.wantError.Code, got.Error.Code)
				assert.Equal(t, tc.wantError.Message, got.Error.Message)
			}
		})
	}
}
