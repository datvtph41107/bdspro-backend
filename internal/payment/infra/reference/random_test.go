package reference

import (
	"regexp"
	"testing"
)

func TestNewOrderReferenceUsesCanonicalUppercase(t *testing.T) {
	reference, err := New().NewOrderReference()
	if err != nil {
		t.Fatalf("NewOrderReference() error = %v", err)
	}
	if !regexp.MustCompile(`^QHP-[0-9A-F]{16}$`).MatchString(reference) {
		t.Fatalf("reference = %q", reference)
	}
}
