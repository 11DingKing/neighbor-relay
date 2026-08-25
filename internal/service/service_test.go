package service

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/config"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func testService(t *testing.T) *Service {
	d, e := storage.Open(context.Background(), filepath.Join(t.TempDir(), "s.db"))
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { d.Close() })
	c := config.Load()
	c.SessionTTL = time.Hour
	s := New(d, c, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	if err := s.EnsureUser(context.Background(), admin(), "admin-password"); err != nil {
		t.Fatal(err)
	}
	if err := s.EnsureUser(context.Background(), worker(), "worker-password"); err != nil {
		t.Fatal(err)
	}
	return s
}
func admin() domain.User {
	return domain.User{ID: "admin", Name: "Admin", Email: "admin@x", Role: domain.RoleAdmin, Active: true}
}
func worker() domain.User {
	return domain.User{ID: "worker", Name: "Worker", Email: "worker@x", Role: domain.RoleFieldWorker, Active: true}
}
func TestEnsureAndLogin(t *testing.T) {
	s := testService(t)
	u := domain.User{ID: "login", Name: "Login", Email: "login@x", Role: domain.RoleAdmin, Active: true}
	if e := s.EnsureUser(context.Background(), u, "secret"); e != nil {
		t.Fatal(e)
	}
	got, token, e := s.Login(context.Background(), u.Email, "secret")
	if e != nil || got.ID != u.ID || token == "" {
		t.Fatalf("login %#v %s %v", got, token, e)
	}
	if _, e = s.Authenticate(context.Background(), token); e != nil {
		t.Fatal(e)
	}
	if e = s.Logout(context.Background(), token); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Authenticate(context.Background(), token); e == nil {
		t.Fatal("logout did not revoke")
	}
}
func TestWrongPassword(t *testing.T) {
	s := testService(t)
	u := domain.User{ID: "wrong", Name: "Wrong", Email: "wrong@x", Role: domain.RoleAdmin, Active: true}
	_ = s.EnsureUser(context.Background(), u, "secret")
	if _, _, e := s.Login(context.Background(), u.Email, "bad"); e != domain.ErrUnauthorized {
		t.Fatalf("err=%v", e)
	}
}
func TestHouseholdAndCase(t *testing.T) {
	s := testService(t)
	u := admin()
	h, e := s.CreateHousehold(context.Background(), u, "1 Main", "A", "123", 2)
	if e != nil {
		t.Fatal(e)
	}
	c, e := s.CreateCase(context.Background(), u, h.ID, "Follow up", "needs visit")
	if e != nil || c.HouseholdID != h.ID {
		t.Fatalf("case=%#v err=%v", c, e)
	}
	if e = s.TransitionCase(context.Background(), u, c.ID, domain.CaseInProgress); e != nil {
		t.Fatal(e)
	}
}
func TestRoleBoundaries(t *testing.T) {
	s := testService(t)
	if _, e := s.CreateCase(context.Background(), worker(), "x", "t", "s"); e != domain.ErrForbidden {
		t.Fatalf("worker created case: %v", e)
	}
	if _, e := s.CreateHousehold(context.Background(), domain.User{Role: "guest", Active: true}, "a", "b", "c", 0); e != domain.ErrForbidden {
		t.Fatalf("guest accepted: %v", e)
	}
}
func TestVisitLifecycle(t *testing.T) {
	s := testService(t)
	a := admin()
	w := worker()
	_ = s.EnsureUser(context.Background(), a, "a")
	_ = s.EnsureUser(context.Background(), w, "w")
	h, e := s.CreateHousehold(context.Background(), a, "road", "A", "p", 1)
	if e != nil {
		t.Fatal(e)
	}
	c, e := s.CreateCase(context.Background(), a, h.ID, "case", "summary")
	if e != nil {
		t.Fatal(e)
	}
	v, e := s.PlanVisit(context.Background(), a, c.ID, w.ID, time.Now().UTC().Add(time.Hour), "key", "req")
	if e != nil {
		t.Fatal(e)
	}
	if e = s.StartVisit(context.Background(), w, v.ID); e != nil {
		t.Fatal(e)
	}
	if e = s.CompleteVisit(context.Background(), w, v.ID, "visited and explained"); e != nil {
		t.Fatal(e)
	}
	got, e := s.Visits.ByID(context.Background(), v.ID)
	if e != nil || got.Status != domain.VisitCompleted {
		t.Fatalf("visit %#v err=%v", got, e)
	}
}
func TestIdempotentPlan(t *testing.T) {
	s := testService(t)
	a := admin()
	h, _ := s.CreateHousehold(context.Background(), a, "road", "A", "p", 1)
	c, _ := s.CreateCase(context.Background(), a, h.ID, "case", "summary")
	when := time.Now().UTC().Add(time.Hour)
	v1, e := s.PlanVisit(context.Background(), a, c.ID, "worker", when, "same", "r")
	if e != nil {
		t.Fatal(e)
	}
	v2, e := s.PlanVisit(context.Background(), a, c.ID, "worker", when, "same", "r")
	if e != nil || v1.ID != v2.ID {
		t.Fatalf("not idempotent %v %v", v1.ID, v2.ID)
	}
}
