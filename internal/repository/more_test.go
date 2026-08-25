package repository

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"testing"
	"time"
)

func TestOutboxLifecycle(t *testing.T) {
	d := repoDB(t)
	r := Outbox{DB: d}
	now := time.Now().UTC()
	j := domain.OutboxJob{ID: "j", Kind: "notice", AggregateID: "a", Payload: "{}", AvailableAt: now, CreatedAt: now, UpdatedAt: now}
	if e := r.Enqueue(context.Background(), j); e != nil {
		t.Fatal(e)
	}
	got, e := r.Claim(context.Background(), now.Add(time.Second))
	if e != nil || got.ID != "j" || got.Attempts != 1 {
		t.Fatalf("job=%#v err=%v", got, e)
	}
	if e = r.Finish(context.Background(), got.ID); e != nil {
		t.Fatal(e)
	}
}
func TestIdempotencyLifecycle(t *testing.T) {
	d := repoDB(t)
	r := Idempotency{DB: d}
	now := time.Now().UTC()
	if e := r.Save(context.Background(), "k", "u", "op", "response", now, now.Add(time.Hour)); e != nil {
		t.Fatal(e)
	}
	v, e := r.Find(context.Background(), "k", "u", "op", now)
	if e != nil || v != "response" {
		t.Fatal(v, e)
	}
	if _, e = r.Find(context.Background(), "k", "other", "op", now); e == nil {
		t.Fatal("cross user hit")
	}
}
func TestAuditQuery(t *testing.T) {
	d := repoDB(t)
	r := Audit{DB: d}
	now := time.Now().UTC()
	for i := 0; i < 2; i++ {
		if e := r.Append(context.Background(), domain.AuditEvent{ID: string(rune('a' + i)), ActorID: "u", EntityType: "visit", EntityID: "v", Action: "a", Result: "success", RequestID: "r", CreatedAt: now}); e != nil {
			t.Fatal(e)
		}
	}
	q := AuditQuery{DB: d}
	items, e := q.ForEntity(context.Background(), "visit", "v", 10)
	if e != nil || len(items) != 2 {
		t.Fatalf("items=%d err=%v", len(items), e)
	}
	n, e := q.Count(context.Background(), "visit")
	if e != nil || n != 2 {
		t.Fatal(n, e)
	}
}
func TestCleanup(t *testing.T) {
	d := repoDB(t)
	r := Cleanup{DB: d}
	n, e := r.ExpireSessions(context.Background(), time.Now().UTC())
	if e != nil || n < 0 {
		t.Fatal(n, e)
	}
	n, e = r.DeleteOldIdempotency(context.Background(), time.Now().UTC())
	if e != nil || n < 0 {
		t.Fatal(n, e)
	}
}
