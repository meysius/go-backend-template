# Go Starter Template

A minimal Go web API using [Gin](https://github.com/gin-gonic/gin), PostgreSQL, [sqlc](https://sqlc.dev), and [Atlas](https://atlasgo.io).

## Requirements

- Go 1.25+
- PostgreSQL

## Getting Started

```bash
cp .env.example .env   # fill in your database credentials
make install           # install air, sqlc, atlas
make db-create         # create the main database
make db-dev-create     # create the atlas dev database (<name>_dev)
make migrate           # apply migrations
make dev               # start dev server with hot-reload
```

## Development

```bash
make dev               # start dev server with hot-reload
make build             # compile binary to bin/app
make test              # run all tests
make clean             # remove build artifacts
make tidy              # sync go.mod and go.sum with imports
make sqlc-gen          # regenerate sqlc code from SQL files
make db-create         # create the database
make db-drop           # drop the database
make db-dev-create     # create the atlas dev database
make db-dev-drop       # drop the atlas dev database
make migrate           # apply all pending migrations
make migrate-down      # roll back last migration
make migrate-diff      # generate a migration from schema diff (prompts for name)
```

## Project Structure

```
.
├── main.go                              # entry point, dependency wiring
├── controllers/
│   └── users_controller.go             # HTTP layer, injects service
├── domain/
│   └── identity/
│       ├── identity_service.go         # business logic
│       ├── identity_repository.go      # interface + postgres implementation
│       ├── identity_schema.sql         # canonical table definitions (source of truth)
│       └── identity_queries.sql        # named queries for sqlc
├── db/                                 # sqlc-generated code (do not edit)
├── migrations/                         # Atlas migration files
├── sqlc.yaml                           # sqlc config
├── Makefile
├── go.mod
└── go.sum
```

## Schema Workflow

Each domain slice owns a canonical `*_schema.sql` file that defines the desired state of its tables. Migrations are generated automatically by diffing the canonical schema against the current database.

To change the schema:
1. Edit `domain/<slice>/<slice>_schema.sql`
2. Run `make migrate-diff` — Atlas generates the migration file
3. Review the generated file in `migrations/`
4. Run `make migrate` to apply it
5. Run `make sqlc-gen` to regenerate Go types

Never edit migration files after they've been committed — Atlas checksums will break.

## Dependencies

```bash
go get github.com/some/package
```
