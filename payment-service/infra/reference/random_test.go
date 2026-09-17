package reference

import (
	"testing"

	domain "payment/internal/domain/payment"
)

func TestNewOrderReferenceUsesCanonicalUppercase(t *testing.T) {
	reference, err := New().NewOrderReference()
	if err != nil {
		t.Fatalf("NewOrderReference() error = %v", err)
	}
	if !domain.IsCanonicalOrderReference(reference) {
		t.Fatalf("reference = %q", reference)
	}
}
