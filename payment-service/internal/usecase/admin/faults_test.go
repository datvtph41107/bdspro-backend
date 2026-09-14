package admin

import (
	"testing"

	"common/fault"
)

func TestNormalizePageReturnsCanonicalPageSizeFault(t *testing.T) {
	_, _, err := normalizePage(1, 101)
	if err == nil {
		t.Fatal("normalizePage returned nil error")
	}
	failure, ok := fault.As(err)
	if !ok {
		t.Fatalf("error type = %T, want canonical fault", err)
	}
	if failure.Kind() != fault.KindValidation {
		t.Fatalf("kind = %q, want %q", failure.Kind(), fault.KindValidation)
	}
	if failure.Code() != "payment.admin.page_size_out_of_range" {
		t.Fatalf("code = %q", failure.Code())
	}
	violations := failure.Violations()
	if len(violations) != 1 || violations[0].Field != "page_size" {
		t.Fatalf("violations = %+v", violations)
	}
}
