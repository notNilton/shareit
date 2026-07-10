package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"shareit/backend/internal/photos"
	"shareit/backend/internal/queue"
	"shareit/backend/internal/storage"
)

const maxUploadSize = 20 << 20 // 20 MiB

type Server struct {
	pool    *pgxpool.Pool
	storage *storage.Storage
	queue   *queue.Queue
	photos  *photos.Repository
}

func NewRouter(pool *pgxpool.Pool, store *storage.Storage, q *queue.Queue) http.Handler {
	s := &Server{
		pool:    pool,
		storage: store,
		queue:   q,
		photos:  photos.NewRepository(pool),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /health/db", s.handleHealthDB)
	mux.HandleFunc("POST /photos", s.handleUploadPhoto)
	mux.HandleFunc("POST /photos/presigned", s.handlePresignedUpload)
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleHealthDB(w http.ResponseWriter, r *http.Request) {
	if err := s.pool.Ping(r.Context()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "error", "error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleUploadPhoto accepts a multipart "file" field, stores the original in
// object storage, records it in Postgres, and enqueues a thumbnail job.
func (s *Server) handleUploadPhoto(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing file field"})
		return
	}
	defer file.Close()

	ctx := r.Context()
	id := uuid.NewString()
	originalKey := "originals/" + id

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	if err := s.storage.Upload(ctx, originalKey, file, header.Size, contentType); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to store photo"})
		return
	}

	if err := s.photos.Create(ctx, id, originalKey); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to save photo metadata"})
		return
	}

	if err := s.queue.PushImageJob(ctx, queue.ImageJob{PhotoID: id, OriginalKey: originalKey}); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to enqueue processing job"})
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{"id": id, "status": "pending"})
}

func (s *Server) handlePresignedUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := uuid.NewString()
	originalKey := "originals/" + id

	uploadURL, err := s.storage.GetPresignedUploadURL(ctx, originalKey, 15*time.Minute)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate upload URL"})
		return
	}

	if err := s.photos.Create(ctx, id, originalKey); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to save photo metadata"})
		return
	}

	if err := s.queue.PushImageJob(ctx, queue.ImageJob{PhotoID: id, OriginalKey: originalKey}); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to enqueue processing job"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"id":         id,
		"upload_url": uploadURL,
		"status":     "pending",
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
