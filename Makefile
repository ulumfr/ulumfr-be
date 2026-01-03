.PHONY: build run dev clean generate migrate help

# Binary name
BINARY_NAME=ulumfr-be

# Build directory
BUILD_DIR=bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Prisma
PRISMA=go run github.com/steebchen/prisma-client-go

help: ## Display this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -ldflags="-s -w" ./cmd/api

run: build ## Build and run the application
	@echo "Running..."
	./$(BUILD_DIR)/$(BINARY_NAME)

dev: ## Run with hot reload (requires air)
	@echo "Running in development mode..."
	air

clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

generate: ## Generate Prisma client
	@echo "Generating Prisma client..."
	$(PRISMA) generate

migrate: ## Push schema to database
	@echo "Pushing schema to database..."
	$(PRISMA) db push

migrate-dev: ## Create development migration
	@echo "Creating migration..."
	$(PRISMA) migrate dev

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOCMD) vet ./...

# Production build with optimizations
build-prod: ## Build for production
	@echo "Building for production..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -ldflags="-s -w" ./cmd/api

# Install development tools
tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
