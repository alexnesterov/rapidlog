package main

import (
	"log"
	"net/http"

	"github.com/alexnesterov/dotline/internal/handler"
)

func main() {
	http.HandleFunc("/", handler.Greet)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
