package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	pkg "github.com/mihakotik/pm/internal/package"
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

// SaveAvailableBatch writes a list of available packages to the database in a single transaction.
func (d *DB) SaveAvailableBatch(ctx context.Context, pkgs []pkg.Package) error {
	if len(pkgs) == 0 {
		return nil
	}

	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT OR REPLACE INTO available (name, version, provider, description, depends, url, checksum)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	for _, p := range pkgs {
		depBytes, err := json.Marshal(p.Depends)
		depStr := ""
		if err == nil {
			depStr = string(depBytes)
		}

		_, err = stmt.ExecContext(ctx,
			p.Name,
			p.Version,
			p.Provider,
			p.Description,
			depStr,
			p.URL,
			p.Checksum,
		)
		if err != nil {
			return fmt.Errorf("failed to insert available package %s (%s): %w", p.Name, p.Provider, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit available packages transaction: %w", err)
	}

	return nil
}

// SearchAvailable queries packages matching name or description.
func (d *DB) SearchAvailable(ctx context.Context, query string) ([]pkg.Package, error) {
	like := "%" + query + "%"
	rows, err := d.conn.QueryContext(ctx, `
		SELECT name, version, provider, description, depends, url, checksum
		FROM available
		WHERE name LIKE ? OR description LIKE ?
		ORDER BY name ASC
	`, like, like)
	if err != nil {
		return nil, fmt.Errorf("failed to search available packages: %w", err)
	}
	defer rows.Close()

	var result []pkg.Package
	for rows.Next() {
		var p pkg.Package
		var depStr string
		if err := rows.Scan(&p.Name, &p.Version, &p.Provider, &p.Description, &depStr, &p.URL, &p.Checksum); err != nil {
			return nil, fmt.Errorf("failed to scan package row: %w", err)
		}
		if depStr != "" {
			_ = json.Unmarshal([]byte(depStr), &p.Depends)
		}
		result = append(result, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading available packages: %w", err)
	}

	return result, nil
}

// GetAvailable retrieves an available package by name and optional provider.
func (d *DB) GetAvailable(ctx context.Context, name, provider string) (*pkg.Package, error) {
	var query string
	var args []any
	if provider != "" {
		query = `SELECT name, version, provider, description, depends, url, checksum FROM available WHERE name = ? AND provider = ? LIMIT 1`
		args = []any{name, provider}
	} else {
		query = `SELECT name, version, provider, description, depends, url, checksum FROM available WHERE name = ? LIMIT 1`
		args = []any{name}
	}

	var p pkg.Package
	var depStr string
	err := d.conn.QueryRowContext(ctx, query, args...).Scan(&p.Name, &p.Version, &p.Provider, &p.Description, &depStr, &p.URL, &p.Checksum)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("package %s not found: %w", name, err)
		}
		return nil, fmt.Errorf("failed to query package %s: %w", name, err)
	}

	if depStr != "" {
		_ = json.Unmarshal([]byte(depStr), &p.Depends)
	}

	return &p, nil
}
