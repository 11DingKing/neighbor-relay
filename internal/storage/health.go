package storage

import (
	"context"
	"database/sql"
	"time"
)

func (d *DB) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	var n int
	return d.SQL.QueryRowContext(ctx, `SELECT 1`).Scan(&n)
}
func (d *DB) TableCount(ctx context.Context) (int, error) {
	var n int
	e := d.SQL.QueryRowContext(ctx, `SELECT COUNT(*) FROM sqlite_master WHERE type='table'`).Scan(&n)
	return n, e
}
func closeRows(rows *sql.Rows) {
	if rows != nil {
		_ = rows.Close()
	}
}
