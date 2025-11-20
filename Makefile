.PHONY: all build run clean build-backend run-backend dev-backend build-frontend dev-frontend test

# --- Project Configuration ---
PROJECT_NAME := jing-sync
GO_MAIN_DIR := ./cmd
GO_MAIN_FILE := main.go
BACKEND_BIN_DIR := ./bin
BACKEND_BIN_NAME := $(PROJECT_NAME)
BACKEND_OUTPUT := $(BACKEND_BIN_DIR)/$(BACKEND_BIN_NAME)

FRONTEND_DIR := ./frontend
FRONTEND_BUILD_DIR := $(FRONTEND_DIR)/dist # Vite 默认输出到 dist 目录
FRONTEND_OUTPUT_TO := ./web # 将前端构建结果复制到后端可访问的路径

# Go related commands
GO := go
GO_BUILD_FLAGS := -ldflags="-s -w" # Strip debug info and DWARF tables

# Node/NPM related commands
NPM := pnpm

# --- Default Target ---
all: build

# --- Backend Targets ---

# Build the Go backend application
build-backend:
	@echo "Building Go backend..."
	@mkdir -p $(BACKEND_BIN_DIR)
	$(GO) build $(GO_BUILD_FLAGS) -o $(BACKEND_OUTPUT) $(GO_MAIN_DIR)/$(GO_MAIN_FILE)
	@echo "Backend built to $(BACKEND_OUTPUT)"

# Run the Go backend application (requires it to be built)
run-backend: build-backend
	@echo "Running Go backend..."
	$(BACKEND_OUTPUT)

# For development, allowing hot-reloading or restarting for backend.
# This usually involves a tool like `air` or `go run` directly.
# For simplicity, we'll just use `go run` here, you might want to replace this with `air` or similar.
dev-backend:
	@echo "Starting Go backend in development mode..."
	$(GO) run $(GO_MAIN_DIR)/$(GO_MAIN_FILE)

db: dev-backend

# Run Go tests
test:
	@echo "Running Go tests..."
	$(GO) test ./...

# --- Frontend Targets ---

# Install frontend dependencies
install-frontend-deps:
	@echo "Installing frontend dependencies..."
	@cd $(FRONTEND_DIR) && $(NPM) install

# Build the frontend application using Vite
build-frontend: install-frontend-deps
	@echo "Building frontend with Vite..."
	@cd $(FRONTEND_DIR) && $(NPM) run build
	@echo "Frontend built to $(FRONTEND_BUILD_DIR)"
	@echo "Copying frontend build to $(FRONTEND_OUTPUT_TO)"
	@rm -rf $(FRONTEND_OUTPUT_TO) # Clean existing dir
	@mkdir -p $(FRONTEND_OUTPUT_TO) # Create target dir
	@cp -R $(FRONTEND_BUILD_DIR)/* $(FRONTEND_OUTPUT_TO)/ # Copy contents

# Start the Vite development server with hot-reloading
dev-frontend: install-frontend-deps
	@echo "Starting Vite development server..."
	@cd $(FRONTEND_DIR) && $(NPM) run dev

df: dev-frontend

# --- Combined Targets ---

# Full build: frontend + backend
build: build-frontend build-backend
	@echo "Full build complete."

# Run the entire application (backend from bin, frontend typically served by backend)
run: build run-backend
	@echo "Application started. Frontend static files are served by the backend."

# Start both backend and frontend development servers.
# Note: This will run in two separate processes. You might need to open two terminals.
# For single-command execution, consider using `concurrently` or similar tools.
# Here, we'll just show the commands.
dev:
	@echo "Starting backend and frontend development servers (may require separate terminals)..."
	@echo "In one terminal, run: make dev-backend"
	@echo "In another terminal, run: make dev-frontend"

# --- Clean Target ---
clean:
	@echo "Cleaning generated files..."
	@rm -rf $(BACKEND_BIN_DIR)
	@rm -rf $(FRONTEND_DIR)/node_modules
	@rm -rf $(FRONTEND_BUILD_DIR)
	@rm -rf $(FRONTEND_OUTPUT_TO)
	@$(GO) clean
	@echo "Clean complete."

.DEFAULT_GOAL := all
