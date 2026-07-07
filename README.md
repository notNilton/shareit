# shareit

Monorepo do projeto shareit.

## Estrutura

```
shareit/
├── apps/
│   ├── backend/      Go: API (cmd/api) + worker de processamento de imagem (cmd/worker)
│   ├── web/           Webapp em React (Vite + TypeScript)
│   ├── backoffice/    Painel administrativo em React (Vite + TypeScript)
│   └── mobile/        App nativo Android em React Native (TypeScript)
├── infra/
│   └── postgres/      Scripts de inicialização do banco
└── docker-compose.yml Sobe Postgres, Redis, MinIO, API e worker localmente
```

## Backend (apps/backend)

Requer Go 1.26+. Depende de Postgres, Redis e um storage S3-compatible (MinIO) já
rodando — veja a seção "Infra local" abaixo.

```
cd apps/backend
go run ./cmd/api      # API HTTP
go run ./cmd/worker   # worker que processa a fila de imagens
```

Variáveis de ambiente:
- `PORT` (default `8080`)
- `DATABASE_URL` (default `postgres://shareit:shareit@localhost:5432/shareit?sslmode=disable`)
- `REDIS_URL` (default `redis://localhost:6379/0`)
- `S3_ENDPOINT` (default `localhost:9000`)
- `S3_ACCESS_KEY` / `S3_SECRET_KEY` (default `shareit` / `shareit123`)
- `S3_BUCKET` (default `media`)
- `S3_USE_SSL` (default `false`)

### Fluxo de upload de foto

1. `POST /photos` (multipart, campo `file`) — a API salva o arquivo original no
   bucket `media` (`originals/<id>`), grava metadados no Postgres (status `pending`)
   e enfileira um job no Redis.
2. O `worker` consome a fila, baixa o original do storage, gera um thumbnail
   (400px, JPEG) e sobe em `thumbnails/<id>.jpg`, atualizando o status para `ready`.

## Web (apps/web) e Backoffice (apps/backoffice)

Requer Node 22+.

```
cd apps/web        # ou apps/backoffice
npm install
npm run dev
```

## Mobile (apps/mobile)

App React Native focado em Android. Requer o ambiente React Native configurado
(JDK, Android SDK, variável `ANDROID_HOME`) — veja o [guia oficial](https://reactnative.dev/docs/set-up-your-environment).

```
cd apps/mobile
npm install
npm run android
```

## Infra local

```
docker compose up -d postgres redis minio minio-init
```

Isso sobe:
- **Postgres** em `localhost:5432` (banco `shareit`, usuário/senha `shareit`),
  aplicando `infra/postgres/init.sql` na primeira inicialização.
- **Redis** em `localhost:6379`, usado como fila de processamento de imagem.
- **MinIO** (S3-compatible) em `localhost:9000` (API) e `localhost:9001` (console web,
  login `shareit` / `shareit123`). O serviço `minio-init` cria o bucket `media`
  automaticamente.

Para subir tudo, incluindo API e worker:

```
docker compose up -d
```
