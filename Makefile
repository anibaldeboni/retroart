# RetroArt TrimUI Smart Pro Makefile
# ===================================

.PHONY: help build run clean shell logs docker-build docker-clean native-build native-run

# Default target
.DEFAULT_GOAL := help

# Project configuration
PROJECT_NAME := retroart
BIN_DIR := bin
BINARY_NAME := $(BIN_DIR)/retroart
TRIMUI_BINARY_NAME := $(BIN_DIR)/retroart-trimui-arm64
GO_CMD := go
DOCKER_COMPOSE := docker compose

# Colors for output
GREEN := \033[0;32m
BLUE := \033[0;34m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

help: ## Show this help message
	@echo "$(BLUE)🐳 RetroArt TrimUI Smart Pro Build System$(NC)"
	@echo "$(BLUE)==========================================$(NC)"
	@echo ""
	@echo "$(GREEN)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: docker-build ## Build ARM64 binary for TrimUI Smart Pro using Docker
	@echo "$(GREEN)🏗️  Building TrimUI Smart Pro cross-compilation...$(NC)"
	@echo ""
	@echo "$(BLUE)📦 Preparing Docker TrimUI SDL2 environment...$(NC)"
	$(DOCKER_COMPOSE) build retroart-build
	@echo ""
	@echo "$(BLUE)🔧 Executing TrimUI SDL2 cross-compilation...$(NC)"
	@mkdir -p $(BIN_DIR)
	$(DOCKER_COMPOSE) run --rm retroart-build
	@if [ -f "$(TRIMUI_BINARY_NAME)" ]; then \
		echo ""; \
		echo "$(GREEN)✅ TrimUI build completed successfully!$(NC)"; \
		echo "$(BLUE)📁 TrimUI ARM64 binary created: $$(pwd)/$(TRIMUI_BINARY_NAME)$(NC)"; \
		echo "$(BLUE)📊 File information:$(NC)"; \
		ls -la $(TRIMUI_BINARY_NAME); \
		file $(TRIMUI_BINARY_NAME); \
		echo ""; \
		echo "$(YELLOW)🎮 To transfer to TrimUI Smart Pro:$(NC)"; \
		echo "   scp $(TRIMUI_BINARY_NAME) user@trimui:/path/to/destination/"; \
	else \
		echo ""; \
		echo "$(RED)❌ Build failed - binary not found$(NC)"; \
		exit 1; \
	fi

native-build: ## Build native binary for local development
	@echo "$(GREEN)🔨 Building native binary...$(NC)"
	@mkdir -p $(BIN_DIR)
	$(GO_CMD) build -o $(BINARY_NAME) cmd/retroart/main.go
	@echo "$(GREEN)✅ Native build completed: $(BINARY_NAME)$(NC)"

run: native-build ## Build and run the application locally
	@echo "$(GREEN)🚀 Running RetroArt locally...$(NC)"
	./$(BINARY_NAME)

native-run: run ## Alias for 'run' target

shell: ## Open interactive shell in Docker container for debugging
	@echo "$(BLUE)🐚 Opening interactive TrimUI SDL2 shell...$(NC)"
	@echo "   Use 'mkdir -p bin && go build -o bin/retroart-trimui-arm64 cmd/retroart/main.go' to compile manually"
	@echo "   Use 'exit' to quit"
	@echo ""
	$(DOCKER_COMPOSE) build retroart-builder
	$(DOCKER_COMPOSE) run --rm retroart-builder

docker-build: ## Build Docker images
	@echo "$(BLUE)📦 Building Docker images...$(NC)"
	$(DOCKER_COMPOSE) build

clean: ## Clean up containers, images and binaries
	@echo "$(YELLOW)🧹 Cleaning up containers and images...$(NC)"
	$(DOCKER_COMPOSE) down --remove-orphans
	@echo "$(YELLOW)🗑️  Removing project-related Docker images...$(NC)"
	@docker images | grep retroart | awk '{print $$3}' | xargs -r docker rmi || true
	$(DOCKER_COMPOSE) down -v
	@echo "$(YELLOW)🗑️  Removing local binaries...$(NC)"
	rm -rf $(BIN_DIR)
	@echo "$(GREEN)✅ Cleanup completed$(NC)"

docker-clean: ## Clean only Docker resources (keep binaries)
	@echo "$(YELLOW)🧹 Cleaning Docker resources...$(NC)"
	$(DOCKER_COMPOSE) down --remove-orphans
	@docker images | grep retroart | awk '{print $$3}' | xargs -r docker rmi || true
	$(DOCKER_COMPOSE) down -v
	@echo "$(GREEN)✅ Docker cleanup completed$(NC)"

logs: ## Show logs from last Docker build
	@echo "$(BLUE)📋 Logs from last build:$(NC)"
	$(DOCKER_COMPOSE) logs retroart-build

test: ## Run tests
	@echo "$(GREEN)🧪 Running tests...$(NC)"
	$(GO_CMD) test ./...

fmt: ## Format Go code
	@echo "$(GREEN)🎨 Formatting Go code...$(NC)"
	$(GO_CMD) fmt ./...

vet: ## Run go vet
	@echo "$(GREEN)🔍 Running go vet...$(NC)"
	$(GO_CMD) vet ./...

mod-tidy: ## Tidy Go modules
	@echo "$(GREEN)📦 Tidying Go modules...$(NC)"
	$(GO_CMD) mod tidy

check: fmt vet test ## Run all checks (format, vet, test)
	@echo "$(GREEN)✅ All checks completed$(NC)"

# Development workflow
dev: clean native-build run ## Clean, build and run for development

# Production workflow  
prod: clean build ## Clean and build for production (TrimUI)

# Check prerequisites
check-docker: ## Check if Docker is running
	@if ! docker info > /dev/null 2>&1; then \
		echo "$(RED)❌ Docker is not running. Please start Docker and try again.$(NC)"; \
		exit 1; \
	else \
		echo "$(GREEN)✅ Docker is running$(NC)"; \
	fi

# Show project status
status: ## Show project status and file information
	@echo "$(BLUE)📊 Project Status$(NC)"
	@echo "$(BLUE)=================$(NC)"
	@echo ""
	@echo "$(YELLOW)Project:$(NC) $(PROJECT_NAME)"
	@echo "$(YELLOW)Native Binary:$(NC) $(BINARY_NAME)"
	@echo "$(YELLOW)TrimUI Binary:$(NC) $(TRIMUI_BINARY_NAME)"
	@echo ""
	@if [ -f "$(BINARY_NAME)" ]; then \
		echo "$(GREEN)✅ Native binary exists:$(NC)"; \
		ls -la $(BINARY_NAME); \
	else \
		echo "$(RED)❌ Native binary not found$(NC)"; \
	fi
	@echo ""
	@if [ -f "$(TRIMUI_BINARY_NAME)" ]; then \
		echo "$(GREEN)✅ TrimUI binary exists:$(NC)"; \
		ls -la $(TRIMUI_BINARY_NAME); \
		file $(TRIMUI_BINARY_NAME); \
	else \
		echo "$(RED)❌ TrimUI binary not found$(NC)"; \
	fi
