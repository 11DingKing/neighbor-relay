package service

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/repository"
	"time"
)

func (s *Service) Cleanup(ctx context.Context) (int64, error) {
	r := repository.Cleanup{DB: s.DB}
	n, e := r.ExpireSessions(ctx, time.Now().UTC())
	if e != nil {
		return n, e
	}
	m, e := r.DeleteOldIdempotency(ctx, time.Now().UTC())
	return n + m, e
}
func (s *Service) RecoverWorkers(ctx context.Context, stale time.Duration) (int64, error) {
	r := repository.Cleanup{DB: s.DB}
	now := time.Now().UTC()
	return r.RecoverJobs(ctx, now.Add(-stale))
}
