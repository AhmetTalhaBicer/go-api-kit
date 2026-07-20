# go-api-kit

A minimal, production-ready Go project template.

**Stack:** `net/http` · `log/slog` · `pgx/v5` · `goose` · `sqlc`

## Getting Started

```bash
# 1. Use this template on GitHub, then clone your new repo

# 2. Update the module path
go mod edit -module github.com/YOUR_USERNAME/YOUR_PROJECT
find . -type f -name "*.go" | xargs sed -i \
  's|github.com/username/go-api-kit|github.com/YOUR_USERNAME/YOUR_PROJECT|g'

# 3. Set up environment
cp .env.example .env

# 4. Install tools
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/pressly/goose/v3/cmd/goose@latest

# 5. Run
make migrate-up
make sqlc
make run
```

## Structure

```
.
├── cmd/api/                # API server entry point (main.go)
├── internal/
│   ├── config/             # Environment-based configuration
│   ├── database/           # pgxpool connection & goose embedded migrations
│   │   └── migrations/     # SQL migration files (.sql)
│   ├── domain1/            # Domain 1: model, repository, service, handler
│   ├── domain2/            # Domain 2: model, repository, service, handler
│   ├── health/             # Health check endpoints (/health, /healthz)
│   ├── logger/             # log/slog wrapper (JSON/text)
│   ├── middleware/         # RequestID, Timeout, Logger, Recover, CORS
│   ├── response/           # Standard JSON response helpers
│   └── validator/          # Input validation without dependencies
├── sqlc/
│   ├── queries/            # SQL query files (input for sqlc)
│   └── gen/                # Generated Go code from sqlc
├── test/
│   └── integration/        # Integration/API tests
├── api/                    # OpenAPI 3.0 spec
├── deployments/            # Dockerfile & docker-compose
├── docs/                   # Architecture documentation
├── scripts/                # Build & migration scripts
└── .github/workflows/      # CI/CD & Release pipelines
```

## Commands

| Command | Description |
|---------|-------------|
| `make run` | Start the server |
| `make watch` | Start live-reloading dev server (Air) |
| `make build` | Compile binary |
| `make test` | Run all tests |
| `make test-integration` | Run integration tests only |
| `make sqlc` | Generate code from SQL |
| `make migrate-up` | Apply migrations |
| `make migrate-down` | Roll back last migration |
| `make migrate-create NAME=x` | New migration file |
| `make lint` | Run golangci-lint |
| `make docker-up` | Start with Docker Compose |

## Adapting to Your Project

1. Rename `internal/domain1/` and `internal/domain2/` to your business concepts.
2. Update entity fields in each `model.go`.
3. Implement the `Repository` interface (wire sqlc `*Queries` inside).
4. Uncomment the DB connect and migrate blocks in `cmd/api/main.go`.
5. Update `internal/database/migrations/` and `sqlc/queries/` for your schema.

## License

MIT
