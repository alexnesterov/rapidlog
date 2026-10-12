package httpapi

import "net/http"

func EmptyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}