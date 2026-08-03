// Package storage wraps an S3-compatible object store (MinIO in development).
package storage

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"

	"github.com/minio/minio-go/v7/pkg/credentials"

	"shareit/backend/internal/config"
)

type Storage struct {
	client *minio.Client
	bucket string
}

func New(cfg config.Config) (*Storage, error) {
	client, err := minio.New(cfg.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		Secure: cfg.S3UseSSL,
	})
	if err != nil {
		return nil, err
	}
	return &Storage{client: client, bucket: cfg.S3Bucket}, nil
}

func (s *Storage) EnsureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{})
}

func (s *Storage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	_, err := s.client.PutObject(ctx, s.bucket, key, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (s *Storage) Download(ctx context.Context, key string) (*minio.Object, error) {
	return s.client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
}

// GetPresignedUploadURL gera uma URL assinada para o cliente fazer upload direto no MinIO.
func (s *Storage) GetPresignedUploadURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	u, err := s.client.PresignedPutObject(ctx, s.bucket, key, expires)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// GetPresignedDownloadURL gera uma URL assinada para leitura direta de imagem/thumbnail.
func (s *Storage) GetPresignedDownloadURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	u, err := s.client.PresignedGetObject(ctx, s.bucket, key, expires, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

