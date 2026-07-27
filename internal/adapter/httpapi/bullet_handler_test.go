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
		wantBody   *entity.Bullet
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
			wantBody: &entity.Bullet{
				ID:    fixedID,
				Title: "Заголовок",
			},
		},
		{
			name:       "invalid json body",
			body:       "не json",
			setupMock:  func(m *mocks.MockBulletService) {},
			wantStatus: http.StatusBadRequest,
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
		},
		{
			name: "validation error",
			body: `{"title": ""}`,
			setupMock: func(m *mocks.MockBulletService) {
				m.EXPECT().
					CreateBullet(port.CreateBulletRequest{Title: ""}).
					Return(nil, &entity.ValidationError{Err: entity.ErrTitleRequired}).
					Once()
			},
			wantStatus: http.StatusBadRequest,
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

			if tc.wantBody != nil {
				var got entity.Bullet
				require.NoError(t, json.NewDecoder(res.Body).Decode(&got))
				assert.Equal(t, *tc.wantBody, got)
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
