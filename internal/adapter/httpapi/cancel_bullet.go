package httpapi

import (
	"errors"
	"net/http"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/google/uuid"
)

func (h *bulletHandler) CancelBullet(w http.ResponseWriter, r *http.Request) {
	userID, _ := UserIDFromContext(r.Context())

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, ErrInvalidID.Error())
		return
	}

	bullet, err := h.usecase.CancelBullet(r.Context(), id, userID)
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

	RespondData(w, http.StatusOK, bullet)
}
