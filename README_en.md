# Jing-Sync

[English](README.me) | [简体中文](README.md)

Jing-Sync is a modern, high-performance, full-stack file synchronization and management web application. It connects with cloud storage engines such as OpenList and AList, providing efficient cross-directory and cross-drive file synchronization, cron job scheduling, mirror alignment, file filtering, and detailed audit logging.

---

## 🌟 Key Features

- **Multi-Engine Management**: Seamless integration with OpenList / AList instances. Configure server URLs and authentication tokens with built-in one-click connection testing.
- **Intuitive Path Selection**:
  - Manually input or paste source and destination paths.
  - Interactive modal Path Browser for visual tree navigation, directory traversal, and rapid folder picking.
- **Flexible Synchronization Modes**:
  - **Add-Only (Incremental)**: Safely copies only missing or updated files to the destination without overwriting or deleting existing destination files.
  - **Full Synchronization (Mirror)**: Accurately mirrors destination with source; obsolete files deleted on source are automatically purged from destination.
  - **Move Files (Cut)**: Transfers files to destination and removes them from the source upon successful copy.
- **Granular Filtering & Safety Protections**:
  - Glob pattern exclusions (e.g., `*.tmp, *.bak, .git/**`).
  - File size threshold filtering (configurable minimum and maximum bytes).
  - Built-in source/destination path overlap detection to prevent infinite recursion loops.
- **Robust Cron Job Scheduling**:
  - Intuitive minute-interval presets or standard 5-field cron syntax (e.g., `0 2 * * *`).
  - Dynamic `CronManager`: Automatically recalibrates and hot-reloads schedules on job creation, modification, toggle, or deletion without restarting the server.
- **Real-Time Audit & Terminal Drawer**:
  - Comprehensive sync metrics tracking: discovered, added, deleted, skipped, and failed files with error details.
  - Sleek terminal-styled drawer for inspecting sync logs and troubleshooting.
- **Security & User Management**:
  - Stateless JWT-based authentication.
  - Auto-generated strong admin password printed on initial startup.
  - User management interface with masked passwords and credential rotation.
  - Multi-language i18n support (`zh-CN` Simplified Chinese and `en` English) and Dark/Light theme switching.

---

## 🛠 Architecture & Tech Stack

### Project Structure

```text
jing-sync/
├── cmd/main.go               # Application entry point
├── boot/                     # Initialization modules (config, logger, i18n, DB, Gin)
├── api/                      # API layer (routes, controllers, middlewares: auth, timeout, logger)
├── internal/                 # Core domain logic
│   ├── models/               # GORM data models (BaseModel, User, Engine, Job, SyncLog)
│   ├── services/             # Business services (SyncService, OpenListClient, CronManager, BaseService)
│   └── utils/                # Utilities (JWT tokens, HTTP requests, file helpers, responses)
├── frontend/                 # Frontend project (Vue 3 + TypeScript + Vite + Element Plus + Tailwind CSS)
└── public/                   # Embedded static assets (frontend dist & i18n toml via Go embed)
```

### Technology Highlights

- **Backend**:
  - Language: Go 1.23+
  - Web Framework: Gin
  - ORM & Persistence: GORM with SQLite (embedded database with automatic schema migrations)
  - Scheduler: `robfig/cron/v3`
  - Logging: Logrus with daily rotation
  - Authentication: `golang-jwt/v5` with salted SHA-256 password hashing
- **Frontend**:
  - Framework: Vue 3 (Composition API / `<script setup>`) + TypeScript
  - Bundler: Vite
  - UI Toolkit: Element Plus + Tailwind CSS 4.x
  - State Management: Pinia
  - Routing & Localization: Vue Router + Vue I18n

---

## 🚀 Quick Start & Build Instructions

The project provides an optimized, cross-platform `Makefile` compatible with both Unix (macOS / Linux) and Windows (PowerShell).

### Command Cheatsheet

```bash
# View command help and available targets
make help

# 1. Full Stack Production Build (builds Vue frontend and embeds it into the Go binary)
make build        # or: make all

# 2. Run the compiled binary
./bin/jing-sync    # macOS / Linux
# .\bin\jing-sync.exe # Windows

# 3. Development Mode
make dev-frontend # Starts Vite dev server (alias: make df, default port: 5173)
make dev-backend  # Starts Go backend in development mode (alias: make db, default port: 8888)

# 4. Standalone Builds
make build-frontend # Compiles frontend and syncs to public/web (alias: make bf)
make build-backend  # Compiles Go backend binary to bin/ (alias: make bb)

# 5. Clean Workspace
make clean        # Cleans bin/, dist/, and public/web/
make clean-all    # Deep clean (cleans build artifacts, node_modules, and logs)
```

---

## 🧪 Testing & Quality Assurance

Both backend and frontend are equipped with complete automated test suites:

### Backend Unit Tests (Go Test)
Covers middlewares, database services, route tree registrations, OpenList client operations, and synchronization logic. Uses an in-memory HTTP transport (`HandlerTransport`) that executes seamlessly in network-restricted sandbox environments:

```bash
# Run all Go backend unit tests
make test

# Run race condition detection
make test-race
```

### Frontend Unit Tests (Vitest)
Covers utility libraries, Pinia stores, and Vue view components:

```bash
# Run frontend unit tests via Vitest
make test-frontend
```

### Full Test Suite
```bash
# Run both backend and frontend tests in sequence
make test-all
```

---

## ⚙️ Configuration & Secrets

Runtime configurations and data files reside in `data/`:
- `data/config_prod.json`: Production configuration (port, DB filename, request timeout).
- `data/secret.key`: Secret key used for signing JWT tokens (automatically generated on first run).
- `data/jing-sync.db`: SQLite database file.
- `data/logs/`: Application execution and HTTP access logs.

> **Note**: On the very first run, the system automatically initializes an `admin` user with a randomly generated high-entropy password printed to the terminal console. Be sure to save it and update your password via the User Settings page after logging in.

---

## 🐳 Docker Deployment

A multi-stage `Dockerfile` is provided for containerized deployment:

```bash
# 1. Build Docker image using Makefile
make docker-build
# Or via docker CLI: docker build -t jing-sync .

# 2. Run the container with persistent volume mounting
make docker-run
# Or via docker CLI:
docker run -d \
  --name jing-sync \
  -p 8888:8888 \
  -v $(pwd)/data:/app/data \
  --restart unless-stopped \
  jing-sync
```

Web interface will be accessible at: `http://127.0.0.1:8888/`

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).
