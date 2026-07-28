package httpapi

import (
	"encoding/json"
	"net/http"
)

type responseStatus string

const (
	statusSuccess responseStatus = "success"
	statusError   responseStatus = "error"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type response[T any] struct {
	Status responseStatus `json:"status"`
	Data   T              `json:"data,omitempty"`
	Error  *errorResponse `json:"error,omitempty"`
}

type listData[T any] struct {
	Items []T `json:"items"`
}

func respondData[T any](w http.ResponseWriter, code int, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(response[T]{
		Status: statusSuccess,
		Data:   data,
	})
}

func respondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(response[any]{
		Status: statusError,
		Error: &errorResponse{
			Code:    code,
			Message: message,
		},
	})
}
