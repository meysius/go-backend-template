-include .env
export

.PHONY: dev build test clean install tidy sqlc-gen migrate migrate-down db-create db-drop db-dev-create db-dev-drop migrate-diff swagger docker-build docker-run

DATABASE_URL=postgres://$(DATABASE_USER):$(DATABASE_PASS)@$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=disable
DATABASE_DEV_URL=postgres://$(DATABASE_USER):$(DATABASE_PASS)@$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)_dev?sslmode=disable

install:
	go install github.com/air-verse/air@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	curl -sSf https://atlasgo.sh | sh

swagger:
	swag init --parseDependency --parseInternal

dev:
	air

build:
	go build -o bin/app main.go

test:
	go test ./...

clean:
	rm -rf bin/

tidy:
	go mod tidy

sqlc-gen:
	sqlc generate

migrate:
	atlas migrate apply --url "$(DATABASE_URL)" --dir "file://migrations"

migrate-down:
	atlas migrate down --url "$(DATABASE_URL)" --dir "file://migrations" --dev-url "$(DATABASE_DEV_URL)"

db-create:
	createdb $(DATABASE_NAME)

db-drop:
	dropdb $(DATABASE_NAME)

db-dev-create:
	createdb $(DATABASE_NAME)_dev

db-dev-drop:
	dropdb $(DATABASE_NAME)_dev

docker-build:
	docker build -t go-starter-template .

docker-run:
	docker run --env-file .env -p 3000:3000 go-starter-template

migrate-diff:
	@read -p "Migration name: " name; \
	atlas migrate diff $$name \
	  $(shell find domain -name '*_schema.sql' | sort | sed 's|^|--to file://|') \
	  --dev-url "$(DATABASE_DEV_URL)" \
	  --dir "file://migrations"; \
	echo "Remember to run: make sqlc-gen"
