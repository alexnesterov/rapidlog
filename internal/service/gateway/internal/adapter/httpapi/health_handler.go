package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type healthResponse struct {
	Status string `json:"status"`
	DB     string `json:"db"`
}

type Pinger interface {
	Ping(ctx context.Context) error
}

func NewHealthHandler(pinger Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		resp := healthResponse{Status: "ok", DB: "ok"}
		status := http.StatusOK

		if err := pinger.Ping(ctx); err != nil {
			resp.Status = "error"
			resp.DB = "error"
			status = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
