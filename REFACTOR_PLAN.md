# Refactor Plan

## Keep
- Go backend with Gin, Ent, PostgreSQL, Redis, Viper, JWT
- Vue 3 frontend with Vite, TypeScript, TailwindCSS, Pinia, Vue Router, axios, vue-i18n
- Basic infra concepts: config loading, logger, unified response/error envelope, middleware patterns, Docker Compose, env examples

## Remove
- AI gateway domain models, services, handlers, routes, menus and UI
- Provider adapters and relay logic: OpenAI, Claude, Gemini, Antigravity, CRS, proxy, relay, upstream account pool
- Billing, quota, pricing, subscription, API key distribution, model scheduling, sticky session, request forwarding
- Product-specific docs, naming, setup wizard, external partner/demo/product references

## Create
- New backend skeleton under `backend/internal` with app/config/server/middleware/web/auth/domain/infrastructure/pkg/util
- New Ent schemas for `User`, `Role`, `Menu`, `AuditLog`, `SystemConfig`
- Seed data and startup auto-migration for minimal demo
- JWT auth + RBAC permission check + audit logging
- New frontend skeleton under `frontend/src` with auth/dashboard/system/user/audit/exception pages and admin layout
- New README and `.env.example` for generic admin starter semantics

## Risks
- Ent generated code replacement requires deleting old generated artifacts to avoid compile conflicts
- Frontend full replacement may leave stale references in build config; build validation will catch this
- Existing compose and deployment scripts are product-oriented; they need aggressive simplification
- Redis is kept as optional runtime dependency for cache/session extension, but minimal demo should still start without custom business queues
