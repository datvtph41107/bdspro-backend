package mapper

import (
	"testing"

	tqdpb "pb/types/tqd"
	qh_domain "tqd/internal/domain/qh"
)

func TestProjectDetailFieldsRoundTripThroughAdminContract(t *testing.T) {
	t.Parallel()
	mapper := NewQHPlanningMapper()
	project := &qh_domain.QHPlanningProject{
		ResearchScope: "Phạm vi toàn phường",
		Indicators:    `[{"name":"Mật độ xây dựng","value":"40%"}]`,
	}

	proto := mapper.ToProjectProto(project)
	if proto.GetResearchScope() != project.ResearchScope {
		t.Fatalf("unexpected research scope: %q", proto.GetResearchScope())
	}
	if proto.GetIndicators() != project.Indicators {
		t.Fatalf("unexpected indicators: %q", proto.GetIndicators())
	}

	domain := mapper.ToProjectDomain(&tqdpb.PlanningProject{
		ResearchScope: proto.GetResearchScope(),
		Indicators:    proto.GetIndicators(),
	})
	if domain.ResearchScope != project.ResearchScope || domain.Indicators != project.Indicators {
		t.Fatalf("detail fields did not round-trip: %#v", domain)
	}
}

func TestProjectIndicatorsDefaultToJSONArray(t *testing.T) {
	t.Parallel()
	domain := NewQHPlanningMapper().ToProjectDomain(&tqdpb.PlanningProject{})
	if domain.Indicators != "[]" {
		t.Fatalf("unexpected indicators default: %q", domain.Indicators)
	}
}
