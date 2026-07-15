package seo_strategy

import (
	"testing"
	"time"
)

func TestCreateSeoIntentHypothesisDraft(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	intent, err := CreateSeoIntentHypothesisDraft(
		10,
		20,
		" Verify Parcel Planning ",
		" Xác minh quy hoạch thửa đất ",
		now,
	)
	if err != nil {
		t.Fatal(err)
	}
	if intent.Code != "verify_parcel_planning" {
		t.Fatalf("unexpected normalized code: %s", intent.Code)
	}
	if intent.Status != IntentHypothesisStatusDraft || intent.LockVersion != 1 {
		t.Fatalf("unexpected defaults: %+v", intent)
	}
}

func TestValidateIntentHypothesisDefinitionRejectsInvalidEnum(t *testing.T) {
	intent := validIntentHypothesisFixture()
	intent.SearchAction = IntentSearchAction("invalid")
	if err := ValidateIntentHypothesisDefinition(intent); err == nil {
		t.Fatal("expected invalid search action error")
	}
}

func TestValidateIntentEvidenceLink(t *testing.T) {
	link := &SeoIntentEvidenceLink{
		ObjectiveID:        1,
		UserNeedID:         2,
		IntentHypothesisID: 3,
		QueryEvidenceID:    4,
		RelationType:       IntentEvidenceSupports,
	}
	if err := ValidateIntentEvidenceLink(link); err != nil {
		t.Fatal(err)
	}
	link.RelationType = IntentEvidenceRelation("invalid")
	if err := ValidateIntentEvidenceLink(link); err == nil {
		t.Fatal("expected invalid relation error")
	}
}

func validIntentHypothesisFixture() *SeoIntentHypothesis {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	return &SeoIntentHypothesis{
		ID:              30,
		ObjectiveID:     10,
		UserNeedID:      20,
		Code:            "verify_parcel_planning",
		Name:            "Xác minh quy hoạch thửa đất",
		GoalStatement:   "Xác minh tác động quy hoạch trước quyết định bất động sản",
		PrimaryQuestion: "Thửa đất này có thuộc quy hoạch hay không?",
		ExpectedResult:  "Nhận được kết quả có nguồn, thời điểm và giới hạn dữ liệu",
		SearchAction:    IntentSearchActionVerify,
		SubjectType:     IntentSubjectParcel,
		ScopeType:       IntentScopeSpecificEntity,
		Status:          IntentHypothesisStatusDraft,
		Confidence:      ObjectiveConfidenceExploratory,
		OwnerCode:       "seo_strategy",
		LockVersion:     1,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
