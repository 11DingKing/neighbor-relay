package service

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/query"
	"github.com/11DingKing/neighbor-relay/internal/repository"
)

func (s *Service) ListAssignedVisits(ctx context.Context, u domain.User, filter query.VisitFilter) ([]domain.Visit, error) {
	if e := s.RequireRole(u, domain.RoleAdmin, domain.RoleFieldWorker); e != nil {
		return nil, e
	}
	f := filter.Normalize()
	if !f.ValidStatus() {
		return nil, domain.ErrInvalid
	}
	id := u.ID
	if u.Role == domain.RoleAdmin && f.AssigneeID != "" {
		id = f.AssigneeID
	}
	return (repository.VisitQuery{DB: s.DB}).ForAssignee(ctx, id, f.Status, f.Limit)
}
func (s *Service) OpenVisitExists(ctx context.Context, caseID string) (bool, error) {
	return (repository.VisitQuery{DB: s.DB}).HasOpenForCase(ctx, caseID)
}
