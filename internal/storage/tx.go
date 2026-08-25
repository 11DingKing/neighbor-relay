package storage

import (
	"context"
	"database/sql"
)

func (d *DB) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := d.SQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
func (d *DB) Exec(ctx context.Context, q string, args ...any) (sql.Result, error) {
	return d.SQL.ExecContext(ctx, q, args...)
}
func (d *DB) Query(ctx context.Context, q string, args ...any) (*sql.Rows, error) {
	return d.SQL.QueryContext(ctx, q, args...)
}
func (d *DB) QueryRow(ctx context.Context, q string, args ...any) *sql.Row {
	return d.SQL.QueryRowContext(ctx, q, args...)
}
