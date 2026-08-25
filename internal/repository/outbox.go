package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"time"
)

type Outbox struct{ DB *storage.DB }

func (r Outbox) Enqueue(ctx context.Context, j domain.OutboxJob) error {
	_, e := r.DB.Exec(ctx, `INSERT INTO outbox_jobs(id,kind,aggregate_id,payload,status,attempts,available_at,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?)`, j.ID, j.Kind, j.AggregateID, j.Payload, "pending", 0, j.AvailableAt.Format(timeLayout), j.CreatedAt.Format(timeLayout), j.UpdatedAt.Format(timeLayout))
	return e
}
func (r Outbox) Claim(ctx context.Context, now time.Time) (domain.OutboxJob, error) {
	var j domain.OutboxJob
	var av, cr, up string
	e := r.DB.QueryRow(ctx, `UPDATE outbox_jobs SET status='running',attempts=attempts+1,updated_at=? WHERE id=(SELECT id FROM outbox_jobs WHERE status='pending' AND available_at<=? ORDER BY available_at LIMIT 1) RETURNING id,kind,aggregate_id,payload,status,attempts,available_at,created_at,updated_at,last_error`, now.Format(timeLayout), now.Format(timeLayout)).Scan(&j.ID, &j.Kind, &j.AggregateID, &j.Payload, &j.Status, &j.Attempts, &av, &cr, &up, &j.LastError)
	if e == sql.ErrNoRows {
		return j, domain.ErrNotFound
	}
	if e != nil {
		return j, e
	}
	j.AvailableAt, _ = time.Parse(timeLayout, av)
	j.CreatedAt, _ = time.Parse(timeLayout, cr)
	j.UpdatedAt, _ = time.Parse(timeLayout, up)
	return j, nil
}
func (r Outbox) Finish(ctx context.Context, id string) error {
	_, e := r.DB.Exec(ctx, `UPDATE outbox_jobs SET status='done',updated_at=? WHERE id=?`, time.Now().UTC().Format(timeLayout), id)
	return e
}
func (r Outbox) Fail(ctx context.Context, j domain.OutboxJob, msg string, next time.Time, permanent bool) error {
	status := "pending"
	if permanent {
		status = "failed"
	}
	_, e := r.DB.Exec(ctx, `UPDATE outbox_jobs SET status=?,available_at=?,updated_at=?,last_error=? WHERE id=?`, status, next.Format(timeLayout), time.Now().UTC().Format(timeLayout), msg, j.ID)
	return e
}
