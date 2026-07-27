package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
)

type bulletHandler struct {
	usecase port.BulletService
}

func NewBulletHandler(uc port.BulletService) *bulletHandler {
	return &bulletHandler{
		usecase: uc,
	}
}

func (h *bulletHandler) CreateBullet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req port.CreateBulletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdBullet, err := h.usecase.CreateBullet(req)
	if err != nil {
		var validationErr *entity.ValidationError
		if errors.As(err, &validationErr) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respondData(w, http.StatusCreated, createdBullet)
}

func (h *bulletHandler) ListBullets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bullets, err := h.usecase.ListBullets()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if bullets == nil {
		bullets = []*entity.Bullet{}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(bullets)
}
