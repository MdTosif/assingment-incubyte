# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
BINARY_NAME=salary-management
BINARY_UNIX=$(BINARY_NAME)_unix

# Test coverage
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

.PHONY: all build clean test coverage coverage-html deps lint help

all: test build

build: 
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server

test:
	$(GOTEST) -v ./...

test-verbose:
	$(GOTEST) -v -race ./...

test-coverage:
	$(GOTEST) -v -coverprofile=$(COVERAGE_FILE) ./...

coverage-html: test-coverage
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

test-unit:
	$(GOTEST) -v ./internal/models/... ./internal/handlers/... ./internal/database/...

test-integration:
	$(GOTEST) -v ./cmd/server/...

test-specific:
	$(GOTEST) -v -run $(RUN) ./...

deps:
	$(GOMOD) download
	$(GOMOD) tidy

lint:
	golangci-lint run

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(COVERAGE_FILE)
	rm -f $(COVERAGE_HTML)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server
	./$(BINARY_NAME)

# Development targets
dev:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server
	PORT=8080 ./$(BINARY_NAME)

# Database targets
seed:
	$(GOBUILD) -o seed -v ./cmd/seed
	./seed
	rm -f seed

# Docker targets
docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run -p 8080:8080 $(BINARY_NAME)

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/server

# TDD specific targets
tdd-watch:
	gow -run "$(GOTEST) ./..."

tdd-models:
	$(GOTEST) -v ./internal/models/...

tdd-handlers:
	$(GOTEST) -v ./internal/handlers/...

tdd-integration:
	$(GOTEST) -v ./cmd/server/...

# Benchmark tests
benchmark:
	$(GOTEST) -bench=. -benchmem ./...

# Help target
help:
	@echo "Available targets:"
	@echo "  all              - Run tests and build"
	@echo "  build            - Build the application"
	@echo "  test             - Run all tests"
	@echo "  test-verbose     - Run tests with verbose output and race detection"
	@echo "  test-coverage    - Run tests with coverage"
	@echo "  coverage-html    - Generate HTML coverage report"
	@echo "  test-unit        - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-specific    - Run specific tests (use RUN variable)"
	@echo "  deps             - Download and tidy dependencies"
	@echo "  lint             - Run linter"
	@echo "  clean            - Clean build artifacts"
	@echo "  run              - Build and run the application"
	@echo "  dev              - Run in development mode"
	@echo "  seed             - Seed the database"
	@echo "  docker-build     - Build Docker image"
	@echo "  docker-run       - Run Docker container"
	@echo "  build-linux      - Cross compile for Linux"
	@echo "  tdd-watch        - Watch for changes and run tests"
	@echo "  tdd-models       - Run model tests only"
	@echo "  tdd-handlers     - Run handler tests only"
	@echo "  tdd-integration - Run integration tests only"
	@echo "  benchmark        - Run benchmark tests"
	@echo "  help             - Show this help message"
