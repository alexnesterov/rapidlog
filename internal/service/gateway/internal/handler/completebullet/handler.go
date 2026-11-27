package completebullet

import (
	"log/slog"
	"net/http"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/platform/server"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/pkg/respond"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

type Handler struct {
	api    journalv1.JournalServiceClient
	logger *slog.Logger
}

func NewHandler(api journalv1.JournalServiceClient, logger *slog.Logger) *Handler {
	return &Handler{api: api, logger: logger}
}

func Setup(mux *http.ServeMux, api journalv1.JournalServiceClient, logger *slog.Logger) {
	h := NewHandler(api, logger)
	mux.HandleFunc("POST /api/bullets/{id}/complete", h.Handle)
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userID, ok := server.UserIDFromContext(r.Context())
	if !ok {
		h.logger.Error("user id not found in context")
		respond.Error(w,
			http.StatusUnauthorized,
			http.StatusText(http.StatusUnauthorized),
		)
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.logger.Error("failed to parse bullet id", "error", err)
		respond.Error(w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	req := &journalv1.CompleteBulletRequest{
		Id:     id.String(),
		UserId: userID.String(),
	}

	resp, err := h.api.CompleteBullet(r.Context(), req)
	if err != nil {
		h.logger.Error("failed to complete bullet", "error", err)
		respond.Error(w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		)
		return
	}

	data, err := protojson.Marshal(resp)
	if err != nil {
		h.logger.Error("failed to marshal response", "error", err)
		respond.Error(w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		)
		return
	}

	respond.Bytes(w, http.StatusOK, data)
}
