package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Migrate applies all pending database migrations using goose from the root ./migrations directory.
func Migrate(ctx context.Context, dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("database: migrate: open: %w", err)
	}
	defer db.Close()

	// Read migration SQL files directly from root ./migrations folder
	goose.SetBaseFS(os.DirFS("."))

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("database: migrate: set dialect: %w", err)
	}

	if err := goose.RunContext(ctx, "up", db, "migrations"); err != nil {
		return fmt.Errorf("database: migrate: run: %w", err)
	}

	return nil
}
