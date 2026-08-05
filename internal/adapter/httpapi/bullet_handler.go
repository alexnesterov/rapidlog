package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/google/uuid"
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
	var input port.CreateBulletInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		RespondError(w, http.StatusBadRequest, errInvalidRequest.Error())
		return
	}

	createdBullet, err := h.usecase.CreateBullet(input)
	if err != nil {
		var validationErr *entity.ValidationError
		if errors.As(err, &validationErr) {
			RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		RespondError(w, http.StatusInternalServerError, errInternal.Error())
		return
	}

	RespondData(w, http.StatusCreated, createdBullet)
}

func (h *bulletHandler) ListBullets(w http.ResponseWriter, r *http.Request) {
	bullets, err := h.usecase.ListBullets()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, errInternal.Error())
		return
	}

	if bullets == nil {
		bullets = []*entity.Bullet{}
	}

	RespondData(w, http.StatusOK, groupBulletsByDay(bullets))
}

func (h *bulletHandler) CompleteBullet(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, errParseID.Error())
		return
	}

	completedBullet, err := h.usecase.CompleteBullet(id)
	if err != nil {
		if errors.Is(err, port.ErrNotFound) {
			RespondError(w, http.StatusNotFound, "bullet not found")
			return
		}

		RespondError(w, http.StatusInternalServerError, errInternal.Error())
		return
	}

	RespondData(w, http.StatusOK, completedBullet)
}
