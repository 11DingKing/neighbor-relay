package worker

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/service"
	"time"
)

func Recover(ctx context.Context, app *service.Service) error {
	_, e := app.RecoverWorkers(ctx, 2*time.Minute)
	return e
}
func Backoff(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	if attempt > 8 {
		attempt = 8
	}
	return time.Duration(1<<attempt) * time.Second
}
func ShouldRetry(attempt, max int) bool { return attempt < max }
