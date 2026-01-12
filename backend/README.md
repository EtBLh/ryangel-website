# Backend Development Guide

This folder contains the Go + Gin API for RyAngel Commerce. The goal is to keep local onboarding to a single command so contributors can focus on feature work quickly.

## Prerequisites

- Go 1.23+ (already installed via `golang-go` APT package).
- PostgreSQL instance reachable with the credentials stored in `.env`.
- Optional: `make` for running helper targets.

## Project layout

```
backend/
├── cmd/server          # main package entry point
├── internal/config     # env parsing and derived settings
├── internal/database   # pgx connection helpers
├── internal/http       # Gin handlers
├── internal/server     # router + middleware wiring
├── go.mod / go.sum     # module definition
└── README.md           # this file
```

## Environment variables

Values are read from the standard process environment and `.env` (loaded automatically at runtime). Ensure the following keys exist:

```
APP_HOST=0.0.0.0
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=ryangel
DB_PASSWORD=RyangelPa33word
DB_NAME=ryangel
DB_SSLMODE=disable
```

## Useful commands

```bash
# Run the API locally (hot reload supported via `go run`)
make run

# Execute unit tests (none yet, but wired up for future use)
make test

# Tidy dependencies
make tidy
```

If you prefer raw Go commands:

```bash
go run ./cmd/server
```

The server exposes:
- `GET /api/healthz` – liveness/readiness with PostgreSQL ping.

## Next steps

- Flesh out routers/controllers based on `api_specs.md`.
- Add migrations (e.g., using `golang-migrate`) once database automation is needed.
- Introduce structured logging (zap/zerolog) and request-level observability.
