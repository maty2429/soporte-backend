# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

```bash
# Run locally (uses .env.local)
make run

# Run tests
make test

# Format and lint
make fmt
make lint

# Docker development
make docker-up          # Start dev stack
make docker-down        # Stop dev stack

# Docker production
make docker-up-prod     # Start prod stack
make docker-down-prod   # Stop prod stack

# Migrations
make migrate-status
make migrate-up
make migrate-down
```

### Running a single test

```bash
go test -v -run TestName ./internal/application/services/tests/
```

### Building binaries

```bash
make build-dev                    # Development binary
make build-prod                   # Optimized production binary (no SQLite, no Swagger)
```

## Architecture

This is a **hexagonal architecture** Go REST API template with clear layer separation:

```
core (Domain)           → Pure business entities and errors (no external deps)
  └── ports             → Repository interfaces (contracts)
       ↓
application/services    → Use cases and business logic
       ↓
adapters/repository     → GORM implementations of ports
       ↓
delivery/http           → Gin handlers, routes, middlewares, contracts (DTOs)
```

### Key architectural rules

- **Domain entities** (`internal/core/domain/`) must have NO Gin or GORM imports
- **Ports** (`internal/core/ports/`) define interfaces that adapters implement
- **Services** orchestrate business logic using port interfaces
- **Repositories** implement ports using GORM models
- **Handlers** convert HTTP ↔ domain, call services, return JSON

### Adding a new resource

Follow this order:
1. Domain entity in `internal/core/domain/`
2. Port interface in `internal/core/ports/`
3. Service types in `internal/application/services/types/`
4. Service implementation in `internal/application/services/`
5. GORM model in `internal/adapters/repository/models/`
6. Repository in `internal/adapters/repository/repository/`
7. HTTP contracts (DTOs) in `internal/delivery/http/contracts/`
8. Handler in `internal/delivery/http/handlers/`
9. Register routes in `internal/delivery/http/routes/router.go`
10. SQL migration in `db/migrations/`

## Testing Patterns

- Service tests use **in-memory SQLite**: `file:testname?mode=memory&cache=shared`
- Tests run in parallel with `t.Parallel()`
- Handler tests are in `internal/delivery/http/handlers/tests/`
- Service tests are in `internal/application/services/tests/`

## Configuration

Environment files load in order (earlier files don't overwrite already-set vars):

**Development:** `.env.local` → `.env.development` → `.env` → `.env.example`
**Production:** `.env.local` → `.env.production` → `.env`

Key variables:
- `DB_FAIL_FAST=false` - API starts in degraded mode if DB unavailable
- `DB_AUTO_MIGRATE=false` - Use SQL migrations, not GORM AutoMigrate
- `DOCS_ENABLED=true` - Enable Swagger UI (dev only)
- `-tags production` - Excludes SQLite and Swagger from binary

## Domain Errors

Use factory functions from `internal/core/domain/errors.go`:
- `domain.ValidationError(msg, err)` → 400
- `domain.NotFoundError(resource, err)` → 404
- `domain.ConflictError(msg, err)` → 409
- `domain.InternalError(msg, err)` → 500
- `domain.ServiceUnavailableError(msg, err)` → 503

## Linting

The project uses `golangci-lint` with security-focused linters:
- **gosec** - Security vulnerabilities
- **bodyclose** - HTTP body close checks
- **errcheck** - Unchecked errors
- **staticcheck** - Advanced static analysis

Run locally: `make lint` or `golangci-lint run`

## CI Pipeline

GitHub Actions runs on push/PR:
1. **test** - `go test` with race detection and coverage
2. **lint** - `golangci-lint`
3. **security** - `govulncheck` + `gosec`
4. **docker-build** - Build + Trivy vulnerability scan

## Endpoints

- `GET /health` - Overall health
- `GET /livez` - Liveness probe
- `GET /readyz` - Readiness probe (checks DB)
- `GET /metrics` - Prometheus metrics
