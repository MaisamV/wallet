GO_VERSION := 1.24.4
APP_NAME := wallet
MIGRATION_DIR := ./scripts/migrations
DATABASE_URL ?= postgres://postgres:password@localhost:5432/go_clean_db?sslmode=disable

.PHONY: help
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

# =============================================================================
# Development Commands
# =============================================================================

.PHONY: run
run: ## Start development environment with docker-compose
	@echo Starting development environment...
	docker-compose up -d app
	@echo open http://localhost:8080/swagger to access APIs

.PHONY: run-no-cache
run-no-cache: ## Start development environment with docker-compose
	@echo Starting development environment...
	docker-compose up --build -d app
	@echo open http://localhost:8080/swagger to access APIs

.PHONY: stop
stop: ## Stop development environment
	@echo "Stopping development environment..."
	docker-compose down

.PHONY: delete
delete: ## Stop development environment
	@echo "Stopping development environment..."
	docker-compose down -v

.PHONY: logs
logs: ## Show development environment logs
	docker-compose logs -f

# =============================================================================
# Testing Commands
# =============================================================================

.PHONY: test
test: ## Run all tests
	@echo Running tests...
	go test -v ./...

.PHONY: bench
bench: ## Run all tests
	@echo Running benchmarks...
	docker-compose up -d postgres_test migrate_test
	go test -v ./... -bench=. -benchmem

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo Running tests with coverage...
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo Coverage report generated: coverage.html

.PHONY: test-integration
test-integration: ## Run integration tests
	@echo Running integration tests...
	go test -v -tags=integration ./test/...

# =============================================================================
# Default target
# =============================================================================

.DEFAULT_GOAL := help