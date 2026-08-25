package worker

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/service"
	"testing"
	"time"
)

func TestPrivateWorkerHonorsCanceledContext(t *testing.T) {
	a := workerApp(t)
	now := time.Now().UTC()
	j := domain.OutboxJob{ID: service.ID(), Kind: "ok", AggregateID: "x", Payload: `{"x":1}`, Status: "running", Attempts: 1, AvailableAt: now, CreatedAt: now, UpdatedAt: now}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if e := New(a, time.Millisecond, nil).handle(ctx, j); e != context.Canceled {
		t.Fatalf("got %v", e)
	}
	var n int
	if e := a.DB.QueryRow(context.Background(), "SELECT COUNT(*) FROM audit_events").Scan(&n); e != nil {
		t.Fatal(e)
	}
	if n != 0 {
		t.Fatalf("audit rows=%d", n)
	}
}
