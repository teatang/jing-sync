# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

jing-sync is a full-stack file synchronization web application with:
- **Backend**: Go 1.23+ with Gin framework, GORM ORM, SQLite database
- **Frontend**: Vue 3 + TypeScript with Vite, Element Plus, Tailwind CSS 4.x
- **Authentication**: JWT-based with auto-generated admin user on first run

## Build Commands

```bash
# Full production build (frontend + backend)
make build        # or make all

# Run the built binary
./bin/jing-sync

# Development
make dev-frontend # or make df  - Start Vite dev server (port 5173)
make dev-backend  # or make db  - Run Go backend with go run

# Build only
make build-frontend  # Build Vue and copy to public/web
make build-backend   # Build Go binary to bin/jing-sync

# Clean build artifacts
make clean
```

## Architecture

### Backend Layer Structure

```
cmd/main.go → boot/ (config, logger, i18n, database, app)
    → api/routes → api/controllers → internal/services → internal/models
```

**Key directories**:
- `boot/` - Initialization (Gin server, config, SQLite, cron scheduler)
- `api/controllers/` - HTTP request handlers (User, Engine, Job, OpenList)
- `api/middlewares/` - Auth, Logger, Timeout, i18n
- `api/routes/` - Route definitions
- `internal/services/` - Business logic with `BaseService[T]` generic pattern
- `internal/models/` - GORM models with shared `BaseModel` (ID, CreateTime, UpdateTime, Status)

### Frontend Structure

```
frontend/src/
├── views/        # Page components (cron-job, engine, setting)
├── stores/       # Pinia stores
├── router/       # Vue Router config
├── i18n/         # vue-i18n locales
└── utils/        # REST client, token helpers
```

### Data Flow

1. Request → Gin router (`boot/app/app.go`)
2. Middleware stack: i18n → Timeout → Logger → Auth (JWT)
3. Controller receives request
4. Controller → Service → Model/DB via GORM
5. Response via `internal/utils/response` helpers

## Configuration

Configuration files in `data/`:
- `config_prod.json` - Production settings (port 8888, site name, DB name)
- `secret.key` - JWT signing key

Environment types: `prod`, `dev`, `test`, `unit_test`

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `POST /api/login` | User authentication |
| `POST/GET/PUT/DELETE /api/user` | User management |
| `POST/GET/PUT/DELETE /api/engine` | Sync engine configuration |
| `POST/GET/PUT/DELETE /api/job` | Cron job management |
| `GET /api/open-list` | List sync operations |

Frontend routes: `/` (cron-job), `/engine`, `/user-setting`

## Key Patterns

- **Generic Base Service**: `internal/services/` uses `BaseService[T]` for reusable CRUD operations
- **Database**: SQLite with auto-migration on startup
- **Cron Jobs**: Uses `robfig/cron` for scheduling sync jobs
- **Logging**: Logrus with log rotation, stored in `data/logs/`

## Docker

```bash
docker build -t jing-sync .
docker run -p 8888:8888 -v /path/to/data:/app/data jing-sync
```

Multi-stage build: Node.js frontend → Go backend → Alpine final image.
