package repository

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"time"
)

type VisitQuery struct{ DB *storage.DB }

func (r VisitQuery) ForAssignee(ctx context.Context, id, status string, limit int) ([]domain.Visit, error) {
	if limit < 1 || limit > 100 {
		limit = 50
	}
	q := `SELECT id,case_id,household_id,assignee_id,status,scheduled_for,version,created_at FROM visits WHERE assignee_id=?`
	args := []any{id}
	if status != "" {
		q += ` AND status=?`
		args = append(args, status)
	}
	q += ` ORDER BY scheduled_for ASC LIMIT ?`
	args = append(args, limit)
	rows, e := r.DB.Query(ctx, q, args...)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := make([]domain.Visit, 0)
	for rows.Next() {
		var v domain.Visit
		var st, sch, cr string
		if e := rows.Scan(&v.ID, &v.CaseID, &v.HouseholdID, &v.AssigneeID, &st, &sch, &v.Version, &cr); e != nil {
			return nil, e
		}
		v.Status = domain.VisitStatus(st)
		v.ScheduledFor, _ = time.Parse(timeLayout, sch)
		v.CreatedAt, _ = time.Parse(timeLayout, cr)
		out = append(out, v)
	}
	return out, rows.Err()
}
func (r VisitQuery) CountByStatus(ctx context.Context, status string) (int, error) {
	var n int
	e := r.DB.QueryRow(ctx, `SELECT COUNT(*) FROM visits WHERE status=?`, status).Scan(&n)
	return n, e
}
func (r VisitQuery) HasOpenForCase(ctx context.Context, caseID string) (bool, error) {
	var n int
	e := r.DB.QueryRow(ctx, `SELECT COUNT(*) FROM visits WHERE case_id=? AND status IN ('planned','in_progress')`, caseID).Scan(&n)
	return n > 0, e
}
