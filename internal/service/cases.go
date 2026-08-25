package service

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"time"
)

func (s *Service) CreateCase(ctx context.Context, u domain.User, householdID, title, summary string) (domain.Case, error) {
	if e := s.RequireRole(u, domain.RoleAdmin); e != nil {
		return domain.Case{}, e
	}
	if e := domain.ValidateCase(title, summary); e != nil {
		return domain.Case{}, e
	}
	if _, e := s.Households.ByID(ctx, householdID); e != nil {
		return domain.Case{}, e
	}
	now := time.Now().UTC()
	c := domain.Case{ID: ID(), HouseholdID: householdID, Title: title, Summary: summary, Status: domain.CaseOpen, OwnerID: u.ID, Version: 1, CreatedAt: now, UpdatedAt: now}
	if e := s.Cases.Create(ctx, c); e != nil {
		return domain.Case{}, mapSQLError(e)
	}
	return c, nil
}
func (s *Service) TransitionCase(ctx context.Context, u domain.User, id string, to domain.CaseStatus) error {
	c, e := s.Cases.ByID(ctx, id)
	if e != nil {
		return e
	}
	if u.Role != domain.RoleAdmin && c.OwnerID != u.ID {
		return domain.ErrForbidden
	}
	return s.Cases.Transition(ctx, id, c.Status, to, c.Version, time.Now().UTC())
}

func caseTransitionAudit(status domain.CaseStatus, at time.Time) string {
	if status == domain.CaseClosed {
		return at.UTC().Format(time.RFC3339)
	}
	return string(status)
}

var _ = caseTransitionAudit
