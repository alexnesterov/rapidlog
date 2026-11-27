package createbullet

import (
	"io"
	"log/slog"
	"net/http"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/platform/server"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/pkg/respond"
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
	mux.HandleFunc("POST /api/bullets", h.Handle)
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userID, ok := server.UserIDFromContext(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed to read request body", "error", err)
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var req journalv1.CreateBulletRequest
	if err := protojson.Unmarshal(bodyBytes, &req); err != nil {
		h.logger.Error("failed to unmarshal request body", "error", err)
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.UserId = userID.String()

	resp, err := h.api.CreateBullet(r.Context(), &req)
	if err != nil {
		h.logger.Error("failed to create bullet", "error", err)
		respond.Error(w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		)
		return
	}

	data, err := protojson.Marshal(resp)
	if err != nil {
		h.logger.Error("failed to marshal response", "error", err)
		respond.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respond.Bytes(w, http.StatusCreated, data)
}
