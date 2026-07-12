package seo_strategy

import (
	"testing"
	"time"
)

func TestObjectiveReadinessRequiresCoreDefinition(
	t *testing.T,
) {
	objective, err := CreateSeoObjectiveDraft(
		"public_search_value",
		"Public Search Value",
		time.Now(),
	)
	if err != nil {
		t.Fatal(err)
	}

	result := EvaluateObjectiveReadiness(
		objective,
		nil,
		nil,
		time.Now(),
	)

	if result.Ready {
		t.Fatal("objective rỗng không được ready")
	}

	if len(result.Blockers) == 0 {
		t.Fatal("phải có blockers")
	}
}

func TestObjectiveReadinessAllowsExploratoryObjective(
	t *testing.T,
) {
	objective, err := CreateSeoObjectiveDraft(
		"public_search_value",
		"Public Search Value",
		time.Now(),
	)
	if err != nil {
		t.Fatal(err)
	}

	objective.AudienceSummary = "Người dùng đang tra cứu"
	objective.Situation = "Trước một quyết định"
	objective.ProblemStatement = "Thông tin phân tán"
	objective.UserOutcome = "Hoàn thành tra cứu"
	objective.ProductOutcome = "Sử dụng capability tra cứu"
	objective.BusinessOutcome = "Tăng qualified usage"
	objective.OwnerCode = "product_search"

	evidence := []SeoObjectiveEvidence{
		{
			ObjectiveID: 1,
			Type:        EvidenceTypeAssumption,
			Stance:      EvidenceStanceContext,
			Finding:     "Nhu cầu hiện mới là giả thuyết",
			Confidence:  ObjectiveConfidenceExploratory,
			Status:      EvidenceStatusActive,
		},
	}

	metrics := []SeoObjectiveMetric{
		{
			ObjectiveID:           1,
			Code:                  "task_completion_rate",
			Name:                  "Task completion rate",
			Layer:                 MetricLayerUser,
			SourceSystem:          "internal_analytics",
			CalculationMethod:     "completed_tasks / eligible_sessions",
			Unit:                  "percent",
			TargetDirection:       TargetDirectionIncrease,
			IsPrimary:             true,
			InstrumentationStatus: InstrumentationStatusPlanned,
			Status:                MetricStatusActive,
		},
	}

	result := EvaluateObjectiveReadiness(
		objective,
		evidence,
		metrics,
		time.Now(),
	)

	if !result.Ready {
		t.Fatalf(
			"objective phải ready, blockers=%v",
			result.Blockers,
		)
	}

	if len(result.Warnings) == 0 {
		t.Fatal(
			"objective exploratory phải còn warnings",
		)
	}
}

func TestValidatedObjectiveRequiresValidatedEvidence(
	t *testing.T,
) {
	objective, _ := CreateSeoObjectiveDraft(
		"public_search_value",
		"Public Search Value",
		time.Now(),
	)

	objective.AudienceSummary = "Audience"
	objective.Situation = "Situation"
	objective.ProblemStatement = "Problem"
	objective.UserOutcome = "User outcome"
	objective.ProductOutcome = "Product outcome"
	objective.BusinessOutcome = "Business outcome"
	objective.OwnerCode = "product_search"
	objective.Confidence = ObjectiveConfidenceValidated

	evidence := []SeoObjectiveEvidence{
		{
			ObjectiveID: 1,
			Type:        EvidenceTypeAssumption,
			Stance:      EvidenceStanceContext,
			Finding:     "Assumption",
			Confidence:  ObjectiveConfidenceExploratory,
			Status:      EvidenceStatusActive,
		},
	}

	metrics := []SeoObjectiveMetric{
		{
			ObjectiveID:           1,
			Code:                  "task_completion_rate",
			Name:                  "Task completion",
			Layer:                 MetricLayerUser,
			SourceSystem:          "internal_analytics",
			CalculationMethod:     "a / b",
			Unit:                  "percent",
			TargetDirection:       TargetDirectionIncrease,
			IsPrimary:             true,
			InstrumentationStatus: InstrumentationStatusInstrumented,
			Status:                MetricStatusActive,
		},
	}

	result := EvaluateObjectiveReadiness(
		objective,
		evidence,
		metrics,
		time.Now(),
	)

	if result.Ready {
		t.Fatal(
			"validated objective không được dựa trên assumption",
		)
	}
}
