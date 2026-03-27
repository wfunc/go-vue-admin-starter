# Go Vue Admin Starter

A reusable admin starter extracted from the original `Wei-Shaw/sub2api` repository.

This repository is no longer positioned as an AI gateway. The active codebase is now a generic admin skeleton focused on reusable infrastructure and common back-office capabilities.

## Positioning

This starter is intended to be the starting point for future admin projects.

Included:

- Go backend with `Gin + Ent`
- Vue 3 frontend with `Vite + TypeScript + TailwindCSS`
- Built-in `zh / en` frontend i18n with topbar language switching
- PostgreSQL and Redis integration
- JWT authentication and RBAC permission checks
- Unified response envelope and error handling
- Audit logging middleware
- Minimal CRUD examples for users, roles, menus and system configs
- Built-in self-service account features for profile update and password change

Removed from the active skeleton:

- AI provider integration
- upstream account pool and relay logic
- API key distribution
- token billing, quota, pricing, subscription logic
- model routing, scheduling, sticky session and gateway semantics

Legacy product-specific materials have been isolated under `_legacy_*` directories and are no longer part of the active starter runtime.

## Tech Stack

| Layer | Stack |
| --- | --- |
| Backend | Go 1.25, Gin, Ent, Viper, JWT |
| Frontend | Vue 3, Vite 5, TypeScript, TailwindCSS, Pinia, Vue Router |
| Database | PostgreSQL |
| Cache | Redis |
| DevOps | Docker Compose |

## Active Modules

Backend modules:

- `auth`
- `user`
- `role`
- `menu`
- `audit_log`
- `system_config`

Frontend pages:

- Login
- Dashboard
- User Management
- Role Management
- Menu Management
- System Config
- Audit Log
- Profile Center
- 403 Forbidden
- 404 Not Found

## Directory Structure

```text
backend/
  cmd/server/
  ent/
  internal/
    app/
    auth/
    config/
    domain/
      audit/
      auth/
      menu/
      role/
      system/
      user/
    infrastructure/
      cache/
      db/
    middleware/
    pkg/
    server/
    util/
    web/
      errorx/
      response/
  migrations/

frontend/
  src/
    api/
      modules/
    components/
      base/
      business/
    composables/
    i18n/
    layouts/
    router/
    stores/
    styles/
    types/
    utils/
    views/
      audit/
      auth/
      dashboard/
      exception/
      system/
      user/
```

## Local Start

Important: the backend depends on PostgreSQL, and the health check also verifies Redis by default. Start `postgres` and `redis` first, otherwise the backend bootstrap will fail.

Recommended startup order:

1. Start PostgreSQL and Redis.
2. Run database migration and seed data.
3. Start the backend API.
4. Start the frontend dev server.

### 0. Start PostgreSQL And Redis First

Use Docker Compose from the repository root:

```bash
cd /Volumes/NewDisk/Github/Wei-Shaw/sub2api
docker compose --env-file .env.example up -d postgres redis
```

Check dependency status:

```bash
docker compose --env-file .env.example ps -a
```

Expected result:

- `admin-starter-postgres` is `healthy`
- `admin-starter-redis` is `healthy`

### 1. Backend

Run migration before the first backend startup:

```bash
cd /Volumes/NewDisk/Github/Wei-Shaw/sub2api/backend
GOTOOLCHAIN=local GOSUMDB=off go run ./cmd/server --migrate-only
```

If PostgreSQL is not running, this command will fail immediately with a database connection hint.

```bash
cd /Volumes/NewDisk/Github/Wei-Shaw/sub2api/backend
GOTOOLCHAIN=local GOSUMDB=off go run ./cmd/server
```

Backend default address:

- `http://localhost:18080`
- health check: `GET /healthz`

### 2. Frontend

```bash
cd /Volumes/NewDisk/Github/Wei-Shaw/sub2api/frontend
npm install
npm run dev
```

Frontend default address:

- `http://localhost:5173`

The Vite dev server proxies `/api` to `http://localhost:18080` by default.

## Docker Compose Start

```bash
cd /Volumes/NewDisk/Github/Wei-Shaw/sub2api
cp .env.example .env
docker compose up --build
```

Default addresses:

- frontend: `http://localhost:5173`
- backend: `http://localhost:18080`
- postgres: `localhost:5432`
- redis: `localhost:6379`

## Database Initialization

Schema creation and seed data are handled during backend bootstrap.

Run migration only:

```bash
cd /Volumes/NewDisk/Github/Wei-Shaw/sub2api/backend
GOTOOLCHAIN=local GOSUMDB=off go run ./cmd/server --migrate-only
```

Seeded demo account:

- username: `admin`
- password: `admin123456`

Verified local dependency-first workflow:

```bash
cd /Volumes/NewDisk/Github/Wei-Shaw/sub2api
docker compose --env-file .env.example up -d postgres redis

cd /Volumes/NewDisk/Github/Wei-Shaw/sub2api/backend
GOTOOLCHAIN=local GOSUMDB=off go run ./cmd/server --migrate-only
GOTOOLCHAIN=local GOSUMDB=off go run ./cmd/server

cd /Volumes/NewDisk/Github/Wei-Shaw/sub2api/frontend
npm install
npm run dev
```

Health check after backend startup:

```bash
curl http://127.0.0.1:18080/healthz
```

Expected response:

```json
{"code":"success","message":"success","data":{"database":"up","redis":"up"}}
```

## Build and Verification

Backend:

```bash
cd backend
GOTOOLCHAIN=local GOSUMDB=off go build ./...
```

Frontend:

```bash
cd frontend
npm install
npm run build
```

## Environment Variables

The starter reads configuration from environment variables through Viper.

Common variables:

- `SERVER_PORT`
- `DATABASE_HOST`
- `DATABASE_PORT`
- `DATABASE_USER`
- `DATABASE_PASSWORD`
- `DATABASE_NAME`
- `REDIS_ENABLED`
- `REDIS_ADDR`
- `JWT_SECRET`
- `SEED_ADMIN_USERNAME`
- `SEED_ADMIN_PASSWORD`

See [`.env.example`](./.env.example) for a complete local template.

## How To Add a New Business Module

Use the existing modules as a template.

1. Add a new Ent schema in `backend/ent/schema`.
2. Generate Ent code.
3. Create backend slices under `backend/internal/domain/<module>/`:
   - `entity`
   - `repository`
   - `service`
   - `handler`
4. Register the handler in `backend/internal/app/app.go` and `backend/internal/server/router.go`.
5. Add frontend API methods under `frontend/src/api/modules`.
6. Add the page under `frontend/src/views/<module>`.
7. Register the route in `frontend/src/router/index.ts`.
8. Seed a menu item and assign a permission if the module should appear in the admin navigation.

## Notes

- The repository name and Go module path are still inherited from the original GitHub repository for convenience.
- Active business semantics have been replaced with starter-oriented naming and modules.
- `_legacy_*` directories are retained only as isolated historical material during this refactor.
