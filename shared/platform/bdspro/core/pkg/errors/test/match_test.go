package apperr

import (
	apperr "bdspro/core/pkg/errors"
	"testing"
)

func TestRetryableAndFaultKinds(t *testing.T) {
	e1 := apperr.New(apperr.KindExternal, apperr.CodeUnavailable, "service unavailable")
	if !apperr.IsRetryable(e1) {
		t.Fatal("external errors should be retryable by default")
	}
	if !apperr.IsSystemFault(e1) {
		t.Fatal("external must be system fault")
	}
	if apperr.IsClientFault(e1) {
		t.Fatal("external must not be client fault")
	}

	e2 := apperr.New(apperr.KindValidation, apperr.CodeInvalidArgument, "invalid")
	if !apperr.IsClientFault(e2) {
		t.Fatal("validation must be client fault")
	}
	if apperr.IsSystemFault(e2) {
		t.Fatal("validation must not be system fault")
	}
}
