package domain

import (
	"strings"
	"time"
)

func ValidateHousehold(address, contact string) error {
	if strings.TrimSpace(address) == "" || len(address) > 240 || strings.TrimSpace(contact) == "" {
		return ErrInvalid
	}
	return nil
}
func ValidateCase(title, summary string) error {
	if strings.TrimSpace(title) == "" || len(title) > 160 || len(summary) > 4000 {
		return ErrInvalid
	}
	return nil
}
func ValidateVisit(t time.Time) error {
	if t.IsZero() || t.Before(time.Now().Add(-365*24*time.Hour)) {
		return ErrInvalid
	}
	return nil
}
func ValidateNote(body string) error {
	n := strings.TrimSpace(body)
	if n == "" || len(n) > 8000 {
		return ErrInvalid
	}
	return nil
}
