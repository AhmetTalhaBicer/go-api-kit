package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for database/sql
	"github.com/pressly/goose/v3"
)

// embedMigrations embeds all SQL migration files at compile time.
// Goose will apply them in order when Migrate() is called.
//
//go:embed migrations/*.sql
var embedMigrations embed.FS

// Migrate applies all pending database migrations using goose.
// It reads SQL files from the embedded migrations/ directory.
//
// Call this once at application startup, before serving traffic.
func Migrate(ctx context.Context, dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("database: migrate: open: %w", err)
	}
	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("database: migrate: set dialect: %w", err)
	}

	if err := goose.RunContext(ctx, "up", db, "migrations"); err != nil {
		return fmt.Errorf("database: migrate: run: %w", err)
	}

	return nil
}
