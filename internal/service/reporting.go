package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"time"
)

type VisitSummary struct{ Planned, InProgress, Completed, Cancelled int }

func (s *Service) Summary(ctx context.Context, u domain.User) (VisitSummary, error) {
	if e := s.RequireRole(u, domain.RoleAdmin, domain.RoleFieldWorker); e != nil {
		return VisitSummary{}, e
	}
	rows, e := s.DB.Query(ctx, `SELECT status,COUNT(*) FROM visits GROUP BY status`)
	if e != nil {
		return VisitSummary{}, e
	}
	defer rows.Close()
	var out VisitSummary
	for rows.Next() {
		var st string
		var n int
		if e := rows.Scan(&st, &n); e != nil {
			return out, e
		}
		switch st {
		case "planned":
			out.Planned = n
		case "in_progress":
			out.InProgress = n
		case "completed":
			out.Completed = n
		case "cancelled":
			out.Cancelled = n
		}
	}
	return out, rows.Err()
}
func (s *Service) AddNote(ctx context.Context, u domain.User, visitID, body string) error {
	if e := domain.ValidateNote(body); e != nil {
		return e
	}
	v, e := s.Visits.ByID(ctx, visitID)
	if e != nil {
		return e
	}
	if u.Role == domain.RoleFieldWorker && v.AssigneeID != u.ID {
		return domain.ErrForbidden
	}
	_, e = s.DB.Exec(ctx, `INSERT INTO visit_notes(id,visit_id,author_id,body,created_at) VALUES(?,?,?,?,?)`, ID(), visitID, u.ID, body, time.Now().UTC().Format(timeLayout))
	return e
}
func scanCount(row *sql.Row) (int, error) { var n int; e := row.Scan(&n); return n, e }
