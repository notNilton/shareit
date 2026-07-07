// Package photos persists photo upload metadata and processing status.
package photos

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusReady   Status = "ready"
	StatusFailed  Status = "failed"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, id, originalKey string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO photos (id, original_key, status) VALUES ($1, $2, $3)`,
		id, originalKey, StatusPending,
	)
	return err
}

func (r *Repository) MarkReady(ctx context.Context, id, thumbKey string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE photos SET thumb_key = $2, status = $3 WHERE id = $1`,
		id, thumbKey, StatusReady,
	)
	return err
}

func (r *Repository) MarkFailed(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE photos SET status = $2 WHERE id = $1`,
		id, StatusFailed,
	)
	return err
}
