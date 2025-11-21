.PHONY: help setup build run stop clean test migrate-up migrate-down swagger


help:
	@echo "Available commands:"
	@echo "  make setup       - Install dependencies"
	@echo "  make swagger     - Generate Swagger docs for all services"
	@echo "  make build       - Build all services"
	@echo "  make run         - Start all services via docker-compose"
	@echo "  make stop        - Stop all services"
	@echo "  make clean       - Clean up containers and volumes"
	@echo "  make test        - Run tests"
	@echo "  make migrate-up  - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"


setup:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Installing swag..."
	@go install github.com/swaggo/swag/cmd/swag@latest || true


swagger:
	@echo "Generating Swagger docs..."
	@cd services/auth && command swag init -g cmd/main.go -o docs --parseDependency --parseInternal
	@cd services/contact && command swag init -g cmd/main.go -o docs --parseDependency --parseInternal
	@cd services/inventory && command swag init -g cmd/main.go -o docs --parseDependency --parseInternal
	@cd services/sales && command swag init -g cmd/main.go -o docs --parseDependency --parseInternal
	@cd services/purchase && command swag init -g cmd/main.go -o docs --parseDependency --parseInternal
	@echo "✅ Swagger docs generated"

sqlc:
	@echo "Generating SQL code with sqlc..."
	@cd services/auth && sqlc generate
	@echo "✅ Auth service sqlc code generated"
	@echo "⚠️  Other services pending sqlc setup"


build:
	@echo "Building all services..."
	docker compose build


run:
	@echo "Starting all services..."
	docker compose up -d
	@echo "Services are starting..."
	@echo "Check status with: docker compose ps"


stop:
	@echo "Stopping all services..."
	docker compose down


clean:
	@echo "Cleaning up..."
	docker compose down -v
	docker system prune -f


test:
	@echo "Running tests..."
	go test ./...


migrate-up:
	@echo "Running migrations..."
	@echo "TODO: Implement with golang-migrate"

migrate-down:
	@echo "Rolling back migrations..."
	@echo "TODO: Implement with golang-migrate"


run-auth:
	docker compose up -d db-auth nats
	@echo "Auth service dependencies started"

run-gateway:
	docker compose up -d gateway
	@echo "Gateway started"


health:
	@echo "Checking service health..."
	docker compose ps

