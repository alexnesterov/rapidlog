package main

import (
	"log"
	"net/http"

	"github.com/alexnesterov/rapidlog-api/internal/adapter/httpapi"
	"github.com/alexnesterov/rapidlog-api/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()

	router.HandleFunc("/", httpapi.Greet)
	router.HandleFunc("/health", httpapi.Health)

	log.Printf("%s is starting on port %s", cfg.App.Name, cfg.HTTP.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.HTTP.Port, router))
}
