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
		RespondError(w, http.StatusBadRequest, ErrInvalidRequestBody.Error())
		return
	}

	userID, _ := UserIDFromContext(r.Context())
	input.UserID = userID

	createdBullet, err := h.usecase.CreateBullet(r.Context(), input)
	if err != nil {
		var validationErr *entity.ValidationError
		if errors.As(err, &validationErr) {
			RespondError(w, http.StatusBadRequest, validationErr.Error())
			return
		}

		RespondError(w, http.StatusInternalServerError, ErrInternal.Error())
		return
	}

	RespondData(w, http.StatusCreated, createdBullet)
}

func (h *bulletHandler) ListBullets(w http.ResponseWriter, r *http.Request) {
	userID, _ := UserIDFromContext(r.Context())

	bullets, err := h.usecase.ListBullets(r.Context(), userID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, ErrInternal.Error())
		return
	}

	if bullets == nil {
		bullets = []*entity.Bullet{}
	}

	RespondData(w, http.StatusOK, groupBulletsByDay(bullets))
}

func (h *bulletHandler) CompleteBullet(w http.ResponseWriter, r *http.Request) {
	userID, _ := UserIDFromContext(r.Context())

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, ErrInvalidID.Error())
		return
	}

	completedBullet, err := h.usecase.CompleteBullet(r.Context(), id, userID)
	if err != nil {
		if errors.Is(err, port.ErrNotFound) {
			RespondError(w, http.StatusNotFound, "bullet not found")
			return
		}

		RespondError(w, http.StatusInternalServerError, ErrInternal.Error())
		return
	}

	RespondData(w, http.StatusOK, completedBullet)
}

func (h *bulletHandler) MigrateBullet(w http.ResponseWriter, r *http.Request) {
	userID, _ := UserIDFromContext(r.Context())

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	migratedBullet, err := h.usecase.MigrateBullet(r.Context(), id, userID)
	if err != nil {
		if errors.Is(err, port.ErrNotFound) {
			RespondError(w, http.StatusNotFound, "bullet not found")
			return
		}

		var validationErr *entity.ValidationError
		if errors.As(err, &validationErr) {
			RespondError(w, http.StatusBadRequest, validationErr.Error())
			return
		}

		RespondError(w, http.StatusInternalServerError, ErrInternal.Error())
		return
	}

	RespondData(w, http.StatusCreated, migratedBullet)
}
