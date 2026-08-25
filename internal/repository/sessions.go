package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"time"
)

type Sessions struct{ DB *storage.DB }

func (r Sessions) Create(ctx context.Context, s domain.Session) error {
	_, e := r.DB.Exec(ctx, `INSERT INTO sessions(id,user_id,token_hash,expires_at,created_at) VALUES(?,?,?,?,?)`, s.ID, s.UserID, s.TokenHash, s.ExpiresAt.Format(timeLayout), s.CreatedAt.Format(timeLayout))
	return e
}
func (r Sessions) Active(ctx context.Context, hash string, now time.Time) (domain.Session, error) {
	var s domain.Session
	var exp, created string
	var revoked sql.NullString
	e := r.DB.QueryRow(ctx, `SELECT id,user_id,token_hash,expires_at,created_at,revoked_at FROM sessions WHERE token_hash=?`, hash).Scan(&s.ID, &s.UserID, &s.TokenHash, &exp, &created, &revoked)
	if e == sql.ErrNoRows {
		return s, domain.ErrUnauthorized
	}
	if e != nil {
		return s, e
	}
	et, _ := time.Parse(timeLayout, exp)
	ct, _ := time.Parse(timeLayout, created)
	s.ExpiresAt = &et
	s.CreatedAt = &ct
	if revoked.Valid || !et.After(now) {
		return s, domain.ErrExpired
	}
	if revoked.Valid {
		t, _ := time.Parse(timeLayout, revoked.String)
		s.RevokedAt = &t
	}
	return s, nil
}
func (r Sessions) Revoke(ctx context.Context, hash string, at time.Time) error {
	res, e := r.DB.Exec(ctx, `UPDATE sessions SET revoked_at=? WHERE token_hash=? AND revoked_at IS NULL`, at.Format(timeLayout), hash)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
