package query

import "testing"

func TestNormalizeFilter(t *testing.T) {
	f := VisitFilter{Status: " planned ", Limit: 0}.Normalize()
	if f.Status != "planned" || f.Limit != 50 {
		t.Fatal(f)
	}
	f = VisitFilter{Limit: 1000}.Normalize()
	if f.Limit != 50 {
		t.Fatal("limit cap")
	}
}
func TestValidStatus(t *testing.T) {
	for _, s := range []string{"", "planned", "in_progress", "completed", "cancelled"} {
		if !(VisitFilter{Status: s}).ValidStatus() {
			t.Fatal(s)
		}
	}
	if (VisitFilter{Status: "unknown"}).ValidStatus() {
		t.Fatal("unknown accepted")
	}
}
func TestCursor(t *testing.T) {
	if NextCursor(nil) != "" {
		t.Fatal("empty cursor")
	}
	if NextCursor([]string{"a", "b"}) != "b" {
		t.Fatal("last cursor")
	}
}
func TestMatch(t *testing.T) {
	if !MatchText("House A", "house") {
		t.Fatal("case insensitive mismatch")
	}
	if MatchText("A", "z") {
		t.Fatal("false match")
	}
}
