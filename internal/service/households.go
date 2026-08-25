package service

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"time"
)

func (s *Service) CreateHousehold(ctx context.Context, u domain.User, address, contact, phone string, priority int) (domain.Household, error) {
	if e := s.RequireRole(u, domain.RoleAdmin, domain.RoleFieldWorker); e != nil {
		return domain.Household{}, e
	}
	if e := domain.ValidateHousehold(address, contact); e != nil {
		return domain.Household{}, e
	}
	now := time.Now().UTC()
	h := domain.Household{ID: ID(), Address: address, ContactName: contact, Phone: phone, Priority: priority, CreatedAt: now, UpdatedAt: now}
	if e := s.Households.Create(ctx, h); e != nil {
		return domain.Household{}, mapSQLError(e)
	}
	return h, nil
}
func (s *Service) ListHouseholds(ctx context.Context, u domain.User, limit int) ([]domain.Household, error) {
	if e := s.RequireRole(u, domain.RoleAdmin, domain.RoleFieldWorker); e != nil {
		return nil, e
	}
	return s.Households.List(ctx, limit)
}
