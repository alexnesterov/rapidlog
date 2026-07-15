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

	http.HandleFunc("/", httpapi.Greet)

	log.Printf("%s is starting on port %s", cfg.App.Name, cfg.HTTP.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.HTTP.Port, nil))
}
