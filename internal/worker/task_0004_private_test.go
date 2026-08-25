package worker

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/service"
	"log/slog"
	"testing"
	"time"
)

func TestPrivateTransientOutboxFailureRemainsPending(t *testing.T) {
	a := workerApp(t)
	now := time.Now().UTC()
	j := domain.OutboxJob{ID: service.ID(), Kind: "bad", AggregateID: "x", Payload: "{", Status: "pending", AvailableAt: now, CreatedAt: now, UpdatedAt: now}
	if e := a.Outbox.Enqueue(context.Background(), j); e != nil {
		t.Fatal(e)
	}
	w := New(a, time.Millisecond, slog.Default())
	w.tick(context.Background())
	var status string
	var attempts int
	if e := a.DB.QueryRow(context.Background(), "SELECT status,attempts FROM outbox_jobs WHERE id=?", j.ID).Scan(&status, &attempts); e != nil {
		t.Fatal(e)
	}
	if status != "pending" {
		t.Fatalf("status=%s attempts=%d", status, attempts)
	}
	if attempts != 1 {
		t.Fatalf("attempts=%d", attempts)
	}
}
