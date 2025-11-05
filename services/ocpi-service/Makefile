.PHONY: help build run test test-coverage clean docker-up docker-down test-integration test-all

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building..."
	go build -o bin/server cmd/server/main.go

run: ## Run the application
	@echo "Running..."
	go run cmd/server/main.go

test: ## Run unit tests only (fast)
	@echo "Running unit tests..."
	go test -v -short ./...

test-integration: docker-up ## Run integration tests (requires Docker)
	@echo "Running integration tests..."
	@echo "Waiting for MongoDB to be ready..."
	@sleep 3
	go test -v -run Integration ./internal/api/...

test-all: docker-up ## Run all tests (unit + integration)
	@echo "Running all tests..."
	@echo "Waiting for MongoDB to be ready..."
	@sleep 3
	go test -v ./...

test-coverage: ## Run all tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-unit: ## Run unit tests only (alias for test)
	@echo "Running unit tests..."
	go test -v -short ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	cd ../../ && docker-compose up -d

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	cd ../../ && docker-compose down

docker-logs: ## View Docker container logs
	@echo "Viewing logs..."
	cd ../../ && docker-compose logs -f

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

dev: docker-up run ## Start dev environment and run server
