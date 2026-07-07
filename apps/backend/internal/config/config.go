package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	RedisURL    string
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3UseSSL    bool
}

func Load() Config {
	return Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://shareit:shareit@localhost:5434/shareit?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379/0"),
		S3Endpoint:  getEnv("S3_ENDPOINT", "localhost:9000"),
		S3AccessKey: getEnv("S3_ACCESS_KEY", "shareit"),
		S3SecretKey: getEnv("S3_SECRET_KEY", "shareit123"),
		S3Bucket:    getEnv("S3_BUCKET", "media"),
		S3UseSSL:    getEnv("S3_USE_SSL", "false") == "true",
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
