package service

import (
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

var phonePattern = regexp.MustCompile(`^[0-9+() -]{6,32}$`)

func ValidateEmail(v string) error {
	if _, e := mail.ParseAddress(strings.TrimSpace(v)); e != nil {
		return domain.ErrInvalid
	}
	return nil
}
func ValidatePhone(v string) error {
	if !phonePattern.MatchString(strings.TrimSpace(v)) {
		return domain.ErrInvalid
	}
	return nil
}
func NormalizeName(v string) string    { return strings.Join(strings.Fields(v), " ") }
func NormalizeAddress(v string) string { return strings.TrimSpace(strings.ReplaceAll(v, "\n", " ")) }
func ValidatePriority(v int) error {
	if v < 0 || v > 10 {
		return domain.ErrInvalid
	}
	return nil
}
func ValidatePage(limit, offset int) error {
	if limit < 1 || limit > 100 || offset < 0 {
		return domain.ErrInvalid
	}
	return nil
}
func ValidateStatus(status string) error {
	switch status {
	case string(domain.CaseOpen), string(domain.CaseInProgress), string(domain.CaseResolved), string(domain.CaseClosed):
		return nil
	}
	return domain.ErrInvalid
}
func ValidateVisitStatus(status string) error {
	switch status {
	case string(domain.VisitPlanned), string(domain.VisitInProgress), string(domain.VisitCompleted), string(domain.VisitCancelled):
		return nil
	}
	return domain.ErrInvalid
}
func ValidID(v string) bool { return len(v) >= 8 && len(v) <= 64 && !strings.ContainsAny(v, " /\\") }
func ValidateCredentials(email, password string) error {
	if ValidateEmail(email) != nil || len(password) < 8 {
		return domain.ErrInvalid
	}
	return nil
}
func TrimOptional(v *string) {
	if v != nil {
		*v = strings.TrimSpace(*v)
	}
}
func ChooseMessage(err error) string {
	if err == nil {
		return "ok"
	}
	if err == domain.ErrNotFound {
		return "not found"
	}
	if err == domain.ErrConflict {
		return "conflict"
	}
	if err == domain.ErrForbidden {
		return "forbidden"
	}
	return "request failed"
}
func EnsureRole(role domain.Role) error {
	if role != domain.RoleAdmin && role != domain.RoleFieldWorker {
		return domain.ErrInvalid
	}
	return nil
}
func EnsureActive(u domain.User) error {
	if !u.Active {
		return domain.ErrUnauthorized
	}
	return nil
}
func SafeLimit(limit int) int {
	if limit < 1 {
		return 50
	}
	if limit > 100 {
		return 100
	}
	return limit
}
func SafeOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}
func IsTerminalCase(s domain.CaseStatus) bool { return s == domain.CaseClosed }
func IsTerminalVisit(s domain.VisitStatus) bool {
	return s == domain.VisitCompleted || s == domain.VisitCancelled
}
func CanEditNote(u domain.User, v domain.Visit) bool {
	return u.Role == domain.RoleAdmin || (u.Role == domain.RoleFieldWorker && u.ID == v.AssigneeID)
}
func NonEmpty(values ...string) bool {
	for _, v := range values {
		if strings.TrimSpace(v) == "" {
			return false
		}
	}
	return true
}
func LengthBetween(v string, min, max int) bool { n := len([]rune(v)); return n >= min && n <= max }
func ValidateSummary(v string) error {
	if !LengthBetween(v, 1, 4000) {
		return domain.ErrInvalid
	}
	return nil
}
func ValidateTitle(v string) error {
	if !LengthBetween(v, 1, 160) {
		return domain.ErrInvalid
	}
	return nil
}
func ClampPriority(v int) int {
	if v < 0 {
		return 0
	}
	if v > 10 {
		return 10
	}
	return v
}
func HasPrefixID(id, prefix string) bool { return strings.HasPrefix(id, prefix) }
func JoinErrors(primary, secondary error) error {
	if primary != nil {
		return primary
	}
	return secondary
}
func ZeroTime(t time.Time) bool       { return t.IsZero() }
func IsFuture(t time.Time) bool       { return t.After(time.Now().UTC()) }
func IsPast(t time.Time) bool         { return t.Before(time.Now().UTC()) }
func NormalizeStatus(v string) string { return strings.ToLower(strings.TrimSpace(v)) }
func EqualFold(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func Between(v, min, max int) bool   { return v >= min && v <= max }
func Present(v *string) bool         { return v != nil && strings.TrimSpace(*v) != "" }
func CanonicalEmail(v string) string { return strings.ToLower(strings.TrimSpace(v)) }
func CanonicalPhone(v string) string { return strings.ReplaceAll(strings.TrimSpace(v), " ", "") }
