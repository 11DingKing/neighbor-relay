package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"time"
)

type Idempotency struct{ DB *storage.DB }

func (r Idempotency) Find(ctx context.Context, key, user, op string, now time.Time) (string, error) {
	var response, expires string
	e := r.DB.QueryRow(ctx, `SELECT response,expires_at FROM idempotency_keys WHERE key=? AND user_id=? AND operation=?`, key, user, op).Scan(&response, &expires)
	if e == sql.ErrNoRows {
		return "", domain.ErrNotFound
	}
	if e != nil {
		return "", e
	}
	t, _ := time.Parse(timeLayout, expires)
	if !t.After(now) {
		return "", domain.ErrNotFound
	}
	return response, nil
}
func (r Idempotency) Save(ctx context.Context, key, user, op, response string, now, expires time.Time) error {
	_, e := r.DB.Exec(ctx, `INSERT INTO idempotency_keys(key,user_id,operation,response,created_at,expires_at) VALUES(?,?,?,?,?,?)`, key, user, op, response, now.Format(timeLayout), expires.Format(timeLayout))
	return e
}
