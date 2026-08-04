package main

import (
	"context"
	"log"
	"net/http"

	"github.com/alexnesterov/rapidlog-api/internal/adapter/httpapi"
	"github.com/alexnesterov/rapidlog-api/internal/config"
	"github.com/alexnesterov/rapidlog-api/internal/domain/usecase"
	"github.com/alexnesterov/rapidlog-api/internal/infrastructure/memory"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DB.DSN)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("failed to ping database: ", err)
	}
	log.Println("Database connected successfully")

	router := http.NewServeMux()

	router.HandleFunc("/health", httpapi.HealthHandler)

	bulletRepository := memory.NewBulletRepository()
	bulletService := usecase.NewBulletService(bulletRepository)
	bulletHandler := httpapi.NewBulletHandler(bulletService)

	router.HandleFunc("POST /api/bullets", bulletHandler.CreateBullet)
	router.HandleFunc("GET /api/bullets", bulletHandler.ListBullets)
	router.HandleFunc("POST /api/bullets/{id}/complete", bulletHandler.CompleteBullet)

	log.Printf("%s is starting on port %s", cfg.App.Name, cfg.HTTP.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.HTTP.Port, router))
}
