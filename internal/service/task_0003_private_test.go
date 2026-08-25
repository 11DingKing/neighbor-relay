package service

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/config"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"log/slog"
	"path/filepath"
	"testing"
)

func TestPrivateFieldWorkerCannotCreateCase(t *testing.T) {
	ctx := context.Background()
	db, e := storage.Open(ctx, filepath.Join(t.TempDir(), "c.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	app := New(db, config.Load(), slog.Default())
	admin := domain.User{ID: "a3", Name: "A", Email: "a3@e", Role: domain.RoleAdmin, Active: true}
	w := domain.User{ID: "w3", Name: "W", Email: "w3@e", Role: domain.RoleFieldWorker, Active: true}
	if e = app.EnsureUser(ctx, admin, "p"); e != nil {
		t.Fatal(e)
	}
	h, e := app.CreateHousehold(ctx, admin, "1 Lane", "R", "1", 1)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = app.CreateCase(ctx, w, h.ID, "x", "y"); e != domain.ErrForbidden {
		t.Fatalf("got %v", e)
	}
}
