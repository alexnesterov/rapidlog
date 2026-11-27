package health

import (
	"encoding/json"
	"net/http"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func Setup(mux *http.ServeMux) {
	h := NewHandler()
	mux.HandleFunc("/health", h.Handle)
}

type healthResponse struct {
	Status string `json:"status"`
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	resp := healthResponse{Status: "ok"}
	status := http.StatusOK

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}
