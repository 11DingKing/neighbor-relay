package domain

import (
	"testing"
	"time"
)

func TestCaseTransitions(t *testing.T) {
	tests := []struct {
		name     string
		from, to CaseStatus
		ok       bool
	}{{"open progress", CaseOpen, CaseInProgress, true}, {"open closed", CaseOpen, CaseClosed, true}, {"open resolved", CaseOpen, CaseResolved, false}, {"progress resolved", CaseInProgress, CaseResolved, true}, {"progress open", CaseInProgress, CaseOpen, true}, {"resolved closed", CaseResolved, CaseClosed, true}, {"closed open", CaseClosed, CaseOpen, false}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.from.CanTransition(tt.to); got != tt.ok {
				t.Fatalf("got %v", got)
			}
		})
	}
}
func TestVisitTransitions(t *testing.T) {
	tests := []struct {
		from, to VisitStatus
		ok       bool
	}{{VisitPlanned, VisitInProgress, true}, {VisitPlanned, VisitCancelled, true}, {VisitPlanned, VisitCompleted, false}, {VisitInProgress, VisitCompleted, true}, {VisitInProgress, VisitCancelled, true}, {VisitCompleted, VisitPlanned, false}, {VisitCancelled, VisitInProgress, false}}
	for _, tt := range tests {
		if tt.from.CanTransition(tt.to) != tt.ok {
			t.Errorf("%s to %s", tt.from, tt.to)
		}
	}
}
func TestValidation(t *testing.T) {
	if ValidateHousehold("", "x") == nil {
		t.Fatal("empty address accepted")
	}
	if ValidateHousehold("road", "x") != nil {
		t.Fatal("valid household rejected")
	}
	if ValidateCase("", "") == nil {
		t.Fatal("empty case accepted")
	}
	if ValidateNote(" ") == nil {
		t.Fatal("empty note accepted")
	}
	if ValidateVisit(time.Time{}) == nil {
		t.Fatal("zero visit accepted")
	}
}
func TestErrorIdentity(t *testing.T) {
	if ErrConflict == ErrInvalid {
		t.Fatal("sentinels collide")
	}
}
