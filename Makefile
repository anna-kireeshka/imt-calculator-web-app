include .env
export

DB_DSN := "host=$(DB_HOST) port=$(DB_PORT) dbname=$(DB_NAME) user=$(DB_USER) password=$(DB_PASSWORD) sslmode=$(DB_SSL_MODE)"

## ── Local dev ────────────────────────────────────────────────────────────────

# Перегенерировать *_templ.go после правки .templ.
# Установка инструмента: go install github.com/a-h/templ/cmd/templ@latest
generate:
	templ generate

run: generate
	go run ./cmd/server

build: generate
	go build -o bin/server ./cmd/server


## ── Docker ───────────────────────────────────────────────────────────────────

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-down-v:
	docker compose down -v

docker-logs:
	docker compose logs -f server

docker-build:
	docker compose build --no-cache


## ── Database ─────────────────────────────────────────────────────────────────

migrate-up:
	goose -dir database/migrations postgres $(DB_DSN) up

migrate-down:
	goose -dir database/migrations postgres $(DB_DSN) down

migrate-status:
	goose -dir database/migrations postgres $(DB_DSN) status

migrate-reset:
	goose -dir database/migrations postgres $(DB_DSN) reset

migrate-create:
	@read -p "Migration name: " name; \
	goose -dir database/migrations create $$name sql