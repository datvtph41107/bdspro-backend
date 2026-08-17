package usecase

import (
	"testing"
	"time"

	"crm/internal/dto"
)

func administrativeUnitProjectionFixture(updatedAt *time.Time) *dto.TqdAdministrativeUnitProjection {
	return &dto.TqdAdministrativeUnitProjection{
		SchemaVersion:      dto.TqdAdministrativeUnitProjectionSchemaVersion,
		EntityID:           "ward:42",
		CanonicalKey:       "administrative_unit:ward:00123",
		Slug:               "phuong-ngo-quyen-00123",
		CanonicalPath:      "/tin-tuc/dia-ban/phuong-ngo-quyen-00123",
		UnitType:           "ward",
		UnitTypeName:       "Phường",
		Name:               "Phường Ngô Quyền",
		DisplayName:        "Phường Ngô Quyền",
		AdministrativeCode: "00123",
		Hierarchy: []dto.TqdAdministrativeUnitHierarchyItem{
			{EntityID: "province:1", UnitType: "province", Name: "Thành phố Hà Nội", AdministrativeCode: "01", CanonicalPath: "/tin-tuc/dia-ban/thanh-pho-ha-noi-01"},
			{EntityID: "ward:42", UnitType: "ward", Name: "Phường Ngô Quyền", AdministrativeCode: "00123", CanonicalPath: "/tin-tuc/dia-ban/phuong-ngo-quyen-00123"},
		},
		GeometrySummary: dto.TqdPlanningPreview{
			Centroid:     &dto.TqdPlanningPoint{Latitude: 21.03, Longitude: 105.84},
			Completeness: "centroid_only",
			FitMode:      "center",
		},
		PlanningContext: dto.TqdAdministrativePlanningContext{
			Summary:          "Bối cảnh quy hoạch công khai của Phường Ngô Quyền.",
			DataAvailability: "unavailable",
			RelatedProjects: []dto.TqdAdministrativePlanningProject{
				{
					EntityID:   "planning_project:7",
					Slug:       "do-an-quy-hoach-mau",
					Name:       "Đồ án quy hoạch mẫu",
					Summary:    "Đồ án được liên kết với địa bàn.",
					PublicPath: "/tin-tuc/do-an-quy-hoach/do-an-quy-hoach-mau",
					MapPath:    "/ban-do?focus=planning_project%3A7",
				},
			},
		},
		Source: dto.TqdPlanningProjectionSource{
			System:         "tqd-service",
			Dataset:        "ward_v2",
			Authority:      "aggregated",
			AuthorityName:  "Dữ liệu địa giới hành chính được đồng bộ vào TQD",
			AuthorityClass: "government_derived",
			DataQuality:    "partial",
			UpdatedAt:      updatedAt,
			Limitations:    []string{"Chỉ có điểm trung tâm hành chính."},
		},
	}
}

func TestAdministrativeUnitDocumentUsesBackendOwnedEligibility(t *testing.T) {
	now := time.Date(2026, time.July, 24, 8, 0, 0, 0, time.UTC)
	document := buildAdministrativeUnitPublicDocument(administrativeUnitProjectionFixture(&now))
	warnings := make([]string, 0)
	normalizeSeoPublicDocument(document, &warnings)
	document.SEO.Indexable = true
	applySeoAnswerReadiness(document)

	if document.AnswerReadiness == nil || !document.AnswerReadiness.Ready {
		t.Fatalf("administrative unit answer should be ready: %#v", document.AnswerReadiness)
	}
	if document.IndexDecision == nil || !document.IndexDecision.Eligible {
		t.Fatalf("administrative unit should be index eligible: %#v", document.IndexDecision)
	}
	if document.IndexDecision.PolicyVersion != seoIndexPolicyVersion {
		t.Fatalf("unexpected policy version: %q", document.IndexDecision.PolicyVersion)
	}
	if document.Media == nil || document.Media.Completeness != "centroid_only" || document.Media.FitMode != "center" {
		t.Fatalf("centroid-only geometry contract was not preserved: %#v", document.Media)
	}
	if value := seoFactValue(document.Facts, "unit-type-code"); value != "ward" {
		t.Fatalf("machine-readable unit type is missing: %q", value)
	}
	if document.MapTarget == nil || document.MapTarget.Path != "/ban-do?focus=administrative_unit%3Award%3A42" {
		t.Fatalf("administrative unit map target must use canonical focus contract: %#v", document.MapTarget)
	}
	if len(document.RelatedEntities) != 1 || document.RelatedEntities[0].PublicPath != "/tin-tuc/do-an-quy-hoach/do-an-quy-hoach-mau" {
		t.Fatalf("related planning project handoff is missing: %#v", document.RelatedEntities)
	}
}

func TestAdministrativeUnitDocumentWithoutUpdatedAtFailsReadiness(t *testing.T) {
	document := buildAdministrativeUnitPublicDocument(administrativeUnitProjectionFixture(nil))
	warnings := make([]string, 0)
	normalizeSeoPublicDocument(document, &warnings)
	document.SEO.Indexable = true
	applySeoAnswerReadiness(document)

	if document.AnswerReadiness == nil || document.AnswerReadiness.Ready {
		t.Fatalf("missing updated-at must fail readiness: %#v", document.AnswerReadiness)
	}
	if document.IndexDecision == nil || document.IndexDecision.Eligible {
		t.Fatalf("missing updated-at must fail index eligibility: %#v", document.IndexDecision)
	}
	if !containsString(document.AnswerReadiness.MissingRequiredFacts, "updated-at") {
		t.Fatalf("updated-at missing reason not reported: %#v", document.AnswerReadiness.MissingRequiredFacts)
	}
}

func TestUnknownSeoModuleFailsClosed(t *testing.T) {
	if got := normalizeSeoModuleID("administrative-unit-typo"); got != "" {
		t.Fatalf("unknown module must not fall back to another page family: %q", got)
	}
	if got := resourceTypeForModule("administrative-unit-typo"); got != "" {
		t.Fatalf("unknown resource type must fail closed: %q", got)
	}
	if _, err := publicCanonicalPath("administrative-unit-typo", "example"); err == nil {
		t.Fatal("unknown page family must not receive a public canonical path")
	}
}

func TestNormalizePublicSeoSlugFailsClosed(t *testing.T) {
	invalid := []string{"", "undefined", "null", "[object Object]", "ward/42", `ward\42`, "../ward", "ward?x=1", "ward#section"}
	for _, value := range invalid {
		if normalized, err := normalizePublicSeoSlug(value); err == nil {
			t.Fatalf("unsafe slug %q was accepted as %q", value, normalized)
		}
	}
	if normalized, err := normalizePublicSeoSlug("phuong-ngo-quyen-00123"); err != nil || normalized != "phuong-ngo-quyen-00123" {
		t.Fatalf("valid canonical slug was rejected: slug=%q err=%v", normalized, err)
	}
}

func seoFactValue(facts []dto.SeoFact, id string) string {
	for _, fact := range facts {
		if fact.ID == id {
			return fact.Value
		}
	}
	return ""
}
