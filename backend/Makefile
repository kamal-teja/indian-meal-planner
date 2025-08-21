# Makefile for Go Backend

.PHONY: build run test clean dev install lint format deps

# Variables
BINARY_NAME=meal-planner-server
BUILD_DIR=build
GO_FILES=$(shell find . -name "*.go" -type f)

# Default target
all: build

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Install development tools
install:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/air-verse/air@latest

# Build the application
build: deps
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/server/main.go

# Build for production with optimizations
build-prod: deps
	@echo "Building for production..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME) cmd/server/main.go

# Run the application
run: build
	@echo "Running application..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run in development mode with hot reload
dev:
	@echo "Running in development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not installed. Install with: make install"; \
		echo "Running without hot reload..."; \
		go run cmd/server/main.go; \
	fi

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench=. ./...

# Lint the code
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: make install"; \
	fi

# Format the code
format:
	@echo "Formatting code..."
	go fmt ./...
	go vet ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# Check for security vulnerabilities
security:
	@echo "Checking for security vulnerabilities..."
	go list -json -m all | nancy sleuth

# Generate documentation
docs:
	@echo "Generating documentation..."
	@if command -v godoc > /dev/null; then \
		echo "Starting godoc server at http://localhost:6060"; \
		godoc -http=:6060; \
	else \
		echo "godoc not installed. Install with: go install golang.org/x/tools/cmd/godoc@latest"; \
	fi

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t meal-planner-backend-go .

docker-run:
	@echo "Running Docker container..."
	docker run -p 5000:5000 --env-file .env meal-planner-backend-go

# Database commands
db-seed:
	@echo "Seeding database..."
	go run cmd/server/main.go --seed-only

# Help
help:
	@echo "Available commands:"
	@echo "  deps          - Install dependencies"
	@echo "  install       - Install development tools"
	@echo "  build         - Build the application"
	@echo "  build-prod    - Build for production"
	@echo "  run           - Build and run the application"
	@echo "  dev           - Run in development mode with hot reload"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  bench         - Run benchmarks"
	@echo "  lint          - Run linter"
	@echo "  format        - Format code"
	@echo "  clean         - Clean build artifacts"
	@echo "  security      - Check for security vulnerabilities"
	@echo "  docs          - Generate documentation"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  help          - Show this help message"
