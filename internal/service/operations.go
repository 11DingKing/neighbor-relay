package service

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"time"
)

type OperationResult struct {
	ID      string
	Changed bool
	At      time.Time
	Message string
}

func (s *Service) CancelVisit(ctx context.Context, u domain.User, id string) (OperationResult, error) {
	v, e := s.Visits.ByID(ctx, id)
	if e != nil {
		return OperationResult{}, e
	}
	if u.Role == domain.RoleFieldWorker && v.AssigneeID != u.ID {
		return OperationResult{}, domain.ErrForbidden
	}
	now := time.Now().UTC()
	if e = s.Visits.Transition(ctx, id, v.Status, domain.VisitCancelled, v.Version, now); e != nil {
		return OperationResult{}, e
	}
	return OperationResult{ID: id, Changed: true, At: now, Message: "visit cancelled"}, nil
}
func (s *Service) ReopenCase(ctx context.Context, u domain.User, id string) (OperationResult, error) {
	c, e := s.Cases.ByID(ctx, id)
	if e != nil {
		return OperationResult{}, e
	}
	if u.Role != domain.RoleAdmin && c.OwnerID != u.ID {
		return OperationResult{}, domain.ErrForbidden
	}
	now := time.Now().UTC()
	if e = s.Cases.Transition(ctx, id, c.Status, domain.CaseOpen, c.Version, now); e != nil {
		return OperationResult{}, e
	}
	return OperationResult{ID: id, Changed: true, At: now, Message: "case reopened"}, nil
}
func (s *Service) AssignCase(ctx context.Context, u domain.User, id, owner string) (OperationResult, error) {
	if e := s.RequireRole(u, domain.RoleAdmin); e != nil {
		return OperationResult{}, e
	}
	if !ValidID(owner) {
		return OperationResult{}, domain.ErrInvalid
	}
	if _, e := s.Users.ByID(ctx, owner); e != nil {
		return OperationResult{}, e
	}
	now := time.Now().UTC()
	res, e := s.DB.Exec(ctx, `UPDATE cases SET owner_id=?,version=version+1,updated_at=? WHERE id=?`, owner, now.Format(timeLayout), id)
	if e != nil {
		return OperationResult{}, e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return OperationResult{}, domain.ErrNotFound
	}
	return OperationResult{ID: id, Changed: true, At: now, Message: "case assigned"}, nil
}
func (s *Service) TouchHousehold(ctx context.Context, id string) (OperationResult, error) {
	now := time.Now().UTC()
	res, e := s.DB.Exec(ctx, `UPDATE households SET updated_at=? WHERE id=?`, now.Format(timeLayout), id)
	if e != nil {
		return OperationResult{}, e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return OperationResult{}, domain.ErrNotFound
	}
	return OperationResult{ID: id, Changed: true, At: now, Message: "household touched"}, nil
}
