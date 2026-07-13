package seo_strategy

import (
	"testing"
	"time"
)

func createReadyUserNeedFixture(t *testing.T) (*SeoObjective, *SeoAudienceSegment, *SeoUserNeed) {
	t.Helper()
	now := time.Now()
	objective, err := CreateSeoObjectiveDraft("public_search_value", "Public Search Value", now)
	if err != nil {
		t.Fatal(err)
	}
	objective.ID = 1
	objective.Status = ObjectiveStatusActive

	audience, err := CreateSeoAudienceSegmentDraft(
		"individual_property_decision_maker",
		"Người tự ra quyết định bất động sản",
		now,
	)
	if err != nil {
		t.Fatal(err)
	}
	audience.ID = 1
	audience.Status = AudienceStatusActive

	need, err := CreateSeoUserNeedDraft(
		1,
		1,
		"assess_planning_impact_before_property_decision",
		"Đánh giá tác động quy hoạch",
		now,
	)
	if err != nil {
		t.Fatal(err)
	}
	need.DecisionStage = UserNeedStageVerification
	need.Situation = "Đang cân nhắc một bất động sản"
	need.TriggerStatement = "Nhận thông tin thửa đất"
	need.JobToBeDone = "Hiểu tác động quy hoạch"
	need.BarrierStatement = "Dữ liệu phân tán"
	need.CurrentAlternative = "Tìm nhiều website"
	need.DesiredOutcome = "Hiểu phạm vi và nguồn"
	need.FailureImpact = "Rủi ro quyết định sai"
	need.RiskLevel = UserNeedRiskHigh
	need.Priority = UserNeedPriorityHigh
	need.OwnerCode = "product_search"
	return objective, audience, need
}

func TestUserNeedRequiresActiveObjectiveAndAudience(t *testing.T) {
	now := time.Now()
	objective, _ := CreateSeoObjectiveDraft("public_search_value", "Public Search Value", now)
	audience, _ := CreateSeoAudienceSegmentDraft("individual_property_decision_maker", "Audience", now)
	need, _ := CreateSeoUserNeedDraft(1, 1, "assess_planning_impact", "Assess planning", now)

	result := EvaluateUserNeedReadiness(objective, audience, need, nil, now)
	if result.Ready {
		t.Fatal("objective/audience chưa active không được ready")
	}
}

func TestExploratoryUserNeedCanBeReadyWithWarnings(t *testing.T) {
	objective, audience, need := createReadyUserNeedFixture(t)
	evidence := []LinkedUserNeedEvidence{{
		Link: SeoUserNeedEvidenceLink{RelationType: UserNeedEvidenceContext},
		Evidence: SeoObjectiveEvidence{
			Type:       EvidenceTypeAssumption,
			Status:     EvidenceStatusActive,
			Confidence: ObjectiveConfidenceExploratory,
		},
	}}

	result := EvaluateUserNeedReadiness(objective, audience, need, evidence, time.Now())
	if !result.Ready {
		t.Fatalf("exploratory need phải ready: %+v", result.Blockers)
	}
	if len(result.Warnings) == 0 {
		t.Fatal("assumption/high risk phải có warning")
	}
}

func TestValidatedUserNeedRequiresValidatedSupport(t *testing.T) {
	objective, audience, need := createReadyUserNeedFixture(t)
	need.Confidence = ObjectiveConfidenceValidated
	evidence := []LinkedUserNeedEvidence{{
		Link: SeoUserNeedEvidenceLink{RelationType: UserNeedEvidenceSupports},
		Evidence: SeoObjectiveEvidence{
			Type:       EvidenceTypeAssumption,
			Status:     EvidenceStatusActive,
			Confidence: ObjectiveConfidenceExploratory,
		},
	}}

	result := EvaluateUserNeedReadiness(objective, audience, need, evidence, time.Now())
	if result.Ready {
		t.Fatal("validated need không được dựa trên assumption")
	}
}
