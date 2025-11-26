.PHONY: all build run clean build-backend run-backend dev-backend build-frontend dev-frontend test

# --- Project Configuration ---
PROJECT_NAME := jing-sync
GO := go
GO_BUILD_FLAGS := -ldflags="-s -w" # Strip debug info and DWARF tables
GO_MAIN_FILE := main.go
NPM := pnpm 

GO_MAIN_DIR := ./cmd
BACKEND_BIN_DIR := ./bin
FRONTEND_DIR := ./frontend
FRONTEND_OUTPUT_TO := ./public/web

ifeq ($(OS),Windows_NT)
    PLATFORM = windows
    BACKEND_BIN_NAME = $(PROJECT_NAME).exe
else
    PLATFORM = other
	BACKEND_BIN_NAME = $(PROJECT_NAME)
endif

# 定义 PowerShell 命令前缀，包含必要的参数
POWERSHELL_PREFIX := powershell.exe -NoProfile -Command

# Remove a directory recursively
# $1: Path to remove
define RMDIR
	@echo "RMDIR start $(1)"
	$(if $(filter windows,$(PLATFORM)), \
		@$(POWERSHELL_PREFIX) "& { if (Test-Path -Path \"$1\" -PathType Container) { Remove-Item -Recurse -Force -Path \"$1\" -ErrorAction SilentlyContinue | Out-Null } }"
	)
	$(if $(filter other,$(PLATFORM)), \
		@rm -rf "$(1)"
	)
	@echo "RMDIR end $(1)"
endef

# Create a directory (if not exists)
# $1: Path to create
define MKDIR
	@echo "MKDIR start $(1)"
	$(if $(filter windows,$(PLATFORM)), \
		@$(POWERSHELL_PREFIX) "& { New-Item -ItemType Directory -Path \"$1\" -ErrorAction SilentlyContinue | Out-Null }"
	)
	$(if $(filter other,$(PLATFORM)), \
		@mkdir -p "$(1)"
	)
	@echo "MKDIR end $(1)"
endef

# Copy contents of one directory to another
# '$1': source directory (e.g., frontend/dist)
# '$2': destination directory (e.g., web)
define COPY_CONTENTS
	@echo "COPY_CONTENTS start $(1)/* to $(2)"
	$(if $(filter windows,$(PLATFORM)), \
		@$(POWERSHELL_PREFIX) "& { Copy-Item -Path \"$1\\*\" -Destination \"$2\" -Recurse -Force -ErrorAction SilentlyContinue | Out-Null }"
	)
	$(if $(filter other,$(PLATFORM)), \
		@cp -r "$(1)/" "$(2)"
	)
	@echo "COPY_CONTENTS end $(1)/* to $(2)"
endef

BACKEND_OUTPUT := $(BACKEND_BIN_DIR)/$(BACKEND_BIN_NAME)
 # Vite 默认输出到 dist 目录
FRONTEND_BUILD_DIR := $(FRONTEND_DIR)/dist

# --- Default Target ---
all: build

# --- Frontend Targets ---
# Install frontend dependencies
frontend-deps:
	@echo "Installing frontend dependencies with $(NPM)..."
	@cd $(FRONTEND_DIR) && $(NPM) install
# Build the frontend application using Vite
build-frontend: frontend-deps
	@echo "Building frontend with Vite..."
# 	@cd $(FRONTEND_DIR) && $(NPM) run build
	@echo "Frontend built to $(FRONTEND_BUILD_DIR)"
	@echo "Copying frontend build to $(FRONTEND_OUTPUT_TO)"
	$(call RMDIR,$(FRONTEND_OUTPUT_TO))
	$(call MKDIR,$(FRONTEND_OUTPUT_TO))
	$(call COPY_CONTENTS,$(FRONTEND_BUILD_DIR),$(FRONTEND_OUTPUT_TO))
	@echo "Frontend build copied to $(FRONTEND_OUTPUT_TO)."

# Start the Vite development server with hot-reloading
dev-frontend: frontend-deps
	@echo "Starting Vite development server..."
	@cd $(FRONTEND_DIR) && $(NPM) run dev

df: dev-frontend

# Build the Go backend application
build-backend:
	@echo "Building Go backend..."
	$(call RMDIR,$(BACKEND_BIN_DIR))
	$(call MKDIR,$(BACKEND_BIN_DIR))
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
# test:
# 	@echo "Running Go tests..."
# 	$(GO) test ./...

# --- Combined Targets ---

# Full build: frontend + backend
build: build-frontend build-backend
	@echo "Full build complete."

# --- Clean Target ---
clean:
	@echo "Cleaning generated files..."
	$(call RMDIR,$(BACKEND_BIN_DIR))
	$(call RMDIR,$(FRONTEND_BUILD_DIR))
	$(call RMDIR,$(FRONTEND_OUTPUT_TO))
# 	@rm -rf $(BACKEND_BIN_DIR)
# 	@rm -rf $(FRONTEND_DIR)/node_modules
# 	@rm -rf $(FRONTEND_BUILD_DIR)
# 	@rm -rf $(FRONTEND_OUTPUT_TO)
	@$(GO) clean
	@echo "Clean complete."

.DEFAULT_GOAL := all
