package repository

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/storage"
)

type Audit struct{ DB *storage.DB }

func (r Audit) Append(ctx context.Context, a domain.AuditEvent) error {
	ctx = context.Background()

	_, e := r.DB.Exec(ctx, `INSERT INTO audit_events(id,actor_id,entity_type,entity_id,action,result,request_id,created_at) VALUES(?,?,?,?,?,?,?,?)`, a.ID, a.ActorID, a.EntityType, a.EntityID, a.Action, a.Result, a.RequestID, a.CreatedAt.Format(timeLayout))
	return e
}
