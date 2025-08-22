# WohnFair Makefile
# Provides commands for building, testing, and running the system

.PHONY: help all clean proto build test lint format deps dev compose-up compose-down bootstrap

# Default target
help: ## Show this help message
	@echo "WohnFair - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Quick Start:"
	@echo "  make bootstrap    # Install dependencies and start system"
	@echo "  make dev         # Start development environment"
	@echo "  make test        # Run all tests"
	@echo "  make build       # Build all services"

# Variables
GO_VERSION := 1.21
RUST_VERSION := 1.70
NODE_VERSION := 18
PYTHON_VERSION := 3.11

# Colors for output
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
NC := \033[0m # No Color

# Check if required tools are installed
check-tools:
	@echo "$(YELLOW)Checking required tools...$(NC)"
	@command -v go >/dev/null 2>&1 || { echo "$(RED)Go is required but not installed. Please install Go $(GO_VERSION)+$(NC)"; exit 1; }
	@command -v rustc >/dev/null 2>&1 || { echo "$(RED)Rust is required but not installed. Please install Rust $(RUST_VERSION)+$(NC)"; exit 1; }
	@command -v node >/dev/null 2>&1 || { echo "$(RED)Node.js is required but not installed. Please install Node.js $(NODE_VERSION)+$(NC)"; exit 1; }
	@command -v python3 >/dev/null 2>&1 || { echo "$(RED)Python is required but not installed. Please install Python $(PYTHON_VERSION)+$(NC)"; exit 1; }
	@command -v docker >/dev/null 2>&1 || { echo "$(RED)Docker is required but not installed$(NC)"; exit 1; }
	@command -v docker-compose >/dev/null 2>&1 || { echo "$(RED)Docker Compose is required but not installed$(NC)"; exit 1; }
	@command -v buf >/dev/null 2>&1 || { echo "$(YELLOW)buf is not installed. Installing...$(NC)"; go install github.com/bufbuild/buf/cmd/buf@latest; }
	@echo "$(GREEN)All required tools are installed!$(NC)"

# Install dependencies
deps: check-tools ## Install all dependencies
	@echo "$(YELLOW)Installing Go dependencies...$(NC)"
	@cd services/orchestration && go mod download
	@cd services/fairrent && go mod download
	@cd services/policy-dsl && go mod download
	@cd services/notifications && go mod download
	@echo "$(GREEN)Go dependencies installed$(NC)"
	
	@echo "$(YELLOW)Installing Rust dependencies...$(NC)"
	@cd services/zk-lease && cargo fetch
	@echo "$(GREEN)Rust dependencies installed$(NC)"
	
	@echo "$(YELLOW)Installing Python dependencies...$(NC)"
	@cd services/ml && pip install -e .
	@echo "$(GREEN)Python dependencies installed$(NC)"
	
	@echo "$(YELLOW)Installing Node.js dependencies...$(NC)"
	@cd frontend && npm install
	@echo "$(GREEN)Node.js dependencies installed$(NC)"

# Generate protocol buffers
proto: ## Generate protocol buffer code
	@echo "$(YELLOW)Generating protocol buffers...$(NC)"
	@buf generate
	@echo "$(GREEN)Protocol buffers generated$(NC)"

# Build all services
build: proto ## Build all services
	@echo "$(YELLOW)Building Go services...$(NC)"
	@cd services/orchestration && go build -o bin/gateway ./cmd/gateway
	@cd services/fairrent && go build -o bin/fairrentd ./cmd/fairrentd
	@cd services/policy-dsl && go build -o bin/policy-dsl ./cmd/policy-dsl
	@cd services/notifications && go build -o bin/notifier ./cmd/notifier
	@echo "$(GREEN)Go services built$(NC)"
	
	@echo "$(YELLOW)Building Rust service...$(NC)"
	@cd services/zk-lease && cargo build --release
	@echo "$(GREEN)Rust service built$(NC)"
	
	@echo "$(YELLOW)Building frontend...$(NC)"
	@cd frontend && npm run build
	@echo "$(GREEN)Frontend built$(NC)"

# Run all tests
test: ## Run all tests
	@echo "$(YELLOW)Running Go tests...$(NC)"
	@cd services/orchestration && go test -v ./...
	@cd services/fairrent && go test -v ./...
	@cd services/policy-dsl && go test -v ./...
	@cd services/notifications && go test -v ./...
	@echo "$(GREEN)Go tests passed$(NC)"
	
	@echo "$(YELLOW)Running Rust tests...$(NC)"
	@cd services/zk-lease && cargo test
	@echo "$(GREEN)Rust tests passed$(NC)"
	
	@echo "$(YELLOW)Running Python tests...$(NC)"
	@cd services/ml && python -m pytest tests/ -v
	@echo "$(GREEN)Python tests passed$(NC)"
	
	@echo "$(YELLOW)Running frontend tests...$(NC)"
	@cd frontend && npm test
	@echo "$(GREEN)Frontend tests passed$(NC)"

# Run specific test suites
test-go: ## Run only Go tests
	@echo "$(YELLOW)Running Go tests...$(NC)"
	@cd services/orchestration && go test -v ./...
	@cd services/fairrent && go test -v ./...
	@cd services/policy-dsl && go test -v ./...
	@cd services/notifications && go test -v ./...

test-rust: ## Run only Rust tests
	@echo "$(YELLOW)Running Rust tests...$(NC)"
	@cd services/zk-lease && cargo test

test-python: ## Run only Python tests
	@echo "$(YELLOW)Running Python tests...$(NC)"
	@cd services/ml && python -m pytest tests/ -v

test-js: ## Run only JavaScript/TypeScript tests
	@echo "$(YELLOW)Running frontend tests...$(NC)"
	@cd frontend && npm test

# Linting
lint: ## Run all linting
	@echo "$(YELLOW)Running Go linting...$(NC)"
	@cd services/orchestration && golangci-lint run
	@cd services/fairrent && golangci-lint run
	@cd services/policy-dsl && golangci-lint run
	@cd services/notifications && golangci-lint run
	@echo "$(GREEN)Go linting passed$(NC)"
	
	@echo "$(YELLOW)Running Rust linting...$(NC)"
	@cd services/zk-lease && cargo clippy
	@echo "$(GREEN)Rust linting passed$(NC)"
	
	@echo "$(YELLOW)Running Python linting...$(NC)"
	@cd services/ml && flake8 wohnfair_ml/ tests/
	@echo "$(GREEN)Python linting passed$(NC)"
	
	@echo "$(YELLOW)Running frontend linting...$(NC)"
	@cd frontend && npm run lint
	@echo "$(GREEN)Frontend linting passed$(NC)"
	
	@echo "$(YELLOW)Running protocol buffer linting...$(NC)"
	@buf lint
	@echo "$(GREEN)Protocol buffer linting passed$(NC)"

# Formatting
format: ## Format all code
	@echo "$(YELLOW)Formatting Go code...$(NC)"
	@cd services/orchestration && go fmt ./...
	@cd services/fairrent && go fmt ./...
	@cd services/policy-dsl && go fmt ./...
	@cd services/notifications && go fmt ./...
	@echo "$(GREEN)Go code formatted$(NC)"
	
	@echo "$(YELLOW)Formatting Rust code...$(NC)"
	@cd services/zk-lease && cargo fmt
	@echo "$(GREEN)Rust code formatted$(NC)"
	
	@echo "$(YELLOW)Formatting Python code...$(NC)"
	@cd services/ml && black wohnfair_ml/ tests/
	@echo "$(GREEN)Python code formatted$(NC)"
	
	@echo "$(YELLOW)Formatting frontend code...$(NC)"
	@cd frontend && npm run format
	@echo "$(GREEN)Frontend code formatted$(NC)"

# Check formatting
format-check: ## Check if code is properly formatted
	@echo "$(YELLOW)Checking Go formatting...$(NC)"
	@cd services/orchestration && test -z "$(shell go fmt ./...)"
	@cd services/fairrent && test -z "$(shell go fmt ./...)"
	@cd services/policy-dsl && test -z "$(shell go fmt ./...)"
	@cd services/notifications && test -z "$(shell go fmt ./...)"
	@echo "$(GREEN)Go formatting is correct$(NC)"
	
	@echo "$(YELLOW)Checking Rust formatting...$(NC)"
	@cd services/zk-lease && cargo fmt -- --check
	@echo "$(GREEN)Rust formatting is correct$(NC)"
	
	@echo "$(YELLOW)Checking Python formatting...$(NC)"
	@cd services/ml && black --check wohnfair_ml/ tests/
	@echo "$(GREEN)Python formatting is correct$(NC)"
	
	@echo "$(YELLOW)Checking frontend formatting...$(NC)"
	@cd frontend && npm run format:check
	@echo "$(GREEN)Frontend formatting is correct$(NC)"

# Docker Compose commands
compose-up: ## Start all services with Docker Compose
	@echo "$(YELLOW)Starting WohnFair services...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)Services started!$(NC)"
	@echo "$(YELLOW)Access points:$(NC)"
	@echo "  Frontend: http://localhost:3000"
	@echo "  Gateway:  http://localhost:8080"
	@echo "  Keycloak: http://localhost:8081"
	@echo "  Grafana:  http://localhost:3001"
	@echo "  Prometheus: http://localhost:9090"
	@echo "  Jaeger:   http://localhost:16686"

compose-down: ## Stop all services
	@echo "$(YELLOW)Stopping WohnFair services...$(NC)"
	docker-compose down
	@echo "$(GREEN)Services stopped$(NC)"

compose-logs: ## Show logs from all services
	docker-compose logs -f

# Development environment
dev: ## Start development environment
	@echo "$(YELLOW)Starting development environment...$(NC)"
	@echo "$(YELLOW)Make sure you have the following running:$(NC)"
	@echo "  - PostgreSQL on :5432"
	@echo "  - Redis on :6379"
	@echo "  - Kafka on :9092"
	@echo ""
	@echo "$(YELLOW)Starting services in development mode...$(NC)"
	@cd services/orchestration && go run ./cmd/gateway &
	@cd services/fairrent && go run ./cmd/fairrentd &
	@cd services/zk-lease && cargo run &
	@cd frontend && npm run dev &
	@echo "$(GREEN)Development environment started$(NC)"
	@echo "$(YELLOW)Press Ctrl+C to stop all services$(NC)"

# Bootstrap the entire system
bootstrap: check-tools deps proto build test compose-up ## Complete system setup
	@echo "$(GREEN)🎉 WohnFair system is ready!$(NC)"
	@echo "$(YELLOW)Access your system at:$(NC)"
	@echo "  Frontend: http://localhost:3000"
	@echo "  Gateway:  http://localhost:8080"
	@echo "  Keycloak: http://localhost:8081 (admin/admin)"
	@echo "  Grafana:  http://localhost:3001 (admin/admin)"
	@echo "  Prometheus: http://localhost:9090"
	@echo "  Jaeger:   http://localhost:16686"

# Clean build artifacts
clean: ## Clean all build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf services/*/bin/
	@rm -rf services/*/dist/
	@rm -rf frontend/.next/
	@rm -rf frontend/out/
	@rm -rf services/zk-lease/target/
	@rm -rf services/ml/build/
	@rm -rf services/ml/*.egg-info/
	@echo "$(GREEN)Build artifacts cleaned$(NC)"

# Performance testing
bench: ## Run performance benchmarks
	@echo "$(YELLOW)Running Go benchmarks...$(NC)"
	@cd services/fairrent && go test -bench=. -benchmem ./...
	@echo "$(YELLOW)Running Rust benchmarks...$(NC)"
	@cd services/zk-lease && cargo bench
	@echo "$(GREEN)Benchmarks completed$(NC)"

# Load testing
load-test: ## Run load tests with k6
	@echo "$(YELLOW)Running load tests...$(NC)"
	@cd benchmark/load/k6 && k6 run basic.js
	@echo "$(GREEN)Load tests completed$(NC)"

# Security scanning
security-scan: ## Run security scans
	@echo "$(YELLOW)Running Go security scan...$(NC)"
	@cd services/orchestration && gosec ./...
	@cd services/fairrent && gosec ./...
	@cd services/policy-dsl && gosec ./...
	@cd services/notifications && gosec ./...
	@echo "$(YELLOW)Running dependency vulnerability scan...$(NC)"
	@npm audit --prefix frontend
	@echo "$(GREEN)Security scans completed$(NC)"

# Database operations
db-migrate: ## Run database migrations
	@echo "$(YELLOW)Running database migrations...$(NC)"
	@cd services/orchestration && go run ./cmd/gateway migrate
	@echo "$(GREEN)Database migrations completed$(NC)"

db-seed: ## Seed database with sample data
	@echo "$(YELLOW)Seeding database...$(NC)"
	@cd services/orchestration && go run ./cmd/gateway seed
	@echo "$(GREEN)Database seeded$(NC)"

# Monitoring and debugging
logs: ## Show logs from specific service
	@echo "$(YELLOW)Available services:$(NC)"
	@echo "  gateway, fairrent, zk-lease, notifier, postgres, redis, kafka, keycloak, grafana, prometheus, jaeger"
	@echo "$(YELLOW)Usage: make logs SERVICE=gateway$(NC)"
	@if [ -n "$(SERVICE)" ]; then docker-compose logs -f $(SERVICE); else echo "$(RED)Please specify SERVICE=service_name$(NC)"; fi

status: ## Show status of all services
	@echo "$(YELLOW)Service Status:$(NC)"
	@docker-compose ps

# Documentation
docs-serve: ## Serve documentation locally
	@echo "$(YELLOW)Starting documentation server...$(NC)"
	@cd docs && python3 -m http.server 8000
	@echo "$(GREEN)Documentation available at http://localhost:8000$(NC)"

# Pre-commit hooks
install-hooks: ## Install pre-commit hooks
	@echo "$(YELLOW)Installing pre-commit hooks...$(NC)"
	@cp scripts/precommit.sh .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "$(GREEN)Pre-commit hooks installed$(NC)"

# All-in-one command
all: clean deps proto build test ## Clean, install deps, build, and test everything
	@echo "$(GREEN)🎉 All tasks completed successfully!$(NC)"

# Default target
.DEFAULT_GOAL := help
