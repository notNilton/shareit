# shareit

Monorepo do projeto shareit.

## Estrutura

```
shareit/
├── apps/
│   ├── backend/      Go API (net/http + pgx), fala com o Postgres
│   ├── web/           Webapp em React (Vite + TypeScript)
│   ├── backoffice/    Painel administrativo em React (Vite + TypeScript)
│   └── mobile/        App nativo Android em React Native (TypeScript)
├── infra/
│   └── postgres/      Scripts de inicialização do banco
└── docker-compose.yml Sobe Postgres + backend localmente
```

## Backend (apps/backend)

Requer Go 1.26+.

```
cd apps/backend
go run ./cmd/api
```

Variáveis de ambiente:
- `PORT` (default `8080`)
- `DATABASE_URL` (default `postgres://shareit:shareit@localhost:5432/shareit?sslmode=disable`)

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

## Banco de dados

```
docker compose up -d postgres
```

Isso sobe um Postgres em `localhost:5432` com o banco `shareit` (usuário/senha `shareit`)
e aplica `infra/postgres/init.sql` na primeira inicialização.

Para subir o backend junto:

```
docker compose up -d
```
