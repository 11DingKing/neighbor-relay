package service

import (
	"context"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	"github.com/11DingKing/neighbor-relay/internal/config"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/storage"
)

func TestPrivateCompleteVisitRollsBackWhenNoteStoreFails(t *testing.T) {
	ctx := context.Background()
	db, err := storage.Open(ctx, filepath.Join(t.TempDir(), "private.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	app := New(db, config.Load(), slog.Default())
	admin := domain.User{ID: "private-admin", Name: "Admin", Email: "private-admin@example.com", Role: domain.RoleAdmin, Active: true}
	worker := domain.User{ID: "private-worker", Name: "Worker", Email: "private-worker@example.com", Role: domain.RoleFieldWorker, Active: true}
	if err := app.EnsureUser(ctx, admin, "admin-password"); err != nil {
		t.Fatal(err)
	}
	if err := app.EnsureUser(ctx, worker, "worker-password"); err != nil {
		t.Fatal(err)
	}
	household, err := app.CreateHousehold(ctx, admin, "1 Lane", "Resident", "10086", 1)
	if err != nil {
		t.Fatal(err)
	}
	caseItem, err := app.CreateCase(ctx, admin, household.ID, "Follow-up", "Needs an in-person check")
	if err != nil {
		t.Fatal(err)
	}
	visit, err := app.PlanVisit(ctx, admin, caseItem.ID, worker.ID, time.Now().UTC().Add(time.Hour), "private-key", "private-request")
	if err != nil {
		t.Fatal(err)
	}
	if err := app.StartVisit(ctx, worker, visit.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(ctx, `CREATE TRIGGER reject_private_notes BEFORE INSERT ON visit_notes BEGIN SELECT RAISE(ABORT, 'notes unavailable'); END`); err != nil {
		t.Fatal(err)
	}
	err = app.CompleteVisit(ctx, worker, visit.ID, "The family requested another visit")
	if err == nil {
		t.Fatal("completion must report note persistence failure")
	}
	got, err := app.Visits.ByID(ctx, visit.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.VisitInProgress {
		t.Fatalf("failed completion changed visit status to %s", got.Status)
	}
	var notes int
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM visit_notes WHERE visit_id=?`, visit.ID).Scan(&notes); err != nil {
		t.Fatal(err)
	}
	if notes != 0 {
		t.Fatalf("failed completion persisted %d notes", notes)
	}
}
