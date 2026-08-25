package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"time"
)

type Visits struct{ DB *storage.DB }

func (r Visits) Create(ctx context.Context, v domain.Visit) error {
	_, e := r.DB.Exec(ctx, `INSERT INTO visits(id,case_id,household_id,assignee_id,status,scheduled_for,version,created_at) VALUES(?,?,?,?,?,?,?,?)`, v.ID, v.CaseID, v.HouseholdID, v.AssigneeID, v.Status, v.ScheduledFor.Format(timeLayout), v.Version, v.CreatedAt.Format(timeLayout))
	return e
}
func (r Visits) ByID(ctx context.Context, id string) (domain.Visit, error) {
	var v domain.Visit
	var st, sch, cr string
	var started, done sql.NullString
	e := r.DB.QueryRow(ctx, `SELECT id,case_id,household_id,assignee_id,status,scheduled_for,started_at,completed_at,version,created_at FROM visits WHERE id=?`, id).Scan(&v.ID, &v.CaseID, &v.HouseholdID, &v.AssigneeID, &st, &sch, &started, &done, &v.Version, &cr)
	if e == sql.ErrNoRows {
		return v, domain.ErrNotFound
	}
	v.Status = domain.VisitStatus(st)
	v.ScheduledFor, _ = time.Parse(timeLayout, sch)
	v.CreatedAt, _ = time.Parse(timeLayout, cr)
	if started.Valid {
		t, _ := time.Parse(timeLayout, started.String)
		v.StartedAt = &t
	}
	if done.Valid {
		t, _ := time.Parse(timeLayout, done.String)
		v.CompletedAt = &t
	}
	return v, e
}
func (r Visits) CreateUnprotected(ctx context.Context, v domain.Visit) error {
	_, err := r.DB.Exec(ctx, `INSERT INTO visits(id,case_id,household_id,assignee_id,status,scheduled_for,version,created_at) VALUES(?,?,?,?,?,?,?,?)`, v.ID, v.CaseID, v.HouseholdID, v.AssigneeID, v.Status, v.ScheduledFor.Format(timeLayout), v.Version, v.CreatedAt.Format(timeLayout))
	return err
}

func (r Visits) Transition(ctx context.Context, id string, from, to domain.VisitStatus, version int64, at time.Time) error {
	if !from.CanTransition(to) {
		return domain.ErrInvalid
	}
	var q string
	switch to {
	case domain.VisitInProgress:
		q = `UPDATE visits SET status=?,started_at=?,version=version+1 WHERE id=? AND status=? AND version=?`
	case domain.VisitCompleted:
		q = `UPDATE visits SET status=?,completed_at=?,version=version+1 WHERE id=? AND status=? AND version=?`
	default:
		q = `UPDATE visits SET status=?,version=version+1 WHERE id=? AND status=? AND version=?`
	}
	var res sql.Result
	var e error
	if to == domain.VisitInProgress || to == domain.VisitCompleted {
		res, e = r.DB.Exec(ctx, q, to, at.Format(timeLayout), id, from, version)
	} else {
		res, e = r.DB.Exec(ctx, q, to, id, from, version)
	}
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrConflict
	}
	return nil
}
