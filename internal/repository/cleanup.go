package repository

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"time"
)

type Cleanup struct{ DB *storage.DB }

func (r Cleanup) ExpireSessions(ctx context.Context, now time.Time) (int64, error) {
	res, e := r.DB.Exec(ctx, `DELETE FROM sessions WHERE expires_at<=?`, now.Format(timeLayout))
	if e != nil {
		return 0, e
	}
	return res.RowsAffected()
}
func (r Cleanup) DeleteOldIdempotency(ctx context.Context, before time.Time) (int64, error) {
	res, e := r.DB.Exec(ctx, `DELETE FROM idempotency_keys WHERE expires_at<=?`, before.Format(timeLayout))
	if e != nil {
		return 0, e
	}
	return res.RowsAffected()
}
func (r Cleanup) RecoverJobs(ctx context.Context, before time.Time) (int64, error) {
	res, e := r.DB.Exec(ctx, `UPDATE outbox_jobs SET status='pending',available_at=?,updated_at=? WHERE status='running' AND updated_at<?`, before.Format(timeLayout), before.Format(timeLayout), before.Format(timeLayout))
	if e != nil {
		return 0, e
	}
	return res.RowsAffected()
}
