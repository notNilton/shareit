# shareit

Photo sharing platform monorepo with image processing pipeline and multi-platform clients.

## Architecture

```
apps/
  backend/            Go API and background image processing worker
  web/                React web application (Vite + TypeScript)
  backoffice/         React admin dashboard (Vite + TypeScript)
  mobile/             React Native Android application
infra/
  postgres/           Database initialization scripts
docker-compose.yml     Docker orchestration
```

### Components

- `apps/backend`: HTTP API (`cmd/api`) and Redis queue worker (`cmd/worker`) for thumbnail generation.
- `apps/web`: User web interface.
- `apps/backoffice`: Administrative dashboard.
- `apps/mobile`: Android native mobile application.

### Image Processing Pipeline

1. `POST /photos`: Uploads original image to MinIO (`media` bucket, `originals/<id>`), saves pending metadata in PostgreSQL, and enqueues processing job to Redis.
2. `worker`: Consumes Redis queue, downloads original file, generates a 400px JPEG thumbnail (`thumbnails/<id>.jpg`), and updates metadata status to `ready`.

## Development

### Prerequisites

- Go 1.26+
- Node.js 22+
- Docker / Podman
- JDK and Android SDK (for mobile)

### Running Services

Start local infrastructure (PostgreSQL, Redis, MinIO):

```bash
docker compose up -d postgres redis minio minio-init
```

Start the full stack including API and worker:

```bash
docker compose up -d
```

### Manual Execution

Backend services:

```bash
cd apps/backend
go run ./cmd/api
go run ./cmd/worker
```

Web client:

```bash
cd apps/web
npm install
npm run dev
```

Mobile client:

```bash
cd apps/mobile
npm install
npm run android
```

### Service Endpoints

| Service | Type | Port | Endpoint |
|---------|------|------|----------|
| Web App | Frontend | `5173` | http://localhost:5173 |
| Backoffice | Frontend | `5174` | http://localhost:5174 |
| HTTP API | Backend | `8080` | http://localhost:8080 |
| PostgreSQL | Database | `5432` | localhost:5432/shareit |
| Redis | Cache/Queue | `6379` | localhost:6379 |
| MinIO API | Storage | `9000` | http://localhost:9000 |
| MinIO Console | Storage UI | `9001` | http://localhost:9001 |


## Documentation

- [📋 Roadmap & TODOs](docs/TODO.md) - Planned features and project roadmap
- [📐 Architecture](docs/ARCHITECTURE.md) - System architecture and components
- [📄 License](LICENSE) - MIT License
