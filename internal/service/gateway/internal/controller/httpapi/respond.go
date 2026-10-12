package httpapi

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type response[T any] struct {
	Data  T              `json:"data"`
	Error *errorResponse `json:"error,omitempty"`
}

func RespondData[T any](w http.ResponseWriter, code int, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(response[T]{
		Data: data,
	})
}

func RespondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(response[any]{
		Error: &errorResponse{
			Code:    code,
			Message: message,
		},
	})
}
