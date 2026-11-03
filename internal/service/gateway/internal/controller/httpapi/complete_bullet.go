package httpapi

import (
	"context"
	"log/slog"
	"net/http"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/pkg/respond"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

type CompleteBulletClient interface {
	CompleteBullet(
		ctx context.Context,
		in *journalv1.CompleteBulletRequest,
		opts ...grpc.CallOption,
	) (*journalv1.CompleteBulletResponse, error)
}

type CompleteBulletHandler struct {
	api    CompleteBulletClient
	logger *slog.Logger
}

func NewCompleteBulletHandler(api CompleteBulletClient, logger *slog.Logger) *CompleteBulletHandler {
	return &CompleteBulletHandler{api: api, logger: logger}
}

func (h *CompleteBulletHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
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
