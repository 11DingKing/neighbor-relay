package service

import (
	"github.com/11DingKing/neighbor-relay/internal/domain"
	"testing"
)

func TestEmailValidation(t *testing.T) {
	for _, v := range []string{"a@example.com", "person+tag@example.org"} {
		if ValidateEmail(v) != nil {
			t.Fatal(v)
		}
	}
	for _, v := range []string{"", "bad", "@x"} {
		if ValidateEmail(v) == nil {
			t.Fatal(v)
		}
	}
}
func TestPhoneValidation(t *testing.T) {
	for _, v := range []string{"123456", "+86 138-0000-0000", "(123) 456"} {
		if ValidatePhone(v) != nil {
			t.Fatal(v)
		}
	}
	for _, v := range []string{"", "abc", "123"} {
		if ValidatePhone(v) == nil {
			t.Fatal(v)
		}
	}
}
func TestNormalization(t *testing.T) {
	if NormalizeName("  A   B ") != "A B" {
		t.Fatal("name")
	}
	if NormalizeAddress(" a\nb ") != "a b" {
		t.Fatal("address")
	}
	v := " x "
	TrimOptional(&v)
	if v != "x" {
		t.Fatal(v)
	}
}
func TestRanges(t *testing.T) {
	for i := 0; i <= 10; i++ {
		if ValidatePriority(i) != nil {
			t.Fatal(i)
		}
	}
	if ValidatePriority(-1) == nil || ValidatePriority(11) == nil {
		t.Fatal("priority")
	}
	if ValidatePage(10, 0) != nil || ValidatePage(0, 0) == nil || ValidatePage(10, -1) == nil {
		t.Fatal("page")
	}
}
func TestStatusValidation(t *testing.T) {
	for _, v := range []string{"open", "in_progress", "resolved", "closed"} {
		if ValidateStatus(v) != nil {
			t.Fatal(v)
		}
	}
	if ValidateStatus("bad") == nil {
		t.Fatal("bad case")
	}
	for _, v := range []string{"planned", "in_progress", "completed", "cancelled"} {
		if ValidateVisitStatus(v) != nil {
			t.Fatal(v)
		}
	}
	if ValidateVisitStatus("bad") == nil {
		t.Fatal("bad visit")
	}
}
func TestIDs(t *testing.T) {
	if !ValidID("12345678") || ValidID("short") || ValidID("bad id") {
		t.Fatal("id")
	}
}
func TestCredentialValidation(t *testing.T) {
	if ValidateCredentials("a@x.com", "12345678") != nil {
		t.Fatal("valid credentials")
	}
	if ValidateCredentials("bad", "12345678") == nil || ValidateCredentials("a@x.com", "short") == nil {
		t.Fatal("invalid credentials")
	}
}
func TestChooseMessage(t *testing.T) {
	if ChooseMessage(nil) != "ok" || ChooseMessage(domain.ErrNotFound) != "not found" || ChooseMessage(domain.ErrConflict) != "conflict" || ChooseMessage(domain.ErrForbidden) != "forbidden" {
		t.Fatal("messages")
	}
}
