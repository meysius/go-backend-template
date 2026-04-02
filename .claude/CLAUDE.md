# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Dev Commands

```bash
make dev               # hot-reload dev server (air)
make run               # run without hot-reload
make build             # compile to bin/app
make test              # run all tests
go test ./domain/...   # run tests for a specific package subtree
go test ./domain/identity/ -run TestGetUser  # run a single test
make swagger           # regenerate OpenAPI docs (swag init)
```

## Database & Schema Commands

```bash
make sqlc-gen          # regenerate Go code from SQL files
make migrate           # apply pending migrations
make migrate-down      # roll back last migration
make migrate-diff      # generate migration from schema diff (prompts for name)
make db-create         # create the database
make db-dev-create     # create the atlas dev database (<name>_dev)
```

## Architecture

Domain-slice architecture with three layers: **controllers** (HTTP) → **services** (business logic) → **repositories** (data access via interface).

### Composition Root

`app.go` (`NewApp`) is the composition root — it instantiates all concrete types, wires dependencies, and sets up routes via `Mount()`. `main.go` only calls `NewApp()` and starts the server.

### Domain Slices

Each domain lives under `domain/<name>/` and owns:
- `<name>_service.go` — business logic, depends on repository interface only
- `<name>_repo.go` — repository interface definition + pgx/sqlc implementation in the same file
- `<name>_schema.sql` — canonical table DDL (source of truth for schema)
- `<name>_queries.sql` — SQL queries consumed by sqlc

Current slices: `identity` (users), `ordering` (products).

### Generated Code

`db/` contains sqlc-generated code — never edit manually. All slices' generated code lands in a single `db` package. Both schema and query SQL files must be registered in `sqlc.yaml`.

### Error Handling Pattern

Each domain defines its own sentinel errors (e.g. `identity.ErrNotFound`). Repositories translate `pgx.ErrNoRows` → domain sentinel. Controllers match sentinels with `errors.Is()` to choose HTTP status codes.

### API Docs

Swagger annotations on controller methods → `make swagger` → Scalar UI served at `/docs`, OpenAPI JSON at `/docs/openapi.json`.

## Key Workflows

**Add a SQL query:** edit `domain/<slice>/<slice>_queries.sql` → `make sqlc-gen`

**Change schema:** edit `domain/<slice>/<slice>_schema.sql` → `make migrate-diff` → review migration → `make migrate` → `make sqlc-gen`

**Add a new domain slice:**
1. Create `domain/<name>/` with service, repo, schema, and queries files
2. Register both SQL files in `sqlc.yaml`
3. `make sqlc-gen` → `make migrate-diff` → `make migrate`
4. Create `controllers/<name>_controller.go`
5. Wire in `app.go` (`NewApp` + `Mount`)

## Constraints

- Migrations are append-only — never modify committed migration files (Atlas checksums)
- `<slice>_schema.sql` is the source of truth — migrations must make the DB match it
- Config comes from `.env` (godotenv) or real env vars; required: `DATABASE_NAME`, `DATABASE_USER`, `DATABASE_PASS`
