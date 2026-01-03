include .env
export

.PHONY: help setup migrate-up migrate-down crawler indexer api docker-up docker-down test clean

help:
	@echo "Available commands:"
	@echo "  make setup        - Setup database and dependencies"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"
	@echo "  make crawler      - Run crawler"
	@echo "  make indexer      - Run indexer"
	@echo "  make api          - Run search API"
	@echo "  make docker-up    - Start all services with Docker"
	@echo "  make docker-down  - Stop all services"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"

setup:
	@echo "Setting up project..."
	cd crawler && go mod download
	cd search-api && go mod download
	cd indexer && cargo build

migrate-up:
	@echo "Running migrations..."
	psql "$(DATABASE_URL)" < migrations/001_create_documents.up.sql
	psql "$(DATABASE_URL)" < migrations/002_create_index.up.sql

migrate-down:
	@echo "Rolling back migrations..."
	psql "$(DATABASE_URL)" < migrations/002_create_index.down.sql
	psql "$(DATABASE_URL)" < migrations/001_create_documents.down.sql

crawler:
	cd crawler && go run cmd/crawler/main.go

indexer:
	cd indexer && cargo run --release

api:
	cd search-api && go run cmd/api/main.go

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

test:
	cd crawler && go test ./...
	cd search-api && go test ./...
	cd indexer && cargo test

clean:
	cd crawler && go clean
	cd search-api && go clean
	cd indexer && cargo clean
	docker-compose down -v
