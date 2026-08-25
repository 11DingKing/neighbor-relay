package policy

import (
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"testing"
)

func users() (domain.User, domain.User) {
	return domain.User{ID: "a", Role: domain.RoleAdmin, Active: true}, domain.User{ID: "f", Role: domain.RoleFieldWorker, Active: true}
}
func TestCasePolicy(t *testing.T) {
	a, f := users()
	if !CanCreateCase(a) || CanCreateCase(f) {
		t.Fatal("case policy")
	}
	a.Active = false
	if CanCreateCase(a) {
		t.Fatal("inactive admin")
	}
}
func TestVisitPolicy(t *testing.T) {
	_, f := users()
	v := domain.Visit{AssigneeID: f.ID}
	if !CanOperateVisit(f, v) {
		t.Fatal("assigned worker denied")
	}
	v.AssigneeID = "other"
	if CanOperateVisit(f, v) {
		t.Fatal("unassigned worker allowed")
	}
}
func TestReadPolicy(t *testing.T) {
	a, f := users()
	if !CanReadHouseholds(a) || !CanReadHouseholds(f) {
		t.Fatal("read denied")
	}
	if CanReadHouseholds(domain.User{Active: true}) {
		t.Fatal("unknown role read")
	}
}
func TestLabels(t *testing.T) {
	if RoleLabel(domain.RoleAdmin) == "" || RoleLabel(domain.RoleFieldWorker) == "" {
		t.Fatal("missing label")
	}
}
