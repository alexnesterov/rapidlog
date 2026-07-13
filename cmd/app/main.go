package main

import (
	"log"
	"net/http"

	"github.com/alexnesterov/dotline/internal/adapter/httpapi"
)

func main() {
	http.HandleFunc("/", httpapi.Greet)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
