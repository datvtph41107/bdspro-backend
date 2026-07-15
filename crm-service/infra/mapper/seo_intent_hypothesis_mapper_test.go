package mapper

import (
	"testing"
	"time"

	seo_strategy "crm/internal/domain/seo_strategy"
)

func TestSeoIntentHypothesisMapperAggregateToPb(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	review := now.Add(24 * time.Hour)
	intent := &seo_strategy.SeoIntentHypothesis{
		ID: 1, ObjectiveID: 10, UserNeedID: 20,
		Code: "verify_parcel_planning", Name: "Verify parcel planning",
		GoalStatement: "Verify planning impact", PrimaryQuestion: "Is this parcel planned?",
		ExpectedResult: "Source-backed answer", SearchAction: seo_strategy.IntentSearchActionVerify,
		SubjectType: seo_strategy.IntentSubjectParcel, ScopeType: seo_strategy.IntentScopeSpecificEntity,
		Status: seo_strategy.IntentHypothesisStatusDraft, Confidence: seo_strategy.ObjectiveConfidenceSupported,
		OwnerCode: "seo", ReviewDueAt: &review, LockVersion: 2, CreatedAt: now, UpdatedAt: now,
	}
	evidence := seo_strategy.SeoQueryEvidence{
		ID: 40, ObjectiveID: 10, UserNeedID: 20,
		RawQuery: "parcel planning", NormalizedQuery: "parcel planning",
		QueryHash:    seo_strategy.BuildQueryEvidenceHash("parcel planning", "vi-VN"),
		SourceSystem: "search_console", SourceKind: seo_strategy.QueryEvidenceSourceFirstParty,
		ObservationType: seo_strategy.QueryObservationObserved, SourceRecordKey: "40",
		Locale: "vi-VN", CountryCode: "VN", DeviceCategory: seo_strategy.QueryDeviceAll,
		Confidence: seo_strategy.ObjectiveConfidenceValidated, Status: seo_strategy.QueryEvidenceStatusActive,
		CreatedAt: now, UpdatedAt: now,
	}
	mapper := NewSeoIntentHypothesisMapper(
		&SeoUserNeedMapper{},
		NewSeoQueryEvidenceMapper(),
	)
	response := mapper.AggregateToPb(&seo_strategy.SeoIntentHypothesisAggregate{
		Intent: intent,
		UserNeed: &seo_strategy.SeoUserNeed{
			ID: 20, ObjectiveID: 10, Code: "need", Name: "Need",
			DecisionStage: seo_strategy.UserNeedStageVerification,
			RiskLevel:     seo_strategy.UserNeedRiskHigh,
			Priority:      seo_strategy.UserNeedPriorityHigh,
			Status:        seo_strategy.UserNeedStatusActive,
			Confidence:    seo_strategy.ObjectiveConfidenceSupported,
			LockVersion:   1, CreatedAt: now, UpdatedAt: now,
		},
		Evidence: []seo_strategy.LinkedIntentQueryEvidence{{
			Link: seo_strategy.SeoIntentEvidenceLink{
				RelationType:  seo_strategy.IntentEvidenceSupports,
				RelevanceNote: "support", CreatedAt: now,
			},
			Evidence: evidence,
		}},
		Readiness: seo_strategy.IntentHypothesisReadiness{Ready: true},
	})
	if response == nil {
		t.Fatal("response must not be nil")
	}
	if response.Id != 1 || response.SearchAction != "verify" || response.LockVersion != 2 {
		t.Fatalf("unexpected response: %+v", response)
	}
	if len(response.Evidence) != 1 || response.Evidence[0].QueryEvidence.GetId() != 40 {
		t.Fatalf("unexpected evidence projection: %+v", response.Evidence)
	}
	if response.Readiness == nil || !response.Readiness.Ready {
		t.Fatalf("unexpected readiness: %+v", response.Readiness)
	}
}
