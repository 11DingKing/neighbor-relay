package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"time"
)

type Households struct{ DB *storage.DB }

func (r Households) Create(ctx context.Context, h domain.Household) error {
	_, e := r.DB.Exec(ctx, `INSERT INTO households(id,address,contact_name,phone,priority,created_at,updated_at) VALUES(?,?,?,?,?,?,?)`, h.ID, h.Address, h.ContactName, h.Phone, h.Priority, h.CreatedAt.Format(timeLayout), h.UpdatedAt.Format(timeLayout))
	return e
}
func (r Households) ByID(ctx context.Context, id string) (domain.Household, error) {
	var h domain.Household
	var c, u string
	e := r.DB.QueryRow(ctx, `SELECT id,address,contact_name,phone,priority,created_at,updated_at FROM households WHERE id=?`, id).Scan(&h.ID, &h.Address, &h.ContactName, &h.Phone, &h.Priority, &c, &u)
	if e == sql.ErrNoRows {
		return h, domain.ErrNotFound
	}
	h.CreatedAt, _ = time.Parse(timeLayout, c)
	h.UpdatedAt, _ = time.Parse(timeLayout, u)
	return h, e
}
func (r Households) List(ctx context.Context, limit int) ([]domain.Household, error) {
	if limit < 1 || limit > 100 {
		limit = 50
	}
	rows, e := r.DB.Query(ctx, `SELECT id,address,contact_name,phone,priority,created_at,updated_at FROM households ORDER BY priority DESC,created_at DESC LIMIT ?`, limit)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := make([]domain.Household, 0)
	for rows.Next() {
		var h domain.Household
		var c, u string
		if e := rows.Scan(&h.ID, &h.Address, &h.ContactName, &h.Phone, &h.Priority, &c, &u); e != nil {
			return nil, e
		}
		h.CreatedAt, _ = time.Parse(timeLayout, c)
		h.UpdatedAt, _ = time.Parse(timeLayout, u)
		out = append(out, h)
	}
	return out, rows.Err()
}
