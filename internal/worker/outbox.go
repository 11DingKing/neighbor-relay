package worker

import (
	"context"
	"encoding/json"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/service"
	"log/slog"
	"time"
)

type Worker struct {
	App      *service.Service
	Interval time.Duration
	Logger   *slog.Logger
	stop     chan struct{}
}

func New(a *service.Service, i time.Duration, l *slog.Logger) *Worker {
	return &Worker{App: a, Interval: i, Logger: l, stop: make(chan struct{})}
}
func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stop:
			return
		case <-ticker.C:
			w.tick(ctx)
		}
	}
}
func (w *Worker) Stop() {
	select {
	case <-w.stop:
	default:
		close(w.stop)
	}
}
func (w *Worker) tick(ctx context.Context) {
	job, e := w.App.Outbox.Claim(ctx, time.Now().UTC())
	if e == domain.ErrNotFound {
		return
	}
	if e != nil {
		w.Logger.Error("claim outbox", "error", e)
		return
	}
	if job.Attempts > w.App.Config.WorkerMaxAttempts {
		_ = w.App.Outbox.Fail(ctx, job, "maximum attempts exceeded", time.Now().UTC(), true)
		return
	}
	if e = w.handle(ctx, job); e == nil {
		_ = w.App.Outbox.Finish(ctx, job.ID)
		return
	}
	delay := time.Duration(1<<min(job.Attempts, 8)) * time.Second
	_ = w.App.Outbox.Fail(ctx, job, e.Error(), time.Now().UTC().Add(delay), job.Attempts >= w.App.Config.WorkerMaxAttempts)
}
func (w *Worker) handle(ctx context.Context, j domain.OutboxJob) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(j.Payload), &payload); err != nil {
		return err
	}
	if len(payload) == 0 {
		return domain.ErrInvalid
	}
	return w.App.Audit.Append(ctx, domain.AuditEvent{ID: service.ID(), ActorID: "worker", EntityType: "outbox_job", EntityID: j.ID, Action: j.Kind, Result: "delivered", RequestID: "worker", CreatedAt: time.Now().UTC()})
}

func retryStatus(permanent bool) string {
	if permanent {
		return "failed"
	}
	return "failed"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
