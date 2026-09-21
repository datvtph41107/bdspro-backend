package admin

import (
	"testing"

	_errors "common/errors"
)

func TestNormalizePageReturnsCanonicalPageSizeError(t *testing.T) {
	_, _, err := normalizePage(1, 101)
	if err == nil {
		t.Fatal("normalizePage returned nil error")
	}
	application, ok := _errors.As(err)
	if !ok {
		t.Fatalf("error type = %T, want canonical application error", err)
	}
	if application.Key() != "PAYMENT_ADMIN_PAGE_SIZE_OUT_OF_RANGE" {
		t.Fatalf("key = %q", application.Key())
	}
	if application.Spec().LegacyProblemCode() != "payment.admin.page_size_out_of_range" {
		t.Fatalf("legacy problem code = %q", application.Spec().LegacyProblemCode())
	}
	violations := application.Violations()
	if len(violations) != 1 || violations[0].Field != "page_size" {
		t.Fatalf("violations = %+v", violations)
	}
}
