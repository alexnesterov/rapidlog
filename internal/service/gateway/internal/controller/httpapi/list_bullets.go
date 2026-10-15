package httpapi

import (
	"context"
	"net/http"
	"time"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/pkg/protoenum"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/pkg/respond"
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
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		respond.Error(w,
			http.StatusUnauthorized,
			http.StatusText(http.StatusUnauthorized),
		)
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

	bullets := resp.GetData()

	data := make([]struct {
		ID        string    `json:"id"`
		Type      string    `json:"type"`
		Signifier string    `json:"signifier"`
		Content   string    `json:"content"`
		UserId    string    `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}, len(bullets))

	for i, b := range bullets {
		data[i].ID = b.Id
		data[i].Type = protoenum.ToCleanString(b.Type)
		data[i].Signifier = protoenum.ToCleanString(b.Signifier)
		data[i].Content = b.Content
		data[i].UserId = b.UserId
		data[i].CreatedAt = b.CreatedAt.AsTime()
		data[i].UpdatedAt = b.UpdatedAt.AsTime()
	}

	respond.Data(w, http.StatusOK, data)
}
