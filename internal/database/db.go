package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect creates and verifies a new PostgreSQL connection pool.
//
// dsn accepts either a keyword/value string or a URL:
//
//	"host=localhost port=5432 user=postgres password=secret dbname=mydb sslmode=disable"
//	"postgres://postgres:secret@localhost:5432/mydb?sslmode=disable"
func Connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("database: connect: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database: ping: %w", err)
	}
	return pool, nil
}
