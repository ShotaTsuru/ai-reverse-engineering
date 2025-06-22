.PHONY: help build up down logs clean dev-frontend dev-backend test

# Help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Docker commands
build: ## Build all Docker images
	docker-compose build

up: ## Start all services
	docker-compose up -d

down: ## Stop all services
	docker-compose down

logs: ## Show logs for all services
	docker-compose logs -f

clean: ## Clean up Docker resources
	docker-compose down -v --remove-orphans
	docker system prune -f

# Development commands
dev: ## Start development environment
	docker-compose up

dev-frontend: ## Start frontend development server
	cd frontend && npm run dev

dev-backend: ## Start backend development server
	cd backend && go run main.go

# Database commands
db-up: ## Start only database services
	docker-compose up -d postgres redis

db-down: ## Stop database services
	docker-compose stop postgres redis

# Testing
test-frontend: ## Run frontend tests
	cd frontend && npm test

test-backend: ## Run backend tests
	cd backend && go test ./...

test: test-frontend test-backend ## Run all tests

# Installation
install: ## Install dependencies
	cd frontend && npm install
	cd backend && go mod tidy

# Environment setup
setup: ## Set up development environment
	cp env.example .env
	@echo "Please edit .env file with your configuration"

# Build for production
build-prod: ## Build for production
	cd frontend && npm run build
	cd backend && go build -o main main.go 