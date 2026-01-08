.PHONY: help test test-short coverage coverage-report mocks mocks-clean mocks-regenerate build run docker-build docker-up docker-down swagger clean clean-all dev test-all ci

# Variables
APP_NAME=tc-fiap-order
MAIN_PATH=./cmd/api
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Detect OS
ifeq ($(OS),Windows_NT)
	BINARY_EXT=.exe
	RM=cmd /C del /Q /F
	RMDIR=cmd /C rd /S /Q
	OPEN=start
	COVERAGE_SCRIPT=powershell -ExecutionPolicy Bypass -File scripts/coverage.ps1
else
	BINARY_EXT=
	RM=rm -f
	RMDIR=rm -rf
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Darwin)
		OPEN=open
	else
		OPEN=xdg-open
	endif
	COVERAGE_SCRIPT=./scripts/coverage.sh
endif

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  help                 Show this help message"
	@echo "  test                 Run all tests"
	@echo "  test-short           Run tests without verbose output"
	@echo "  coverage             Run tests with coverage"
	@echo "  coverage-report      Generate and open coverage report"
	@echo "  mocks                Generate mocks using mockery"
	@echo "  mocks-clean          Clean generated mocks"
	@echo "  mocks-regenerate     Clean and regenerate all mocks"
	@echo "  build                Build the application"
	@echo "  run                  Run the application"
	@echo "  docker-build         Build Docker image"
	@echo "  docker-up            Start services with docker-compose"
	@echo "  docker-down          Stop services with docker-compose"
	@echo "  swagger              Generate Swagger documentation"
	@echo "  clean                Clean build artifacts and coverage files"
	@echo "  clean-all            Clean everything including mocks"
	@echo "  dev                  Start development environment"
	@echo "  test-all             Run mocks generation, tests and coverage"
	@echo "  ci                   Run CI pipeline (tidy, mocks, test, build)"

# Testing
test: ## Run all tests
	@echo "Running tests..."
	go test ./... -v -race

test-short: ## Run tests without verbose output
	@echo "Running tests..."
	go test ./...

coverage: ## Run tests with coverage (uses OS-specific script)
	@echo "Running tests with coverage..."
	$(COVERAGE_SCRIPT)

coverage-report: coverage ## Generate and open coverage report
	@echo "Opening coverage report..."
	$(OPEN) $(COVERAGE_HTML)

# Mocking
mocks: ## Generate mocks using mockery
	@echo "Generating mocks..."
	@if ! command -v mockery &> /dev/null; then \
		echo "Installing mockery..."; \
		go install github.com/vektra/mockery/v2@latest; \
	fi
	mockery

mocks-clean: ## Clean generated mocks
	@echo "Cleaning mocks..."
	$(RMDIR) mocks 2>nul || true

mocks-regenerate: mocks-clean mocks ## Clean and regenerate all mocks

# Building
build: ## Build the application
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME)$(BINARY_EXT) $(MAIN_PATH)
	@echo "Build complete: bin/$(APP_NAME)$(BINARY_EXT)"

run: ## Run the application
	@echo "Running $(APP_NAME)..."
	go run $(MAIN_PATH)/main.go

# Docker
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -f Dockerfile -t $(APP_NAME):latest .

docker-up: ## Start services with docker-compose
	@echo "Starting services..."
	docker-compose up -d

docker-down: ## Stop services with docker-compose
	@echo "Stopping services..."
	docker-compose down

# Swagger
swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	@if ! command -v swag &> /dev/null; then \
		echo "Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	swag init -g $(MAIN_PATH)/main.go -o ./docs
	@echo "Swagger docs generated in ./docs"

# Dependencies
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download

deps-tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	go mod tidy

# Clean
clean: ## Clean build artifacts and coverage files
	@echo "Cleaning..."
	$(RMDIR) bin 2>nul || true
	$(RM) $(COVERAGE_FILE) 2>nul || true
	$(RM) $(COVERAGE_HTML) 2>nul || true
	@echo "Clean complete"

clean-all: clean mocks-clean ## Clean everything including mocks

# Development workflow
dev: docker-up run ## Start development environment

test-all: mocks test coverage ## Run mocks generation, tests and coverage

ci: deps-tidy mocks test build ## Run CI pipeline (tidy, mocks, test, build)

# Quick aliases
t: test ## Alias for test
tc: coverage ## Alias for coverage
b: build ## Alias for build
r: run ## Alias for run
m: mocks ## Alias for mocks
