package main

import (
	"bytes"
	"context"
	"errors"
	"log"
	"time"

	"github.com/disintegration/imaging"

	"shareit/backend/internal/config"
	"shareit/backend/internal/db"
	"shareit/backend/internal/photos"
	"shareit/backend/internal/queue"
	"shareit/backend/internal/storage"
)

const (
	popTimeout    = 5 * time.Second
	thumbnailSize = 400
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

	q, err := queue.New(cfg.RedisURL)
	if err != nil {
		log.Fatalf("failed to init queue: %v", err)
	}
	defer q.Close()

	repo := photos.NewRepository(pool)

	const maxConcurrentWorkers = 3
	sem := make(chan struct{}, maxConcurrentWorkers)

	log.Println("worker started, waiting for image jobs (max concurrency:", maxConcurrentWorkers, ")")
	for {
		job, err := q.PopImageJob(ctx, popTimeout)
		if err != nil {
			log.Printf("error popping job: %v", err)
			continue
		}
		if job == nil {
			continue // timeout, no job available
		}

		sem <- struct{}{}
		go func(j queue.ImageJob) {
			defer func() { <-sem }()
			if err := processWithRetry(ctx, store, repo, j, 3); err != nil {
				log.Printf("failed to process photo %s after retries: %v", j.PhotoID, err)
				if err := repo.MarkFailed(ctx, j.PhotoID); err != nil {
					log.Printf("failed to mark photo %s as failed: %v", j.PhotoID, err)
				}
				return
			}
			log.Printf("processed photo %s", j.PhotoID)
		}(*job)
	}
}

func processWithRetry(ctx context.Context, store *storage.Storage, repo *photos.Repository, job queue.ImageJob, retries int) error {
	var err error
	for i := 0; i < retries; i++ {
		if err = process(ctx, store, repo, job); err == nil {
			return nil
		}
		time.Sleep(time.Duration(i+1) * 2 * time.Second)
	}
	return err
}

func process(ctx context.Context, store *storage.Storage, repo *photos.Repository, job queue.ImageJob) error {
	original, err := store.Download(ctx, job.OriginalKey)
	if err != nil {
		return err
	}
	defer original.Close()

	img, err := imaging.Decode(original, imaging.AutoOrientation(true))
	if err != nil {
		return errors.New("decode image: " + err.Error())
	}

	thumb := imaging.Resize(img, thumbnailSize, 0, imaging.Lanczos)

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, thumb, imaging.JPEG, imaging.JPEGQuality(85)); err != nil {
		return errors.New("encode thumbnail: " + err.Error())
	}

	thumbKey := "thumbnails/" + job.PhotoID + ".jpg"
	if err := store.Upload(ctx, thumbKey, &buf, int64(buf.Len()), "image/jpeg"); err != nil {
		return err
	}

	return repo.MarkReady(ctx, job.PhotoID, thumbKey)
}
