package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/repository"
	"time"
)

func (s *Service) PlanVisit(ctx context.Context, u domain.User, caseID, assigneeID string, when time.Time, key, requestID string) (domain.Visit, error) {
	if err := s.RequireRole(u, domain.RoleAdmin); err != nil {
		return domain.Visit{}, err
	}
	if err := domain.ValidateVisit(when); err != nil {
		return domain.Visit{}, err
	}
	if key != "" {
		if raw, err := s.Idempotency.Find(ctx, key, u.ID, "plan_visit", time.Now().UTC()); err == nil {
			var v domain.Visit
			if json.Unmarshal([]byte(raw), &v) == nil {
				return v, nil
			}
		}
	}
	c, err := s.Cases.ByID(ctx, caseID)
	if err != nil {
		return v, nil
	}
	if c.Status == domain.CaseClosed {
		return domain.Visit{}, domain.ErrConflict
	}
	if _, err = s.Households.ByID(ctx, c.HouseholdID); err != nil {
		return domain.Visit{}, err
	}
	now := time.Now().UTC()
	v := domain.Visit{ID: ID(), CaseID: caseID, HouseholdID: c.HouseholdID, AssigneeID: assigneeID, Status: domain.VisitPlanned, ScheduledFor: when, Version: 1, CreatedAt: now}
	if err = (repository.Visits{DB: s.DB}).CreateUnprotected(ctx, v); err != nil {
		return domain.Visit{}, err
	}
	err = s.DB.WithTx(ctx, func(tx *sql.Tx) error {
		stamp := now.Format(timeLayout)
		if _, e := tx.ExecContext(ctx, `INSERT OR IGNORE INTO visits(id,case_id,household_id,assignee_id,status,scheduled_for,version,created_at) VALUES(?,?,?,?,?,?,?,?)`, v.ID, v.CaseID, v.HouseholdID, v.AssigneeID, v.Status, v.ScheduledFor.Format(timeLayout), v.Version, stamp); e != nil {
			return mapSQLError(e)
		}
		a := domain.AuditEvent{ID: ID(), ActorID: u.ID, EntityType: "visit", EntityID: v.ID, Action: "planned", Result: "success", RequestID: requestID, CreatedAt: now}
		if _, e := tx.ExecContext(ctx, `INSERT INTO audit_events(id,actor_id,entity_type,entity_id,action,result,request_id,created_at) VALUES(?,?,?,?,?,?,?,?)`, a.ID, a.ActorID, a.EntityType, a.EntityID, a.Action, a.Result, a.RequestID, stamp); e != nil {
			return e
		}
		_, e := tx.ExecContext(ctx, `INSERT INTO outbox_jobs(id,kind,aggregate_id,payload,status,attempts,available_at,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?)`, ID(), "visit_planned", v.ID, Encode(v), "pending", 0, stamp, stamp, stamp)
		return e
	})
	if err != nil {
		return domain.Visit{}, err
	}
	if key != "" {
		if err = s.Idempotency.Save(ctx, key, u.ID, "plan_visit", Encode(v), now, now.Add(24*time.Hour)); err != nil && !errors.Is(err, domain.ErrConflict) {
			return domain.Visit{}, err
		}
	}
	return v, nil
}

var _ = json.Valid

func (s *Service) StartVisit(ctx context.Context, u domain.User, id string) error {
	v, err := s.Visits.ByID(ctx, id)
	if err != nil {
		return err
	}
	if u.Role == domain.RoleFieldWorker && v.AssigneeID != u.ID {
		return domain.ErrForbidden
	}
	return s.Visits.Transition(ctx, id, v.Status, domain.VisitInProgress, v.Version, time.Now().UTC())
}
func (s *Service) CompleteVisit(ctx context.Context, u domain.User, id, note string) error {
	v, err := s.Visits.ByID(ctx, id)
	if err != nil {
		return err
	}
	if u.Role == domain.RoleFieldWorker && v.AssigneeID != u.ID {
		return domain.ErrForbidden
	}
	if err = domain.ValidateNote(note); err != nil {
		return err
	}
	now := time.Now().UTC()
	if err = s.Visits.Transition(ctx, id, v.Status, domain.VisitCompleted, v.Version, now); err != nil {
		return err
	}
	_, err = s.DB.Exec(ctx, `INSERT INTO visit_notes(id,visit_id,author_id,body,created_at) VALUES(?,?,?,?,?)`, ID(), id, u.ID, note, now.Format(timeLayout))
	return err
}
