# Database Migration Guide

The starter uses Ent schema auto-migration during application bootstrap.

## Initialize schema and seed data

```bash
cd backend
GOTOOLCHAIN=local GOSUMDB=off go run ./cmd/server --migrate-only
```

This command will:

- create or update tables for `users`, `roles`, `menus`, `audit_logs`, `system_configs`
- seed built-in roles
- seed the default admin account
- seed starter menus and public system configuration entries

## Start the server normally

```bash
cd backend
GOTOOLCHAIN=local GOSUMDB=off go run ./cmd/server
```

## Future production migration strategy

For larger projects, replace the auto-migration flow with versioned SQL or Ent migration files in this directory, then invoke them before application startup.
