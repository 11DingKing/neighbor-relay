package worker

import (
	"context"
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
func (w *Worker) handle(ctx context.Context, j domain.OutboxJob) error { return nil }

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
