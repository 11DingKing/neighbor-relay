package audit

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/repository"
	"time"
)

type Recorder struct{ Store repository.Audit }

func (r Recorder) Record(ctx context.Context, actor, entityType, entityID, action, result, requestID string) error {
	return r.Store.Append(ctx, domain.AuditEvent{ID: time.Now().UTC().Format("20060102150405.000000000"), ActorID: actor, EntityType: entityType, EntityID: entityID, Action: action, Result: result, RequestID: requestID, CreatedAt: time.Now().UTC()})
}
func Outcome(ok bool) string {
	if ok {
		return "success"
	}
	return "failure"
}
