package httpapi

import (
	"encoding/json"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response[T any] struct {
	Data  T              `json:"data"`
	Error *ErrorResponse `json:"error,omitempty"`
}

func RespondData[T any](w http.ResponseWriter, code int, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(Response[T]{
		Data: data,
	})
}

func RespondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(Response[any]{
		Error: &ErrorResponse{
			Code:    code,
			Message: message,
		},
	})
}

type ProtoJSONWrapper struct {
	Msg proto.Message
}

func (w ProtoJSONWrapper) MarshalJSON() ([]byte, error) {
	marshaler := protojson.MarshalOptions{
		UseProtoNames: true,
	}
	return marshaler.Marshal(w.Msg)
}

func WrapProtoSlice[T proto.Message](slice []T) []ProtoJSONWrapper {
	wrapped := make([]ProtoJSONWrapper, len(slice))
	for i, msg := range slice {
		wrapped[i] = ProtoJSONWrapper{Msg: msg}
	}
	return wrapped
}
