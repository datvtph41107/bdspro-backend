package usecase

import (
	"testing"
	"time"

	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
)

func TestBuildSeoGovernanceDocumentPreservesBackendDecisions(t *testing.T) {
	now := time.Date(2026, time.July, 24, 9, 30, 0, 0, time.UTC)
	page := &seo_domain.SeoDomain{
		ID:              42,
		Slug:            "phuong-ben-nghe-26734",
		CanonicalURL:    "/tin-tuc/dia-ban/phuong-ben-nghe-26734",
		PageStatus:      seo_domain.SeoPageStatusPublished,
		Published:       true,
		Scope:           seo_domain.SeoScopePublic,
		IsIndex:         true,
		IsSiteMap:       true,
		RenderStatus:    seo_domain.SeoRenderStatusSuccess,
		NeedGenerate:    false,
		TemplateVersion: dto.SeoPublicDocumentSchemaVersion,
		StaticHtmlHash:  "rendered-hash",
		UpdatedAt:       now,
		GeneratedAt:     &now,
	}
	publicDocument := &dto.SeoPublicDocument{
		SchemaVersion: dto.SeoPublicDocumentSchemaVersion,
		ModuleID:      "administrative-unit",
		ResourceType:  "administrative_area",
		Identity: dto.SeoPublicIdentity{
			ID:            "administrative_unit:ward:26734",
			Slug:          page.Slug,
			CanonicalPath: page.CanonicalURL,
		},
		Facts: []dto.SeoFact{
			{ID: "unit-type", Value: "Phường"},
			{ID: "administrative-hierarchy", Value: "Thành phố Hồ Chí Minh"},
			{ID: "updated-at", Value: now.Format(time.RFC3339)},
		},
		Sections: []dto.SeoSection{
			{ID: "overview", Blocks: []dto.SeoContentBlock{{Type: "paragraph", Text: "Thông tin tổng quan đầy đủ cho đơn vị hành chính này."}}},
			{ID: "planning-context", Blocks: []dto.SeoContentBlock{{Type: "paragraph", Text: "Bối cảnh quy hoạch được tổng hợp từ dữ liệu nguồn hiện hành."}}},
			{ID: "source", Blocks: []dto.SeoContentBlock{{Type: "paragraph", Text: "Nguồn dữ liệu TQD được cập nhật và công bố rõ ràng."}}},
		},
		AnswerReadiness: &dto.SeoAnswerReadiness{
			SpecificationVersion: "seo-m4/administrative-unit/v1",
			RequiredFacts:        []string{"unit-type", "administrative-hierarchy", "updated-at"},
			RequiredSections:     []string{"overview", "planning-context", "source"},
			RenderedSections:     []string{"overview", "planning-context", "source"},
			Ready:                true,
		},
		IndexDecision: &dto.SeoIndexDecision{
			PolicyVersion: "seo-m11/index-eligibility/v1",
			Eligible:      false,
			ReasonCodes:   []string{"non_production_source"},
		},
		Source: dto.SeoPublicSource{
			Environment: "development",
			Authority:   "government_derived",
			SourceName:  "TQD",
			Quality:     "complete",
			UpdatedAt:   now.Format(time.RFC3339),
		},
	}

	governance := buildSeoGovernanceDocument(page, publicDocument, "seo_domain.generated_at", now, nil)
	if governance.AnswerReadiness == nil || !governance.AnswerReadiness.Ready {
		t.Fatalf("governance must preserve backend answer readiness: %#v", governance.AnswerReadiness)
	}
	if governance.IndexDecision == nil || governance.IndexDecision.Eligible {
		t.Fatalf("governance must preserve backend noindex decision: %#v", governance.IndexDecision)
	}
	if governance.AnswerSpecificationVersion != "seo-m4/administrative-unit/v1" {
		t.Fatalf("unexpected answer specification version: %s", governance.AnswerSpecificationVersion)
	}
	for _, item := range governance.RequiredFacts {
		if !item.Present {
			t.Fatalf("required fact unexpectedly absent: %#v", item)
		}
	}
	for _, item := range governance.RequiredSections {
		if !item.Present {
			t.Fatalf("required section unexpectedly absent: %#v", item)
		}
	}
	if governance.Review.AuditAvailable {
		t.Fatal("audit must not be reported available before an audit repository exists")
	}
}
