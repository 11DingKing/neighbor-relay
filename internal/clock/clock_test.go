package clock

import (
	"testing"
	"time"
)

func TestFixedClock(t *testing.T) {
	want := time.Date(2026, 8, 25, 9, 0, 0, 0, time.UTC)
	if (Fixed{Value: want}).Now() != want {
		t.Fatal("fixed clock changed")
	}
}
func TestUTC(t *testing.T) {
	loc := time.FixedZone("x", 3600)
	got := UTC(time.Date(2020, 1, 1, 1, 0, 0, 0, loc))
	if got.Location() != time.UTC {
		t.Fatal("not utc")
	}
}
func TestSameDay(t *testing.T) {
	a := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
	b := a.Add(20 * time.Hour)
	if !SameDay(a, b) {
		t.Fatal("same day false")
	}
	if SameDay(a, a.Add(48*time.Hour)) {
		t.Fatal("different day true")
	}
}
func TestStart(t *testing.T) {
	a := time.Date(2024, 1, 1, 22, 4, 0, 0, time.UTC)
	if AtStartOfDay(a).Hour() != 0 {
		t.Fatal("not start")
	}
}
