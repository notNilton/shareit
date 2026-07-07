package main

import (
	"context"
	"log"
	"net/http"

	"shareit/backend/internal/config"
	"shareit/backend/internal/db"
	"shareit/backend/internal/httpserver"
	"shareit/backend/internal/queue"
	"shareit/backend/internal/storage"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	store, err := storage.New(cfg)
	if err != nil {
		log.Fatalf("failed to init object storage: %v", err)
	}
	if err := store.EnsureBucket(ctx); err != nil {
		log.Fatalf("failed to ensure bucket: %v", err)
	}

	q, err := queue.New(cfg.RedisURL)
	if err != nil {
		log.Fatalf("failed to init queue: %v", err)
	}
	defer q.Close()
	if err := q.Ping(ctx); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	router := httpserver.NewRouter(pool, store, q)

	log.Printf("listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
