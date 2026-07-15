package seo_strategy

import (
	"testing"
	"time"
)

func TestEvaluateIntentHypothesisReadinessRequiresSupportingEvidence(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	need := &SeoUserNeed{ID: 20, ObjectiveID: 10, Status: UserNeedStatusActive}
	intent := validIntentHypothesisFixture()

	result := EvaluateIntentHypothesisReadiness(need, intent, nil, now)
	if result.Ready {
		t.Fatal("intent without supporting evidence must not be ready")
	}
	if !hasIntentIssue(result.Blockers, "missing_supporting_query_evidence") {
		t.Fatalf("missing blocker: %+v", result.Blockers)
	}
}

func TestEvaluateIntentHypothesisReadinessWithObservedSupport(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	review := now.Add(30 * 24 * time.Hour)
	need := &SeoUserNeed{ID: 20, ObjectiveID: 10, Status: UserNeedStatusActive}
	intent := validIntentHypothesisFixture()
	intent.ReviewDueAt = &review

	evidence := queryEvidenceFixture(now)
	linked := []LinkedIntentQueryEvidence{
		{
			Link:     SeoIntentEvidenceLink{RelationType: IntentEvidenceSupports},
			Evidence: evidence,
		},
		{
			Link:     SeoIntentEvidenceLink{RelationType: IntentEvidenceContext},
			Evidence: evidence,
		},
	}
	result := EvaluateIntentHypothesisReadiness(need, intent, linked, now)
	if !result.Ready {
		t.Fatalf("expected ready, blockers=%+v", result.Blockers)
	}
}

func TestEvaluateIntentHypothesisReadinessValidatedNeedsValidatedObservedEvidence(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	need := &SeoUserNeed{ID: 20, ObjectiveID: 10, Status: UserNeedStatusActive}
	intent := validIntentHypothesisFixture()
	intent.Confidence = ObjectiveConfidenceValidated

	evidence := queryEvidenceFixture(now)
	evidence.Confidence = ObjectiveConfidenceSupported
	result := EvaluateIntentHypothesisReadiness(
		need,
		intent,
		[]LinkedIntentQueryEvidence{{
			Link:     SeoIntentEvidenceLink{RelationType: IntentEvidenceSupports},
			Evidence: evidence,
		}},
		now,
	)
	if !hasIntentIssue(result.Blockers, "validated_without_validated_evidence") {
		t.Fatalf("missing validated evidence blocker: %+v", result.Blockers)
	}
}

func TestEvaluateIntentHypothesisReadinessWarnsOnContradiction(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	need := &SeoUserNeed{ID: 20, ObjectiveID: 10, Status: UserNeedStatusActive}
	intent := validIntentHypothesisFixture()
	evidence := queryEvidenceFixture(now)
	result := EvaluateIntentHypothesisReadiness(
		need,
		intent,
		[]LinkedIntentQueryEvidence{
			{Link: SeoIntentEvidenceLink{RelationType: IntentEvidenceSupports}, Evidence: evidence},
			{Link: SeoIntentEvidenceLink{RelationType: IntentEvidenceContradicts}, Evidence: evidence},
		},
		now,
	)
	if !hasIntentIssue(result.Warnings, "contradicting_query_evidence") {
		t.Fatalf("missing contradiction warning: %+v", result.Warnings)
	}
}

func TestValidateIntentHypothesisStatusTransition(t *testing.T) {
	cases := []struct {
		current IntentHypothesisStatus
		target  IntentHypothesisStatus
		valid   bool
	}{
		{IntentHypothesisStatusDraft, IntentHypothesisStatusActive, true},
		{IntentHypothesisStatusActive, IntentHypothesisStatusPaused, true},
		{IntentHypothesisStatusPaused, IntentHypothesisStatusActive, true},
		{IntentHypothesisStatusArchived, IntentHypothesisStatusActive, false},
		{IntentHypothesisStatusDraft, IntentHypothesisStatusPaused, false},
	}
	for _, tc := range cases {
		err := ValidateIntentHypothesisStatusTransition(tc.current, tc.target)
		if (err == nil) != tc.valid {
			t.Fatalf("transition %s -> %s valid=%v err=%v", tc.current, tc.target, tc.valid, err)
		}
	}
}

func queryEvidenceFixture(now time.Time) SeoQueryEvidence {
	return SeoQueryEvidence{
		ID:              40,
		ObjectiveID:     10,
		UserNeedID:      20,
		RawQuery:        "tra cứu quy hoạch thửa đất",
		NormalizedQuery: "tra cứu quy hoạch thửa đất",
		QueryHash:       BuildQueryEvidenceHash("tra cứu quy hoạch thửa đất", "vi-VN"),
		SourceSystem:    "search_console",
		SourceKind:      QueryEvidenceSourceFirstParty,
		ObservationType: QueryObservationObserved,
		SourceRecordKey: "fixture-40",
		Locale:          "vi-VN",
		CountryCode:     "VN",
		DeviceCategory:  QueryDeviceAll,
		Confidence:      ObjectiveConfidenceValidated,
		Status:          QueryEvidenceStatusActive,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func hasIntentIssue(items []IntentHypothesisReadinessIssue, code string) bool {
	for _, item := range items {
		if item.Code == code {
			return true
		}
	}
	return false
}
