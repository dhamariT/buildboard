.PHONY: help build run dev test clean install docker-build docker-up docker-down docker-logs

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install dependencies
	go mod download
	go mod tidy

build: ## Build the server binary
	@echo "Building server..."
	@mkdir -p bin
	go build -o bin/server cmd/server/main.go
	@echo "Build complete: bin/server"

run: build ## Build and run the server
	@echo "Starting server..."
	./bin/server

dev: ## Run the server in development mode (with go run)
	@echo "Starting server in development mode..."
	go run cmd/server/main.go

test: ## Run tests
	go test -v ./...

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: fmt vet ## Run formatters and linters

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	go clean
	@echo "Clean complete"

setup: ## Copy .env.example to .env
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo ".env file created. Please update it with your configuration."; \
	else \
		echo ".env file already exists."; \
	fi

# Docker commands
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t buildboard-backend:latest .
	@echo "Docker image built successfully"

docker-up: ## Start services with docker-compose
	@echo "Starting services..."
	docker-compose up -d
	@echo "Services started. API available at http://localhost:8080"
	@echo "Run 'make docker-logs' to view logs"

docker-up-tools: ## Start services with Adminer (database UI)
	@echo "Starting services with tools..."
	docker-compose --profile tools up -d
	@echo "Services started."
	@echo "API: http://localhost:8080"
	@echo "Adminer (DB UI): http://localhost:8081"

docker-down: ## Stop services
	@echo "Stopping services..."
	docker-compose down
	@echo "Services stopped"

docker-logs: ## View logs
	docker-compose logs -f

docker-ps: ## Show running containers
	docker-compose ps

docker-restart: ## Restart services
	@echo "Restarting services..."
	docker-compose restart
	@echo "Services restarted"

docker-clean: ## Stop and remove all containers and volumes
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	@echo "Docker resources cleaned"

docker-prod-build: ## Build production Docker image
	@echo "Building production Docker image..."
	docker build -t buildboard-backend:prod .
	@echo "Production image built"

docker-prod-up: ## Start production services
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Run 'make setup' first."; \
		exit 1; \
	fi
	@echo "Starting production services..."
	docker-compose -f docker-compose.prod.yml up -d
	@echo "Production services started"

docker-prod-down: ## Stop production services
	docker-compose -f docker-compose.prod.yml down