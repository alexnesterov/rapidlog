package main

import (
	"log"
	"net/http"

	"github.com/alexnesterov/rapidlog-api/internal/adapter/httpapi"
	"github.com/alexnesterov/rapidlog-api/internal/config"
	"github.com/alexnesterov/rapidlog-api/internal/domain/usecase"
	"github.com/alexnesterov/rapidlog-api/internal/infrastructure/memory"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()

	router.HandleFunc("/health", httpapi.HealthHandler)

	bulletRepository := memory.NewBulletRepository()
	bulletService := usecase.NewBulletService(bulletRepository)
	bulletHandler := httpapi.NewBulletHandler(bulletService)

	router.HandleFunc("POST /api/bullets", bulletHandler.CreateBullet)
	router.HandleFunc("GET /api/bullets", bulletHandler.ListBullets)

	log.Printf("%s is starting on port %s", cfg.App.Name, cfg.HTTP.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.HTTP.Port, router))
}
