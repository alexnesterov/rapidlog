package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeBulletService struct {
	createFn func(port.CreateBulletRequest) (*entity.Bullet, error)
}

func (f *fakeBulletService) CreateBullet(req port.CreateBulletRequest) (*entity.Bullet, error) {
	return f.createFn(req)
}

func (f *fakeBulletService) ListBullets() ([]*entity.Bullet, error) {
	return nil, nil
}

func (f *fakeBulletService) ReadBullet(id uuid.UUID) (*entity.Bullet, error) {
	return nil, nil
}

func (f *fakeBulletService) UpdateBullet(req port.UpdateBulletRequest) (*entity.Bullet, error) {
	return nil, nil
}

func (f *fakeBulletService) DeleteBullet(id uuid.UUID) error {
	return nil
}

func TestCreateBullet(t *testing.T) {
	body := strings.NewReader(`{"title": "Заголовок"}`)

	want := &entity.Bullet{
		ID:    uuid.New(),
		Title: "Заголовок",
	}

	var got port.CreateBulletRequest
	bulletHandler := &BulletHandler{
		BulletService: &fakeBulletService{
			createFn: func(req port.CreateBulletRequest) (*entity.Bullet, error) {
				got = req
				return want, nil
			},
		},
	}

	handler := http.HandlerFunc(bulletHandler.CreateBullet)

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bullets", body)
	req.Header.Set("Content-Type", "application/json")

	handler.ServeHTTP(res, req)

	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	assert.Equal(t, "Заголовок", got.Title)
	assert.Equal(t, http.StatusCreated, res.Code)

	wantJSON, err := json.Marshal(want)
	require.NoError(t, err)
	assert.JSONEq(t, string(wantJSON), res.Body.String())
}

func TestCreateBullet_InvalidJSON(t *testing.T) {
	body := strings.NewReader("не json")

	bulletHandler := &BulletHandler{
		BulletService: &fakeBulletService{},
	}

	handler := http.HandlerFunc(bulletHandler.CreateBullet)

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bullets", body)
	req.Header.Set("Content-Type", "application/json")

	handler.ServeHTTP(res, req)

	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestCreateBullet_ServiceError(t *testing.T) {
	body := strings.NewReader(`{"title": "Заголовок"}`)

	bulletHandler := &BulletHandler{
		BulletService: &fakeBulletService{
			createFn: func(req port.CreateBulletRequest) (*entity.Bullet, error) {
				return nil, errors.New("service error")
			},
		},
	}

	handler := http.HandlerFunc(bulletHandler.CreateBullet)

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bullets", body)
	req.Header.Set("Content-Type", "application/json")

	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}
