package httpapi

import "net/http"

func ListBullets(w http.ResponseWriter, r *http.Request) {
	RespondData(w, http.StatusOK, []string{})
}