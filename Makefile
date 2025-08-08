# Chinese Bridge Game Makefile

.PHONY: help build run test clean docker-up docker-down k8s-up k8s-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Go Development
build: ## Build all Go services
	@echo "Building Go services..."
	go mod tidy
	go build -o bin/auth-service ./cmd/auth-service
	go build -o bin/user-service ./cmd/user-service
	go build -o bin/game-service ./cmd/game-service

run-auth: ## Run auth service
	@echo "Starting auth service..."
	./bin/auth-service

run-user: ## Run user service
	@echo "Starting user service..."
	./bin/user-service

run-game: ## Run game service
	@echo "Starting game service..."
	./bin/game-service

test: ## Run Go tests
	@echo "Running Go tests..."
	go test -v ./...

# Flutter Development
flutter-get: ## Get Flutter dependencies
	@echo "Getting Flutter dependencies..."
	cd flutter_app && flutter pub get

flutter-build: ## Build Flutter app
	@echo "Building Flutter app..."
	cd flutter_app && flutter build apk

flutter-run: ## Run Flutter app
	@echo "Running Flutter app..."
	cd flutter_app && flutter run

flutter-test: ## Run Flutter tests
	@echo "Running Flutter tests..."
	cd flutter_app && flutter test

# Docker Development
docker-up: ## Start all services with Docker Compose
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

docker-down: ## Stop all Docker services
	@echo "Stopping Docker services..."
	docker-compose down

docker-logs: ## Show Docker logs
	docker-compose logs -f

docker-build: ## Build Docker images
	@echo "Building Docker images..."
	docker build -t chinese-bridge/auth-service -f docker/auth-service.Dockerfile .
	docker build -t chinese-bridge/user-service -f docker/user-service.Dockerfile .
	docker build -t chinese-bridge/game-service -f docker/game-service.Dockerfile .

# Kubernetes Development
k8s-up: ## Deploy to Kubernetes
	@echo "Deploying to Kubernetes..."
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/postgres.yaml
	kubectl apply -f k8s/redis.yaml
	kubectl apply -f k8s/kafka.yaml
	kubectl apply -f k8s/auth-service.yaml

k8s-down: ## Remove from Kubernetes
	@echo "Removing from Kubernetes..."
	kubectl delete -f k8s/auth-service.yaml
	kubectl delete -f k8s/kafka.yaml
	kubectl delete -f k8s/redis.yaml
	kubectl delete -f k8s/postgres.yaml
	kubectl delete -f k8s/namespace.yaml

k8s-status: ## Check Kubernetes status
	kubectl get all -n chinese-bridge

# Database
db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	# This will be implemented when GORM migrations are added

db-seed: ## Seed database with test data
	@echo "Seeding database..."
	# This will be implemented when seed scripts are added

# Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	cd flutter_app && flutter clean
	docker-compose down -v
	docker system prune -f

# Development Setup
setup: ## Setup development environment
	@echo "Setting up development environment..."
	cp .env.example .env
	go mod tidy
	cd flutter_app && flutter pub get
	docker-compose up -d postgres redis kafka

# Code Quality
lint: ## Run linters
	@echo "Running Go linters..."
	golangci-lint run
	@echo "Running Flutter linters..."
	cd flutter_app && flutter analyze

format: ## Format code
	@echo "Formatting Go code..."
	go fmt ./...
	@echo "Formatting Flutter code..."
	cd flutter_app && dart format .