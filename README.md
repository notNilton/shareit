# shareit

Monorepo for the shareit project.

## Structure

```
shareit/
├── apps/
│   ├── backend/      Go: API (cmd/api) + image processing worker (cmd/worker)
│   ├── web/           React webapp (Vite + TypeScript)
│   ├── backoffice/    React admin panel (Vite + TypeScript)
│   └── mobile/        React Native Android app (TypeScript)
├── infra/
│   └── postgres/      Database initialization scripts
└── docker-compose.yml Runs Postgres, Redis, MinIO, API, and worker locally
```

## Backend (apps/backend)

Requires Go 1.26+. Depends on Postgres, Redis, and an S3-compatible storage (MinIO) running — see the "Local Infrastructure" section below.

```bash
cd apps/backend
go run ./cmd/api      # HTTP API
go run ./cmd/worker   # Worker processing the image queue
```

Environment variables:
- `PORT` (default `8080`)
- `DATABASE_URL` (default `postgres://shareit:shareit@localhost:5432/shareit?sslmode=disable`)
- `REDIS_URL` (default `redis://localhost:6379/0`)
- `S3_ENDPOINT` (default `localhost:9000`)
- `S3_ACCESS_KEY` / `S3_SECRET_KEY` (default `shareit` / `shareit123`)
- `S3_BUCKET` (default `media`)
- `S3_USE_SSL` (default `false`)

### Photo Upload Flow

1. `POST /photos` (multipart, `file` field) — the API saves the original file to the `media` bucket (`originals/<id>`), writes metadata to Postgres (status `pending`), and queues a job in Redis.
2. The `worker` consumes the queue, downloads the original from storage, generates a thumbnail (400px, JPEG), uploads it to `thumbnails/<id>.jpg`, and updates the status to `ready`.

## Web (apps/web) and Backoffice (apps/backoffice)

Requires Node 22+.

```bash
cd apps/web        # or apps/backoffice
npm install
npm run dev
```

## Mobile (apps/mobile)

Android-focused React Native application. Requires a configured React Native environment (JDK, Android SDK, `ANDROID_HOME` variable) — see the [official guide](https://reactnative.dev/docs/set-up-your-environment).

```bash
cd apps/mobile
npm install
npm run android
```

## Local Infrastructure

```bash
docker compose up -d postgres redis minio minio-init
```

This starts:
- **Postgres** on `localhost:5432` (`shareit` database, `shareit`/`shareit` credentials), applying `infra/postgres/init.sql` on first startup.
- **Redis** on `localhost:6379`, used for image processing queues.
- **MinIO** (S3-compatible) on `localhost:9000` (API) and `localhost:9001` (web console, login `shareit` / `shareit123`). The `minio-init` service automatically creates the `media` bucket.

To run everything, including API and worker:

```bash
docker compose up -d
```
