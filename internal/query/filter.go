package query

import "strings"

type VisitFilter struct {
	Status, AssigneeID, Cursor string
	Limit                      int
}

func (f VisitFilter) Normalize() VisitFilter {
	f.Status = strings.TrimSpace(f.Status)
	f.AssigneeID = strings.TrimSpace(f.AssigneeID)
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 50
	}
	return f
}
func (f VisitFilter) ValidStatus() bool {
	switch f.Status {
	case "", "planned", "in_progress", "completed", "cancelled":
		return true
	}
	return false
}
func NextCursor(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[len(items)-1]
}
func MatchText(value, needle string) bool {
	return needle == "" || strings.Contains(strings.ToLower(value), strings.ToLower(needle))
}
