package mapper

import (
	"testing"

	tqdpb "pb/types/tqd"
	qh "userv4/modules/tqd/entity/qh"
)

func TestProjectDetailFieldsRoundTripThroughAdminContract(t *testing.T) {
	t.Parallel()
	mapper := NewQHPlanningMapper()
	project := &qh.QHPlanningProject{
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

func TestProjectClientViewDefaultsTrueAndRoundTrips(t *testing.T) {
	t.Parallel()
	mapper := NewQHPlanningMapper()

	domain := mapper.ToProjectDomain(&tqdpb.PlanningProject{})
	if !domain.ClientView {
		t.Fatal("clientView must default to true when omitted")
	}

	hidden := false
	domain = mapper.ToProjectDomain(&tqdpb.PlanningProject{ClientView: &hidden})
	if domain.ClientView {
		t.Fatal("explicit clientView=false must be preserved")
	}

	proto := mapper.ToProjectProto(&qh.QHPlanningProject{ClientView: false})
	if proto.GetClientView() {
		t.Fatal("hidden project must serialize clientView=false")
	}
}

func TestProjectIndicatorsDefaultToJSONArray(t *testing.T) {
	t.Parallel()
	domain := NewQHPlanningMapper().ToProjectDomain(&tqdpb.PlanningProject{})
	if domain.Indicators != "[]" {
		t.Fatalf("unexpected indicators default: %q", domain.Indicators)
	}
}

func TestDocumentLegalMetadataRoundTrip(t *testing.T) {
	t.Parallel()
	mapper := NewQHPlanningMapper()
	document := &qh.QHPlanningDocument{
		Indicators:     `[{"name":"Tỷ lệ đất ở","value":"60%"}]`,
		Authority:      "UBND thành phố",
		Sector:         "Xây dựng",
		Field:          "Quy hoạch",
		Signer:         "Nguyễn Văn A",
		ValidityStatus: "Còn hiệu lực",
		RelativePath:   "ban-do/QH.pdf",
		Status:         20,
	}

	proto := mapper.ToDocumentProto(document)
	if proto.GetIndicators() != document.Indicators ||
		proto.GetAuthority() != document.Authority ||
		proto.GetSector() != document.Sector ||
		proto.GetField() != document.Field ||
		proto.GetSigner() != document.Signer ||
		proto.GetValidityStatus() != document.ValidityStatus ||
		proto.GetRelativePath() != document.RelativePath ||
		proto.GetStatus() != document.Status {
		t.Fatalf("unexpected document proto: %#v", proto)
	}

	domain := mapper.ToDocumentDomain(proto)
	if domain.Indicators != document.Indicators ||
		domain.Authority != document.Authority ||
		domain.Sector != document.Sector ||
		domain.Field != document.Field ||
		domain.Signer != document.Signer ||
		domain.ValidityStatus != document.ValidityStatus ||
		domain.RelativePath != document.RelativePath ||
		domain.Status != document.Status {
		t.Fatalf("document legal metadata did not round-trip: %#v", domain)
	}
}

func TestDocumentIssuedDateAliasRoundTrip(t *testing.T) {
	t.Parallel()
	domain := NewQHPlanningMapper().ToDocumentDomain(&tqdpb.PlanningDocument{
		IssuedDate:    "2025-08-01",
		EffectiveDate: "2025-08-05",
	})
	if domain.IssueDate == nil || domain.IssueDate.Format("2006-01-02") != "2025-08-01" {
		t.Fatalf("issuedDate alias was not mapped: %#v", domain.IssueDate)
	}
}

func TestDocumentIndicatorsDefaultToJSONArray(t *testing.T) {
	t.Parallel()
	domain := NewQHPlanningMapper().ToDocumentDomain(&tqdpb.PlanningDocument{})
	if domain.Indicators != "[]" {
		t.Fatalf("unexpected indicators default: %q", domain.Indicators)
	}
}
