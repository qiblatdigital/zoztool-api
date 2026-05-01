.PHONY: dev build test migrate

dev:
	@echo "Starting dev server..."
	go run ./cmd/api

build:
	@echo "Building binary..."
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/api ./cmd/api
	@echo "Build complete: bin/api"

test:
	go test ./... -v

migrate:
	@echo "Running migrations..."
	psql "host=$${DB_HOST} port=$${DB_PORT} user=$${DB_USER} password=$${DB_PASSWORD} dbname=$${DB_NAME} sslmode=disable" \
		-f migrations/001_init.sql
	@echo "Migration complete"

tidy:
	go mod tidy

lint:
	golangci-lint run ./...

docker-build:
	docker build -t zoztool-api:latest .

docker-up:
	docker compose up -d

docker-down:
	docker compose down
