package policy

import "github.com/11DingKing/neighbor-relay/internal/domain"

// CanCreateCase centralizes the privileged case-creation boundary.
func CanCreateCase(u domain.User) bool {
	return u.Active && (u.Role == domain.RoleAdmin || u.Role == domain.RoleFieldWorker)
}
func CanPlanVisit(u domain.User) bool { return u.Active && u.Role == domain.RoleAdmin }
func CanReadHouseholds(u domain.User) bool {
	return u.Active && (u.Role == domain.RoleAdmin || u.Role == domain.RoleFieldWorker)
}
func CanOperateVisit(u domain.User, v domain.Visit) bool {
	return u.Active && (u.Role == domain.RoleAdmin || (u.Role == domain.RoleFieldWorker && v.AssigneeID == u.ID))
}
func RoleLabel(r domain.Role) string {
	if r == domain.RoleAdmin {
		return "administrator"
	}
	if r == domain.RoleFieldWorker {
		return "field worker"
	}
	return "unknown"
}
