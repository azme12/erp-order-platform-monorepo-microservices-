.PHONY: help setup build run stop clean test migrate-up migrate-down

# Default target
help:
	@echo "Available commands:"
	@echo "  make setup       - Install dependencies"
	@echo "  make build       - Build all services"
	@echo "  make run         - Start all services via docker-compose"
	@echo "  make stop        - Stop all services"
	@echo "  make clean       - Clean up containers and volumes"
	@echo "  make test        - Run tests"
	@echo "  make migrate-up  - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"

# Install dependencies
setup:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Build all services
build:
	@echo "Building all services..."
	docker compose build

# Run all services
run:
	@echo "Starting all services..."
	docker compose up -d
	@echo "Services are starting..."
	@echo "Check status with: docker compose ps"

# Stop all services
stop:
	@echo "Stopping all services..."
	docker compose down

# Clean up
clean:
	@echo "Cleaning up..."
	docker compose down -v
	docker system prune -f

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Database migrations (placeholder - will be implemented with golang-migrate)
migrate-up:
	@echo "Running migrations..."
	@echo "TODO: Implement with golang-migrate"

migrate-down:
	@echo "Rolling back migrations..."
	@echo "TODO: Implement with golang-migrate"

# Individual service commands
run-auth:
	docker compose up -d db-auth nats
	@echo "Auth service dependencies started"

run-gateway:
	docker compose up -d gateway
	@echo "Gateway started"

# Check service health
health:
	@echo "Checking service health..."
	docker compose ps

