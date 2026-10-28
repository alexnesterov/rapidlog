package httpapi

import (
	"context"
	"net/http"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/pkg/respond"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

type listBulletsClient interface {
	ListBullets(ctx context.Context, in *journalv1.ListBulletsRequest, opts ...grpc.CallOption) (*journalv1.ListBulletsResponse, error)
}

type listBulletsHandler struct {
	api listBulletsClient
}

func NewListBulletsHandler(api listBulletsClient) http.HandlerFunc {
	h := &listBulletsHandler{api: api}
	return h.Handle
}

func (h *listBulletsHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	marshalOptions := &protojson.MarshalOptions{
		UseEnumNumbers: false,
		UseProtoNames:  true,
	}

	data, err := marshalOptions.Marshal(resp)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respond.Bytes(w, http.StatusOK, data)
}
