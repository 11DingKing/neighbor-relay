package service

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"testing"
	"time"
)

func TestCancelledContext(t *testing.T) {
	s := testService(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, e := s.CreateHousehold(ctx, admin(), "road", "A", "p", 0)
	if e == nil {
		t.Fatal("cancelled context ignored")
	}
}
func TestCaseConflict(t *testing.T) {
	s := testService(t)
	a := admin()
	h, _ := s.CreateHousehold(context.Background(), a, "road", "A", "p", 0)
	c, _ := s.CreateCase(context.Background(), a, h.ID, "case", "summary")
	if e := s.TransitionCase(context.Background(), a, c.ID, domain.CaseClosed); e != nil {
		t.Fatal(e)
	}
	if e := s.TransitionCase(context.Background(), a, c.ID, domain.CaseOpen); e == nil {
		t.Fatal("closed case reopened")
	}
}
func TestVisitValidation(t *testing.T) {
	s := testService(t)
	a := admin()
	h, _ := s.CreateHousehold(context.Background(), a, "road", "A", "p", 0)
	c, _ := s.CreateCase(context.Background(), a, h.ID, "case", "summary")
	if _, e := s.PlanVisit(context.Background(), a, c.ID, "w", time.Time{}, "", ""); e == nil {
		t.Fatal("zero schedule accepted")
	}
}
