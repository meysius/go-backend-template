# Go Starter Template — Claude Context

## Project Overview
A Go web API starter template using Gin, PostgreSQL (pgx), sqlc, and goose. Follows a domain-slice architecture with clear separation between HTTP, business logic, and data access layers.

## Architecture

### Layer Structure
```
controllers/        → HTTP layer (Gin handlers). Injects services.
domain/<slice>/
  <slice>_service.go      → Business logic. Injects repository interface.
  <slice>_repository.go   → Data access interface + PostgreSQL implementation.
  <slice>_schema.sql      → Canonical table definitions (source of truth for schema + sqlc input)
  <slice>_queries.sql     → Named SQL queries (input to sqlc)
db/                 → sqlc-generated code (all slices). DO NOT edit manually.
migrations/         → Atlas migration files (must match canonical schema files)
main.go             → Wires everything together (composition root)
```

### Dependency Flow
```
main.go → controller → service → repository (interface)
                                      ↑
                              userRepository (pgx impl, uses sqlc Queries)
```

### Domain Slices
Each domain (e.g. `identity`) lives under `domain/<name>/` and owns its service, repository, SQL, and generated code. Adding a new domain means creating a new slice — do not mix concerns across slices.

## Key Conventions

- **File naming** — files are named `<entity>_<layer>.go` (e.g. `users_controller.go`, `identity_service.go`, `identity_repository.go`) for project-wide searchability. The `db/` directory is exempt (sqlc-generated).
- **Errors propagate** — repository methods return errors, services pass them through, controllers map them to HTTP status codes. Never swallow errors silently.
- **`ErrNotFound`** is defined per domain slice and used to map to 404 in controllers.
- **Repository interface** is the boundary — the service only depends on the interface, never the concrete implementation.
- **`<slice>_schema.sql` is the canonical schema** — it defines the desired state of the DB for that slice. Migrations must be written to make the DB match it. Never let them drift.
- **sqlc generates** all DB query code. To change queries: edit `<slice>_queries.sql` → run `make sqlc-gen`. Never edit `db/` by hand.
- **Migrations are append-only** — never modify existing migration files (Atlas checksums will break). Each migration must leave the DB consistent with the canonical schema files.
- **`main.go` is the composition root** — it's the only place that instantiates concrete types and wires dependencies.

## Environment & Configuration
- Config comes from `.env` (loaded by godotenv at startup) or real env vars (production).
- `.env` is gitignored. `.env.example` is the template committed to the repo.
- Key vars: `DATABASE_URL`, `DATABASE_NAME`, `DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_USER`, `DATABASE_PASS`.
- Makefile loads `.env` via `-include .env` + `export`, so all `make` targets have access to these vars.

## Dev Workflow
```bash
cp .env.example .env      # first time setup
make install              # install air, sqlc, atlas
make db-create            # create the database
make db-dev-create        # create the atlas dev database (<name>_dev)
make migrate              # run migrations
make dev                  # start with hot-reload (air)
```

### Common Tasks
| Task | Command |
|------|---------|
| Add a SQL query | Edit `domain/<slice>/<slice>_queries.sql` → `make sqlc-gen` |
| Change schema | Edit `domain/<slice>/<slice>_schema.sql` → `make migrate-diff` (prompts for name) → `make sqlc-gen` |
| Add a migration | `make migrate-diff` (Atlas generates it from schema diff, prompts for name) |
| Roll back migration | `make migrate-down` |
| Rebuild binary | `make build` |
| Run tests | `make test` |

## Tech Stack
| Tool | Purpose |
|------|---------|
| [Gin](https://github.com/gin-gonic/gin) | HTTP router |
| [pgx/v5](https://github.com/jackc/pgx) | PostgreSQL driver |
| [sqlc](https://sqlc.dev) | Type-safe SQL code generation |
| [atlas](https://atlasgo.io) | Database migrations (schema diff + apply) |
| [air](https://github.com/air-verse/air) | Hot-reload dev server |
| [godotenv](https://github.com/joho/godotenv) | .env file loading |

## Adding a New Domain Slice
1. Create `domain/<name>/<name>_service.go` and `<name>_repository.go`
2. Create `domain/<name>/<name>_schema.sql` (canonical table definitions)
3. Create `domain/<name>/<name>_queries.sql` (named SQL queries)
4. Add both files to `sqlc.yaml` under `schema` and `queries`
5. Run `make sqlc-gen`
6. Run `make migrate-diff` to generate the migration (prompts for name)
7. Create `controllers/<name>_controller.go`
8. Wire it in `main.go`
