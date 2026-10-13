package httpapi

import (
	"context"
	"net/http"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"google.golang.org/grpc"
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
	userID, _ := UserIDFromContext(r.Context())

	bullets, err := h.api.ListBullets(r.Context(), &journalv1.ListBulletsRequest{
		UserId: userID.String(),
	})
	if err != nil {
		RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	RespondData(w, http.StatusOK, bullets)
}
