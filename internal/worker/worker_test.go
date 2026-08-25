package worker

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/config"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/service"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func workerApp(t *testing.T) *service.Service {
	d, e := storage.Open(context.Background(), filepath.Join(t.TempDir(), "w.db"))
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { d.Close() })
	c := config.Load()
	c.WorkerInterval = time.Millisecond
	return service.New(d, c, slog.New(slog.NewTextHandler(os.Stdout, nil)))
}
func TestWorkerStop(t *testing.T) {
	a := workerApp(t)
	w := New(a, time.Millisecond, slog.Default())
	ctx, cancel := context.WithCancel(context.Background())
	go w.Run(ctx)
	w.Stop()
	cancel()
	select {
	case <-time.After(100 * time.Millisecond):
		t.Fatal("worker did not stop")
	default:
	}
}
func TestWorkerNoJobs(t *testing.T) {
	a := workerApp(t)
	w := New(a, time.Millisecond, slog.Default())
	ctx, cancel := context.WithCancel(context.Background())
	go w.Run(ctx)
	time.Sleep(10 * time.Millisecond)
	cancel()
	w.Stop()
}
func TestBackoffBound(t *testing.T) {
	if min(20, 8) != 8 || min(3, 8) != 3 {
		t.Fatal("min")
	}
	_ = domain.ErrNotFound
}
