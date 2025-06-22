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

test-frontend-e2e: ## Run frontend E2E tests
	cd frontend && npm run test:e2e

test-backend: ## Run backend tests
	cd backend && go test ./...

test-backend-unit: ## Run backend unit tests
	cd backend && go test -v -race -coverprofile=coverage.out ./...

test-backend-integration: ## Run backend integration tests
	cd backend && go test -v -tags=integration ./tests/integration/...

test-coverage: ## Show test coverage
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: backend/coverage.html"

test: test-frontend test-backend ## Run all tests

test-all: test-frontend test-frontend-e2e test-backend-unit test-backend-integration ## Run all tests including E2E

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