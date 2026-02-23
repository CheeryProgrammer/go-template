package store

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib" // register pgx driver for database/sql
)

// Postgres implements Store backed by a PostgreSQL database.
type Postgres struct {
	db *sql.DB
}

// NewPostgres opens a connection to the PostgreSQL database at dsn and verifies
// it is reachable. The caller must call Close when done.
func NewPostgres(dsn string) (*Postgres, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	if err := db.PingContext(context.Background()); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}
	return &Postgres{db: db}, nil
}

// Ping checks whether the database connection is still alive.
func (p *Postgres) Ping(ctx context.Context) error {
	if err := p.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}
	return nil
}

// Close releases the database connection pool.
func (p *Postgres) Close() {
	if err := p.db.Close(); err != nil {
		slog.Error("closing database connection", "err", err)
	}
}
