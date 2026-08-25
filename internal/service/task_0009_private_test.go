package service

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/config"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"log/slog"
	"path/filepath"
	"testing"
	"time"
)

func TestPrivatePlanVisitRollsBackWhenAuditFails(t *testing.T) {
	ctx := context.Background()
	db, e := storage.Open(ctx, filepath.Join(t.TempDir(), "a.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	a := New(db, config.Load(), slog.Default())
	u := domain.User{ID: "a9", Name: "A", Email: "a9@e", Role: domain.RoleAdmin, Active: true}
	w := domain.User{ID: "w9", Name: "W", Email: "w9@e", Role: domain.RoleFieldWorker, Active: true}
	if e = a.EnsureUser(ctx, u, "p"); e != nil {
		t.Fatal(e)
	}
	h, e := a.CreateHousehold(ctx, u, "1", "R", "1", 1)
	if e != nil {
		t.Fatal(e)
	}
	c, e := a.CreateCase(ctx, u, h.ID, "T", "S")
	if e != nil {
		t.Fatal(e)
	}
	if _, e = db.Exec(ctx, `CREATE TRIGGER reject_audit BEFORE INSERT ON audit_events BEGIN SELECT RAISE(ABORT,'audit down'); END`); e != nil {
		t.Fatal(e)
	}
	_, e = a.PlanVisit(ctx, u, c.ID, w.ID, time.Now().UTC().Add(time.Hour), "k9", "r9")
	if e == nil {
		t.Fatal("expected audit failure")
	}
	var n int
	if e = db.QueryRow(ctx, "SELECT COUNT(*) FROM visits").Scan(&n); e != nil {
		t.Fatal(e)
	}
	if n != 0 {
		t.Fatalf("visits=%d", n)
	}
}
