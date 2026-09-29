package admin

import (
	"testing"

	_errors "common/errors"
)

func TestValidateDraftCommandReturnsCanonicalTierRankError(t *testing.T) {
	_, err := validateDraftCommand(DraftCommand{
		ActorID:            1,
		ProductDisplayName: "Professional",
		TierRank:           0,
	})
	assertValidationError(t, err, "USER_PLAN_TIER_RANK_MUST_BE_POSITIVE", "catalog.plan_version.tier_rank_positive", "tier_rank")
}

func TestNormalizeQueryReturnsCanonicalValidationErrors(t *testing.T) {
	testCases := []struct {
		name       string
		query      Query
		key        _errors.Key
		legacyCode string
		fieldName  string
	}{
		{"page size", Query{PageSize: 101}, "USER_PLAN_PAGE_SIZE_OUT_OF_RANGE", "catalog.plan_version.page_size_invalid", "page_size"},
		{"status", Query{Status: "unknown"}, "USER_PLAN_STATUS_INVALID", "catalog.plan_version.status_invalid", "status"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := normalizeQuery(tc.query)
			assertValidationError(t, err, tc.key, tc.legacyCode, tc.fieldName)
		})
	}
}

func TestValidateDraftCommandWrapsDomainValidationAsCanonicalError(t *testing.T) {
	_, err := validateDraftCommand(DraftCommand{
		ActorID:            1,
		ProductDisplayName: "Professional",
		TierRank:           1,
	})
	assertValidationError(t, err, "USER_PLAN_TERMS_INVALID", "catalog.plan_version.terms_invalid", "")
}

func assertValidationError(t *testing.T, err error, key _errors.Key, legacyCode, fieldName string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected validation error")
	}
	application, ok := _errors.As(err)
	if !ok {
		t.Fatalf("error type = %T, want canonical application error", err)
	}
	if application.Key() != key {
		t.Fatalf("key = %q, want %q", application.Key(), key)
	}
	if application.Spec().LegacyProblemCode() != legacyCode {
		t.Fatalf("legacy code = %q, want %q", application.Spec().LegacyProblemCode(), legacyCode)
	}
	if fieldName == "" {
		return
	}
	violations := application.Violations()
	if len(violations) != 1 || violations[0].Field != fieldName {
		t.Fatalf("violations = %+v, want field %q", violations, fieldName)
	}
}
