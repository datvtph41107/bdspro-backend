package seo_strategy

import (
	"testing"
	"time"
)

func TestAudienceReadinessRequiresDefinition(t *testing.T) {
	audience, err := CreateSeoAudienceSegmentDraft(
		"individual_property_decision_maker",
		"Người tự ra quyết định bất động sản",
		time.Now(),
	)
	if err != nil {
		t.Fatal(err)
	}

	result := EvaluateAudienceReadiness(audience, time.Now())
	if result.Ready {
		t.Fatal("audience rỗng không được ready")
	}
	if len(result.Blockers) == 0 {
		t.Fatal("phải có blockers")
	}
}

func TestAudienceReadinessAllowsCompleteDraft(t *testing.T) {
	audience, _ := CreateSeoAudienceSegmentDraft(
		"individual_property_decision_maker",
		"Người tự ra quyết định bất động sản",
		time.Now(),
	)
	audience.Description = "Người mua, chủ sở hữu hoặc nhà đầu tư cá nhân"
	audience.DecisionRole = AudienceDecisionRoleSelfDecisionMaker
	audience.KnowledgeLevel = AudienceKnowledgeBeginner
	audience.OwnerCode = "product_search"

	result := EvaluateAudienceReadiness(audience, time.Now())
	if !result.Ready {
		t.Fatalf("audience phải ready: %+v", result.Blockers)
	}
}
