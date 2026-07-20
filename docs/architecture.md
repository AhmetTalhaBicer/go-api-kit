# Architecture Overview

## Layers & Domain Structure

This project follows a **Domain-Driven (Package-by-Feature) Clean Architecture** approach.
Each domain lives inside `internal/<domain>/` and encapsulates its own model, repository interface, service, and HTTP handler.

```
cmd/api/main.go (Entry Point / Wire Dependencies)
   │
   └── internal/domain1/
         ├── model.go       (Domain Entities)
         ├── repository.go  (Data Access Interface)
         ├── service.go     (Business Logic)
         └── handler.go     (HTTP Handler & Request DTOs)
```

---

## Package Responsibilities

| Package | Responsibility |
|---------|---------------|
| `cmd/api` | Entry point. Wires dependencies, connects DB, and starts HTTP server. |
| `internal/config` | Loads typed configuration from environment variables. |
| `internal/db` | PostgreSQL `pgxpool` connection and embedded `goose` migration runner. |
| `internal/db/migrations` | SQL migration files (`.sql`) used by goose and sqlc. |
| `internal/<domain>` | Complete domain logic (model, repository interface, service, HTTP handler). |
| `internal/middleware` | HTTP middlewares: logger, panic recovery, CORS. |
| `db/queries` | SQL queries formatted for `sqlc`. |
| `db/sqlc` | Type-safe Go code generated automatically by `sqlc`. |
| `test/integration` | End-to-end and integration tests. |
| `test/mocks` | Directory for generated or hand-written mocks. |
| `pkg/logger` | Thin wrapper around standard `log/slog`. |
| `pkg/response` | Consistent JSON response helpers. |
| `pkg/validator` | Fluent input validation without external dependencies. |

---

## Adding a New Domain Feature

1. Create a new directory under `internal/<domain>/` (e.g. `internal/user/`).
2. Add `model.go` for your domain entity struct.
3. Define the `Repository` interface in `repository.go`.
4. Add SQL migrations in `internal/db/migrations/` and SQL queries in `db/queries/`.
5. Run `make sqlc` to generate type-safe database code in `db/sqlc/`.
6. Implement `service.go` depending on the `Repository` interface.
7. Implement `handler.go` with request DTOs and HTTP routes.
8. Wire the new domain service and handler in `cmd/api/main.go`.

---

## Technology Choices

| Concern | Choice | Rationale |
|---------|--------|-----------|
| HTTP | `net/http` (stdlib) | Zero third-party web framework lock-in; Go 1.22+ ServeMux routing. |
| Logging | `log/slog` (stdlib) | Structured JSON/text logging built into stdlib. |
| Database Driver | `pgx/v5` | High-performance PostgreSQL driver and connection pool. |
| Migrations | `goose` | SQL migrations embedded into binary via `//go:embed`. |
| Code Gen | `sqlc` | Compile-time type-safe Go generated directly from raw SQL queries. |
| Configuration | `os.Getenv` | Simple 12-factor configuration with fallback defaults. |
| Validation | Custom `pkg/validator` | Fluent API with zero external dependencies. |
