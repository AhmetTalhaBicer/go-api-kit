#!/usr/bin/env bash
# scripts/migrate.sh — Run database migrations using the goose CLI.
#
# Requires:
#   goose CLI — install with:
#     go install github.com/pressly/goose/v3/cmd/goose@latest
#
#   DATABASE_URL environment variable set in your .env file, e.g.:
#     DATABASE_URL=postgres://postgres:secret@localhost:5432/mydb?sslmode=disable
#
# Usage:
#   ./scripts/migrate.sh up                # apply all pending migrations
#   ./scripts/migrate.sh down              # roll back the last migration
#   ./scripts/migrate.sh status            # show migration status
#   ./scripts/migrate.sh create <name>     # create a new migration file

set -euo pipefail

DIRECTION=${1:-"up"}
MIGRATIONS_DIR="./internal/database/migrations"

if ! command -v goose &>/dev/null; then
    echo "Error: goose not found."
    echo "Install: go install github.com/pressly/goose/v3/cmd/goose@latest"
    exit 1
fi

if [[ -z "${DATABASE_URL:-}" ]]; then
    echo "Error: DATABASE_URL is not set."
    echo "Copy .env.example to .env and configure DATABASE_URL."
    exit 1
fi

if [[ "$DIRECTION" == "create" ]]; then
    NAME=${2:-"new_migration"}
    goose -dir "$MIGRATIONS_DIR" create "$NAME" sql
    echo "Migration file created in $MIGRATIONS_DIR"
else
    goose -dir "$MIGRATIONS_DIR" postgres "$DATABASE_URL" "$DIRECTION"
fi
