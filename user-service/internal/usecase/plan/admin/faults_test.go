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
	if err == nil {
		t.Fatal("validateDraftCommand returned nil error")
	}
	failure, ok := fault.As(err)
	if !ok {
		t.Fatalf("error type = %T, want canonical fault", err)
	}
	if failure.Kind() != fault.KindValidation {
		t.Fatalf("kind = %q, want %q", failure.Kind(), fault.KindValidation)
	}
	if failure.Code() != "catalog.plan_version.tier_rank_positive" {
		t.Fatalf("code = %q", failure.Code())
	}
	violations := failure.Violations()
	if len(violations) != 1 || violations[0].Field != "tier_rank" {
		t.Fatalf("violations = %+v", violations)
	}
}
