package httpapi

import (
	"context"
	"io"
	"log/slog"
	"net/http"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/pkg/respond"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

type createBulletClient interface {
	CreateBullet(
		ctx context.Context,
		in *journalv1.CreateBulletRequest,
		opts ...grpc.CallOption,
	) (*journalv1.CreateBulletResponse, error)
}

type createBulletHandler struct {
	api    createBulletClient
	logger *slog.Logger
}

func NewCreateBulletHandler(api createBulletClient, logger *slog.Logger) http.HandlerFunc {
	h := &createBulletHandler{api: api, logger: logger}
	return h.Handle
}

func (h *createBulletHandler) Handle(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
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
