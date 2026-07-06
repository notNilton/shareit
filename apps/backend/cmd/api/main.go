package main

import (
	"context"
	"log"
	"net/http"

	"shareit/backend/internal/config"
	"shareit/backend/internal/db"
	"shareit/backend/internal/httpserver"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	router := httpserver.NewRouter(pool)

	log.Printf("listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
