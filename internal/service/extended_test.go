package service

import (
	"context"
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"github.com/11DingKing/neighbor-relay/internal/query"
	"testing"
	"time"
)

func TestVisitQueries(t *testing.T) {
	s := testService(t)
	a := admin()
	w := worker()
	h, _ := s.CreateHousehold(context.Background(), a, "road", "A", "p", 1)
	c, _ := s.CreateCase(context.Background(), a, h.ID, "case", "summary")
	for i := 0; i < 4; i++ {
		_, e := s.PlanVisit(context.Background(), a, c.ID, w.ID, time.Now().UTC().Add(time.Duration(i+1)*time.Hour), "key-"+string(rune('a'+i)), "r")
		if e != nil {
			t.Fatal(e)
		}
	}
	items, e := s.ListAssignedVisits(context.Background(), w, query.VisitFilter{Status: "planned", Limit: 3})
	if e != nil || len(items) != 3 {
		t.Fatalf("items=%d err=%v", len(items), e)
	}
	ok, e := s.OpenVisitExists(context.Background(), c.ID)
	if e != nil || !ok {
		t.Fatalf("open=%v err=%v", ok, e)
	}
}
func TestSummary(t *testing.T) {
	s := testService(t)
	a := admin()
	w := worker()
	h, _ := s.CreateHousehold(context.Background(), a, "road", "A", "p", 1)
	c, _ := s.CreateCase(context.Background(), a, h.ID, "case", "summary")
	v, _ := s.PlanVisit(context.Background(), a, c.ID, w.ID, time.Now().UTC().Add(time.Hour), "sum-key", "r")
	if e := s.StartVisit(context.Background(), w, v.ID); e != nil {
		t.Fatal(e)
	}
	sum, e := s.Summary(context.Background(), a)
	if e != nil || sum.InProgress != 1 {
		t.Fatalf("summary=%#v err=%v", sum, e)
	}
}
func TestNotes(t *testing.T) {
	s := testService(t)
	a := admin()
	w := worker()
	h, _ := s.CreateHousehold(context.Background(), a, "road", "A", "p", 1)
	c, _ := s.CreateCase(context.Background(), a, h.ID, "case", "summary")
	v, _ := s.PlanVisit(context.Background(), a, c.ID, w.ID, time.Now().UTC().Add(time.Hour), "note-key", "r")
	if e := s.AddNote(context.Background(), w, v.ID, "observed needs"); e != nil {
		t.Fatal(e)
	}
	if e := s.AddNote(context.Background(), w, v.ID, " "); e == nil {
		t.Fatal("empty note")
	}
}
func TestMaintenance(t *testing.T) {
	s := testService(t)
	n, e := s.Cleanup(context.Background())
	if e != nil || n < 0 {
		t.Fatalf("cleanup=%d err=%v", n, e)
	}
	if _, e = s.RecoverWorkers(context.Background(), time.Minute); e != nil {
		t.Fatal(e)
	}
}
func TestManyValidationCases(t *testing.T) {
	cases := []struct {
		name           string
		title, summary string
		valid          bool
	}{{"normal", "title", "summary", true}, {"spaces", " ", "summary", false}, {"long summary", "title", string(make([]byte, 4001)), false}, {"unicode", "关怀", "需要回访", true}, {"empty", "", "", false}, {"long title", string(make([]byte, 161)), "s", false}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			e := domain.ValidateCase(c.title, c.summary)
			if (e == nil) != c.valid {
				t.Fatalf("valid=%v err=%v", c.valid, e)
			}
		})
	}
}
func TestContextDeadline(t *testing.T) {
	s := testService(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	time.Sleep(time.Millisecond)
	_, e := s.ListHouseholds(ctx, admin(), 10)
	if e == nil {
		t.Fatal("deadline ignored")
	}
}
func TestUnauthorizedTransitions(t *testing.T) {
	s := testService(t)
	a := admin()
	w := worker()
	h, _ := s.CreateHousehold(context.Background(), a, "road", "A", "p", 0)
	c, _ := s.CreateCase(context.Background(), a, h.ID, "case", "summary")
	if e := s.TransitionCase(context.Background(), w, c.ID, domain.CaseInProgress); e != domain.ErrForbidden {
		t.Fatalf("worker changed case: %v", e)
	}
}
func TestRepeatedVisitCompletion(t *testing.T) {
	s := testService(t)
	a := admin()
	w := worker()
	h, _ := s.CreateHousehold(context.Background(), a, "road", "A", "p", 0)
	c, _ := s.CreateCase(context.Background(), a, h.ID, "case", "summary")
	v, _ := s.PlanVisit(context.Background(), a, c.ID, w.ID, time.Now().UTC().Add(time.Hour), "repeat-key", "r")
	_ = s.StartVisit(context.Background(), w, v.ID)
	if e := s.CompleteVisit(context.Background(), w, v.ID, "done"); e != nil {
		t.Fatal(e)
	}
	if e := s.CompleteVisit(context.Background(), w, v.ID, "again"); e == nil {
		t.Fatal("completed twice")
	}
}
func TestCaseOpenVisitGuard(t *testing.T) {
	s := testService(t)
	a := admin()
	h, _ := s.CreateHousehold(context.Background(), a, "road", "A", "p", 0)
	c, _ := s.CreateCase(context.Background(), a, h.ID, "case", "summary")
	if _, e := s.PlanVisit(context.Background(), a, c.ID, worker().ID, time.Now().UTC().Add(time.Hour), "guard-key", "r"); e != nil {
		t.Fatal(e)
	}
}
func TestListLimits(t *testing.T) {
	s := testService(t)
	if _, e := s.ListAssignedVisits(context.Background(), domain.User{Role: domain.RoleAdmin, Active: true, ID: "admin"}, query.VisitFilter{Limit: -1}); e != nil {
		t.Fatal(e)
	}
	if _, e := s.ListAssignedVisits(context.Background(), domain.User{Role: domain.RoleAdmin, Active: true, ID: "admin"}, query.VisitFilter{Status: "bad"}); e != domain.ErrInvalid {
		t.Fatalf("bad status err=%v", e)
	}
}
func TestUserInactive(t *testing.T) {
	s := testService(t)
	u := domain.User{ID: "inactive", Email: "inactive@x", Role: domain.RoleAdmin, Active: false}
	if e := s.EnsureUser(context.Background(), u, "pw"); e != nil {
		t.Fatal(e)
	}
	if _, _, e := s.Login(context.Background(), u.Email, "pw"); e != domain.ErrUnauthorized {
		t.Fatalf("inactive login=%v", e)
	}
}
func TestCaseVersionConflict(t *testing.T) {
	s := testService(t)
	a := admin()
	h, _ := s.CreateHousehold(context.Background(), a, "road", "A", "p", 0)
	c, _ := s.CreateCase(context.Background(), a, h.ID, "case", "summary")
	if e := s.Cases.Transition(context.Background(), c.ID, c.Status, domain.CaseInProgress, 99, time.Now()); e != domain.ErrConflict {
		t.Fatalf("conflict=%v", e)
	}
}
