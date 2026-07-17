package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
)

type BulletHandler struct {
	BulletService port.BulletService
}

func (h *BulletHandler) CreateBullet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req port.CreateBulletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdBullet, err := h.BulletService.CreateBullet(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(createdBullet)
}
