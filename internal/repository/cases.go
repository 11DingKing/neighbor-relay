package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"time"
)

type Cases struct{ DB *storage.DB }

func (r Cases) Create(ctx context.Context, c domain.Case) error {
	_, e := r.DB.Exec(ctx, `INSERT INTO cases(id,household_id,title,summary,status,owner_id,version,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?)`, c.ID, c.HouseholdID, c.Title, c.Summary, c.Status, c.OwnerID, c.Version, c.CreatedAt.Format(timeLayout), c.UpdatedAt.Format(timeLayout))
	return e
}
func (r Cases) ByID(ctx context.Context, id string) (domain.Case, error) {
	var c domain.Case
	var st, cr, up string
	e := r.DB.QueryRow(ctx, `SELECT id,household_id,title,summary,status,owner_id,version,created_at,updated_at FROM cases WHERE id=?`, id).Scan(&c.ID, &c.HouseholdID, &c.Title, &c.Summary, &st, &c.OwnerID, &c.Version, &cr, &up)
	if e == sql.ErrNoRows {
		return c, domain.ErrNotFound
	}
	c.Status = domain.CaseStatus(st)
	c.CreatedAt, _ = time.Parse(timeLayout, cr)
	c.UpdatedAt, _ = time.Parse(timeLayout, up)
	return c, e
}
func (r Cases) Transition(ctx context.Context, id string, from, to domain.CaseStatus, version int64, at time.Time) error {
	if !from.CanTransition(to) {
		return domain.ErrInvalid
	}
	res, e := r.DB.Exec(ctx, `UPDATE cases SET status=?,version=version+1,updated_at=? WHERE id=? AND status=?`, to, at.Format(timeLayout), id, from)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrConflict
	}
	return nil
}
