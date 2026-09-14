package admin

import (
	"testing"

	"common/fault"
)

func TestValidateDraftCommandReturnsCanonicalTierRankFault(t *testing.T) {
	_, err := validateDraftCommand(DraftCommand{
		ActorID:            1,
		ProductDisplayName: "Professional",
		TierRank:           0,
	})
	assertValidationFault(t, err, "catalog.plan_version.tier_rank_positive", "tier_rank")
}

func TestNormalizeQueryReturnsCanonicalValidationFaults(t *testing.T) {
	testCases := []struct {
		name      string
		query     Query
		code      string
		fieldName string
	}{
		{
			name:      "page size",
			query:     Query{PageSize: 101},
			code:      "catalog.plan_version.page_size_invalid",
			fieldName: "page_size",
		},
		{
			name:      "status",
			query:     Query{Status: "unknown"},
			code:      "catalog.plan_version.status_invalid",
			fieldName: "status",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := normalizeQuery(testCase.query)
			assertValidationFault(t, err, testCase.code, testCase.fieldName)
		})
	}
}

func TestValidateDraftCommandWrapsDomainValidationAsCanonicalFault(t *testing.T) {
	_, err := validateDraftCommand(DraftCommand{
		ActorID:            1,
		ProductDisplayName: "Professional",
		TierRank:           1,
	})
	assertValidationFault(t, err, "catalog.plan_version.terms_invalid", "")
}

func assertValidationFault(t *testing.T, err error, code, fieldName string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected validation error")
	}
	failure, ok := fault.As(err)
	if !ok {
		t.Fatalf("error type = %T, want canonical fault", err)
	}
	if failure.Kind() != fault.KindValidation {
		t.Fatalf("kind = %q, want %q", failure.Kind(), fault.KindValidation)
	}
	if failure.Code() != code {
		t.Fatalf("code = %q, want %q", failure.Code(), code)
	}
	if fieldName == "" {
		return
	}
	violations := failure.Violations()
	if len(violations) != 1 || violations[0].Field != fieldName {
		t.Fatalf("violations = %+v, want field %q", violations, fieldName)
	}
}
