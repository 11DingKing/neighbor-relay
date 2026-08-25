package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/storage"
	"time"
)

type AuditQuery struct{ DB *storage.DB }

func (r AuditQuery) ForEntity(ctx context.Context, typ, id string, limit int) ([]domain.AuditEvent, error) {
	if limit < 1 || limit > 100 {
		limit = 50
	}
	rows, e := r.DB.Query(ctx, `SELECT id,actor_id,entity_type,entity_id,action,result,request_id,created_at FROM audit_events WHERE entity_type=? AND entity_id=? ORDER BY created_at DESC LIMIT ?`, typ, id, limit)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := make([]domain.AuditEvent, 0)
	for rows.Next() {
		var a domain.AuditEvent
		var created string
		if e := rows.Scan(&a.ID, &a.ActorID, &a.EntityType, &a.EntityID, &a.Action, &a.Result, &a.RequestID, &created); e != nil {
			return nil, e
		}
		a.CreatedAt, _ = time.Parse(timeLayout, created)
		out = append(out, a)
	}
	return out, rows.Err()
}
func (r AuditQuery) Count(ctx context.Context, typ string) (int, error) {
	var n int
	e := r.DB.QueryRow(ctx, `SELECT COUNT(*) FROM audit_events WHERE entity_type=?`, typ).Scan(&n)
	return n, e
}
func optionalString(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}
