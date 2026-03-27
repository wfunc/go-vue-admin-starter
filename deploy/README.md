# Deploy Notes

This directory contains the runtime entrypoint used by the generic starter Docker image.

Primary local development and demo deployment now use the repository root files:

- `docker-compose.yml`
- `.env.example`

Quick start:

```bash
cp .env.example .env
docker compose up --build
```

Services:

- `postgres`: PostgreSQL 16
- `redis`: Redis 7
- `backend`: Gin + Ent API server
- `frontend`: Vite development server
