package db

import (
	"context"
	"database/sql"
	"fmt"
)

// DB wraps an sql.DB connection with convenience queries.
type DB struct {
	conn *sql.DB
}

// New wraps an existing *sql.DB.
func New(conn *sql.DB) *DB {
	return &DB{conn: conn}
}

// Close closes the underlying database connection.
func (d *DB) Close() error {
	return d.conn.Close()
}

// Conn returns the raw sql.DB connection.
func (d *DB) Conn() *sql.DB {
	return d.conn
}

// CountAvailable returns the total number of packages in available table.
func (d *DB) CountAvailable(ctx context.Context) (int, error) {
	var count int
	err := d.conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM available").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count available packages: %w", err)
	}
	return count, nil
}

// CountInstalled returns the total number of installed packages.
func (d *DB) CountInstalled(ctx context.Context) (int, error) {
	var count int
	err := d.conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM installed").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count installed packages: %w", err)
	}
	return count, nil
}
