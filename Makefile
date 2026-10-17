.PHONY: help run test build docker-build docker-run clean

# Variables
APP_NAME=devops-demo-app
VERSION?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
DOCKER_IMAGE?=$(APP_NAME):$(VERSION)

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application locally
	@echo "🚀 Starting application..."
	@go run main.go

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v -cover

test-coverage: ## Run tests with coverage report
	@echo "🧪 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

build: ## Build the Go binary
	@echo "🔨 Building binary..."
	@CGO_ENABLED=0 go build \
		-ldflags="-w -s -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)" \
		-o $(APP_NAME) \
		main.go
	@echo "✅ Binary built: ./$(APP_NAME)"

docker-build: ## Build Docker image
	@echo "🐋 Building Docker image..."
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		-t $(DOCKER_IMAGE) .
	@echo "✅ Docker image built: $(DOCKER_IMAGE)"

docker-run: ## Run Docker container locally
	@echo "🐋 Running Docker container..."
	@docker run -p 8080:8080 \
		-e ENVIRONMENT=local \
		-e POD_NAME=local-docker \
		-e POD_NAMESPACE=default \
		$(DOCKER_IMAGE)

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@rm -f $(APP_NAME)
	@rm -f coverage.out coverage.html
	@echo "✅ Clean complete"

fmt: ## Format Go code
	@echo "📝 Formatting code..."
	@go fmt ./...
	@echo "✅ Format complete"

lint: ## Run go vet
	@echo "🔍 Running linter..."
	@go vet ./...
	@echo "✅ Lint complete"

.DEFAULT_GOAL := help
