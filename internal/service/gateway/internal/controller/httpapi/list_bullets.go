package httpapi

import (
	"context"
	"net/http"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/pkg/respond"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

type ListBulletsClient interface {
	ListBullets(
		ctx context.Context,
		in *journalv1.ListBulletsRequest,
		opts ...grpc.CallOption,
	) (*journalv1.ListBulletsResponse, error)
}

type ListBulletsHandler struct {
	api ListBulletsClient
}

func NewListBulletsHandler(api ListBulletsClient) *ListBulletsHandler {
	return &ListBulletsHandler{api: api}
}

func (h *ListBulletsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
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
