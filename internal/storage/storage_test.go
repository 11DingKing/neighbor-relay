package storage

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
)

func openTest(t *testing.T) *DB {
	t.Helper()
	d, e := Open(context.Background(), filepath.Join(t.TempDir(), "test.db"))
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { d.Close() })
	return d
}
func TestMigrationCreatesTables(t *testing.T) {
	d := openTest(t)
	n, e := d.TableCount(context.Background())
	if e != nil {
		t.Fatal(e)
	}
	if n < 10 {
		t.Fatalf("tables=%d", n)
	}
	if e = d.Health(context.Background()); e != nil {
		t.Fatal(e)
	}
}
func TestMigrationRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "restart.db")
	d, e := Open(context.Background(), path)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = d.Exec(context.Background(), `INSERT INTO users(id,name,email,role,active,created_at) VALUES('u','U','u@x','admin',1,'2026-01-01T00:00:00Z')`); e != nil {
		t.Fatal(e)
	}
	d.Close()
	d, e = Open(context.Background(), path)
	if e != nil {
		t.Fatal(e)
	}
	defer d.Close()
	var n int
	if e = d.QueryRow(context.Background(), `SELECT COUNT(*) FROM users`).Scan(&n); e != nil || n != 1 {
		t.Fatalf("restart n=%d err=%v", n, e)
	}
}
func TestForeignKeys(t *testing.T) {
	d := openTest(t)
	_, e := d.Exec(context.Background(), `INSERT INTO cases(id,household_id,title,summary,status,owner_id,version,created_at,updated_at) VALUES('c','missing','t','s','open','missing',1,'x','x')`)
	if e == nil {
		t.Fatal("foreign key disabled")
	}
}
func TestTransactionRollback(t *testing.T) {
	d := openTest(t)
	e := d.WithTx(context.Background(), func(tx *sql.Tx) error {
		_, e := tx.ExecContext(context.Background(), `INSERT INTO users(id,name,email,role,active,created_at) VALUES('u','U','u@x','admin',1,'x')`)
		if e != nil {
			return e
		}
		return errors.New("abort")
	})
	if e == nil {
		t.Fatal("rollback error missing")
	}
	var n int
	if e = d.QueryRow(context.Background(), `SELECT COUNT(*) FROM users`).Scan(&n); e != nil || n != 0 {
		t.Fatalf("rollback failed n=%d err=%v", n, e)
	}
}
