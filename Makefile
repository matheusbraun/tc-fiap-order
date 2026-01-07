.PHONY: help build run test test-unit test-bdd test-coverage mocks clean dev docker-build docker-up docker-down

# Help
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run all tests"
	@echo "  make test-unit      - Run unit tests only"
	@echo "  make test-bdd       - Run BDD tests only"
	@echo "  make test-coverage  - Generate test coverage report"
	@echo "  make mocks          - Generate mocks using mockery"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make dev            - Run with hot reload (requires air)"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-up      - Start services with docker-compose"
	@echo "  make docker-down    - Stop services with docker-compose"

# Build
build:
	@echo "Building application..."
	go build -o bin/order-service ./cmd/api

# Run
run:
	@echo "Running application..."
	go run ./cmd/api/main.go

# Generate mocks
mocks:
	@echo "Generating mocks..."
	@if ! command -v mockery &> /dev/null; then \
		echo "Installing mockery..."; \
		go install github.com/vektra/mockery/v2@latest; \
	fi
	mockery --all

# Run all tests
test:
	@echo "Running all tests..."
	go test ./... -v -race

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	go test ./tests/unit/... -v -race

# Run BDD tests
test-bdd:
	@echo "Running BDD tests..."
	go test ./tests -v

# Generate coverage report
test-coverage:
	@echo "Generating coverage report..."
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@echo "Coverage: $$(go tool cover -func=coverage.out | grep total | awk '{print $$3}')"

# Clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Development with hot reload (requires air: go install github.com/air-verse/air@latest)
dev:
	@if ! command -v air &> /dev/null; then \
		echo "Installing air..."; \
		go install github.com/air-verse/air@latest; \
	fi
	air

# Docker
docker-build:
	@echo "Building Docker image..."
	docker build -f Dockerfile -t order-service:latest .

docker-up:
	@echo "Starting services..."
	docker-compose up -d

docker-down:
	@echo "Stopping services..."
	docker-compose down
