package listbullets

import (
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
	mux.HandleFunc("GET /api/bullets", h.Handle)
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userID, ok := server.UserIDFromContext(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	resp, err := h.api.ListBullets(r.Context(), &journalv1.ListBulletsRequest{
		UserId: userID.String(),
	})
	if err != nil {
		respond.Error(w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		)
		return
	}

	data, err := protojson.Marshal(resp)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respond.Bytes(w, http.StatusOK, data)
}
