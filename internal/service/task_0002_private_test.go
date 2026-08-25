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

func TestPrivateRevokedSessionCannotAuthenticate(t *testing.T) {
	ctx := context.Background()
	db, err := storage.Open(ctx, filepath.Join(t.TempDir(), "s.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	app := New(db, config.Load(), slog.Default())
	u := domain.User{ID: "u2", Name: "Worker", Email: "u2@example.com", Role: domain.RoleFieldWorker, Active: true}
	if err = app.EnsureUser(ctx, u, "pw"); err != nil {
		t.Fatal(err)
	}
	_, tok, err := app.Login(ctx, u.Email, "pw")
	if err != nil {
		t.Fatal(err)
	}
	if err = app.Logout(ctx, tok); err != nil {
		t.Fatal(err)
	}
	if _, err = app.Authenticate(ctx, tok); err == nil {
		t.Fatal("revoked token remained usable")
	}
}
