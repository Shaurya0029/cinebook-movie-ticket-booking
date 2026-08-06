package db

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.sql
var schemaSQL string

// Connect opens a pooled connection to Postgres using the given connection string.
func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("db: create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db: ping: %w", err)
	}
	return pool, nil
}

// EnsureSchema runs the embedded schema.sql, which is written entirely as
// idempotent CREATE TABLE/INDEX IF NOT EXISTS statements, so it is safe to
// call on every startup of both cmd/api and cmd/seed.
func EnsureSchema(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, schemaSQL); err != nil {
		return fmt.Errorf("db: ensure schema: %w", err)
	}
	return nil
}
