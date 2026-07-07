// Package queue provides a minimal Redis list-backed job queue.
package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const ImageProcessingQueue = "shareit:image-processing"

type ImageJob struct {
	PhotoID     string `json:"photo_id"`
	OriginalKey string `json:"original_key"`
}

type Queue struct {
	client *redis.Client
}

func New(redisURL string) (*Queue, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	return &Queue{client: redis.NewClient(opts)}, nil
}

func (q *Queue) Ping(ctx context.Context) error {
	return q.client.Ping(ctx).Err()
}

func (q *Queue) Close() error {
	return q.client.Close()
}

func (q *Queue) PushImageJob(ctx context.Context, job ImageJob) error {
	payload, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return q.client.LPush(ctx, ImageProcessingQueue, payload).Err()
}

// PopImageJob blocks up to timeout waiting for a job. It returns (nil, nil) on timeout.
func (q *Queue) PopImageJob(ctx context.Context, timeout time.Duration) (*ImageJob, error) {
	result, err := q.client.BRPop(ctx, timeout, ImageProcessingQueue).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var job ImageJob
	if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
		return nil, err
	}
	return &job, nil
}
