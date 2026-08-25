package repository

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"path/filepath"
	"testing"
	"time"
)

func repoDB(t *testing.T) *storage.DB {
	d, e := storage.Open(context.Background(), filepath.Join(t.TempDir(), "r.db"))
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { d.Close() })
	return d
}
func seedUser(t *testing.T, d *storage.DB, id, role string) {
	_, e := d.Exec(context.Background(), `INSERT INTO users(id,name,email,role,active,created_at) VALUES(?,?,?,?,?,?)`, id, id, id+"@x", role, 1, time.Now().UTC().Format(timeLayout))
	if e != nil {
		t.Fatal(e)
	}
}
func TestUsersRoundTrip(t *testing.T) {
	d := repoDB(t)
	r := Users{DB: d}
	u := domain.User{ID: "u", Name: "Uma", Email: "u@x", Role: domain.RoleAdmin, Active: true, CreatedAt: time.Now().UTC()}
	if e := r.Create(context.Background(), u, "hash"); e != nil {
		t.Fatal(e)
	}
	got, h, e := r.ByEmail(context.Background(), u.Email)
	if e != nil || got.ID != u.ID || h != "hash" {
		t.Fatalf("%#v %s %v", got, h, e)
	}
}
func TestHouseholdList(t *testing.T) {
	d := repoDB(t)
	seedUser(t, d, "u", "admin")
	r := Households{DB: d}
	now := time.Now().UTC()
	for i := 0; i < 3; i++ {
		h := domain.Household{ID: string(rune('a' + i)), Address: "road", ContactName: "c", Phone: "p", CreatedAt: now, UpdatedAt: now}
		if e := r.Create(context.Background(), h); e != nil {
			t.Fatal(e)
		}
	}
	items, e := r.List(context.Background(), 2)
	if e != nil || len(items) != 2 {
		t.Fatalf("len=%d err=%v", len(items), e)
	}
}
func TestSessionRevocation(t *testing.T) {
	d := repoDB(t)
	seedUser(t, d, "u", "admin")
	r := Sessions{DB: d}
	now := time.Now().UTC()
	s := domain.Session{ID: "s", UserID: "u", TokenHash: "h", CreatedAt: &now, ExpiresAt: ptr(now.Add(time.Hour))}
	if e := r.Create(context.Background(), s); e != nil {
		t.Fatal(e)
	}
	if _, e := r.Active(context.Background(), "h", now); e != nil {
		t.Fatal(e)
	}
	if e := r.Revoke(context.Background(), "h", now); e != nil {
		t.Fatal(e)
	}
	if _, e := r.Active(context.Background(), "h", now); e == nil {
		t.Fatal("revoked session active")
	}
}
func ptr(t time.Time) *time.Time { return &t }
