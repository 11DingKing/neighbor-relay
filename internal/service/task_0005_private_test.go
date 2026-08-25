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

func TestPrivatePlanVisitReusesIdempotencyResult(t *testing.T) {
	ctx := context.Background()
	db, e := storage.Open(ctx, filepath.Join(t.TempDir(), "i.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	a := New(db, config.Load(), slog.Default())
	u := domain.User{ID: "a5", Name: "A", Email: "a5@e", Role: domain.RoleAdmin, Active: true}
	w := domain.User{ID: "w5", Name: "W", Email: "w5@e", Role: domain.RoleFieldWorker, Active: true}
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
	when := time.Now().UTC().Add(time.Hour)
	v1, e := a.PlanVisit(ctx, u, c.ID, w.ID, when, "same-key", "r1")
	if e != nil {
		t.Fatal(e)
	}
	v2, e := a.PlanVisit(ctx, u, c.ID, w.ID, when, "same-key", "r2")
	if e != nil {
		t.Fatal(e)
	}
	if v1.ID != v2.ID {
		t.Fatalf("idempotency returned %s then %s", v1.ID, v2.ID)
	}
}
