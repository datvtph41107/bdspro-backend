package usecase

import (
	"strings"
	"testing"

	"crm/internal/dto"
)

func TestSeoOperatingSystemV2UAppliesAllSixteenModules(t *testing.T) {
	document := &dto.SeoPublicDocument{
		SchemaVersion:   dto.SeoPublicDocumentSchemaVersion,
		DocumentVersion: "qhpro-public-document/test",
		ModuleID:        "administrative-unit",
		ResourceType:    "administrative_area",
		Identity: dto.SeoPublicIdentity{
			ID:            "ward:42",
			Slug:          "phuong-mau-00042",
			CanonicalPath: "/dia-ban/phuong-mau-00042",
		},
		Heading: dto.SeoPublicHeading{
			Title:   "Phường Mẫu",
			Summary: "Thông tin hành chính và bối cảnh quy hoạch của Phường Mẫu.",
		},
		Facts: []dto.SeoFact{
			{ID: "unit-type", Label: "Loại đơn vị", Value: "Phường"},
			{ID: "administrative-hierarchy", Label: "Phân cấp", Value: "Tỉnh Mẫu / Phường Mẫu"},
			{ID: "updated-at", Label: "Cập nhật", Value: "2026-07-24T00:00:00Z"},
		},
		Sections: []dto.SeoSection{
			{ID: "overview", Title: "Tổng quan", Blocks: []dto.SeoContentBlock{{Type: "paragraph", Text: "Phường Mẫu là đơn vị hành chính được công bố trong nguồn dữ liệu chính thức."}}},
			{ID: "planning-context", Title: "Bối cảnh quy hoạch", Blocks: []dto.SeoContentBlock{{Type: "paragraph", Text: "Bối cảnh quy hoạch được tổng hợp từ dữ liệu quy hoạch đang có hiệu lực."}}},
			{ID: "source", Title: "Nguồn", Blocks: []dto.SeoContentBlock{{Type: "notice", Text: "Nguồn dữ liệu TQD chính thức."}}},
		},
		MapTarget: &dto.SeoMapTarget{Path: "/ban-do?focus=administrative_unit%3Award%3A42", Label: "Xem trên bản đồ"},
		SEO: dto.SeoPublicSEO{
			Title:         "Phường Mẫu | QHPro",
			Description:   "Tra cứu Phường Mẫu.",
			CanonicalPath: "/dia-ban/phuong-mau-00042",
			Indexable:     true,
			Follow:        true,
		},
		Source: dto.SeoPublicSource{
			Environment: "production",
			Authority:   "official",
			SourceName:  "TQD",
			Quality:     "complete",
			UpdatedAt:   "2026-07-24T00:00:00Z",
		},
	}

	applySeoAnswerReadiness(document)

	if document.OperatingSystem == nil || len(document.OperatingSystem.Modules) != 16 {
		t.Fatalf("V2.U must expose all 16 SEO OS modules: %#v", document.OperatingSystem)
	}
	if document.OperatingSystem.Release != dto.SeoOperatingSystemRelease {
		t.Fatalf("unexpected release: %q", document.OperatingSystem.Release)
	}
	if document.ClientHandoff == nil || document.ClientHandoff.ContractVersion != dto.SeoPageBootstrapContractVersion {
		t.Fatalf("client handoff missing: %#v", document.ClientHandoff)
	}
	if document.ClientHandoff.CanonicalSlug != document.Identity.Slug || document.ClientHandoff.EntityID != document.Identity.ID || document.ClientHandoff.DocumentVersion != document.DocumentVersion {
		t.Fatalf("client handoff identity mismatch: %#v", document.ClientHandoff)
	}
	if document.ClientHandoff.MapFocus == nil || document.ClientHandoff.MapFocus.ID != "ward:42" {
		t.Fatalf("map focus missing: %#v", document.ClientHandoff.MapFocus)
	}
	if document.Measurement == nil || len(document.Measurement.RequiredEvents) == 0 {
		t.Fatalf("measurement plan missing: %#v", document.Measurement)
	}
	if document.GenerativeAI == nil || !document.GenerativeAI.Ready || len(document.GenerativeAI.Citations) != 1 {
		t.Fatalf("generative answer must be citation ready: %#v", document.GenerativeAI)
	}
}

func TestSeoOperatingSystemV2UManifestHasSixteenModules(t *testing.T) {
	manifest := buildSeoOperatingSystemManifest()
	if manifest.Release != dto.SeoOperatingSystemRelease {
		t.Fatalf("unexpected release: %q", manifest.Release)
	}
	if len(manifest.Modules) != 16 {
		t.Fatalf("expected 16 modules, got %d", len(manifest.Modules))
	}
	if manifest.Modules[0].Number != 1 || manifest.Modules[15].Number != 16 {
		t.Fatalf("module ordering is not stable: %#v", manifest.Modules)
	}
	for _, module := range manifest.Modules {
		if len(module.Evidence) == 0 {
			t.Fatalf("module %d must expose operating evidence: %#v", module.Number, module)
		}
	}
	if len(manifest.RuntimeFlow) != 10 {
		t.Fatalf("expected 10 runtime stages, got %d", len(manifest.RuntimeFlow))
	}
	if len(manifest.PilotGates) != 7 {
		t.Fatalf("expected B1-B7 pilot gates, got %d", len(manifest.PilotGates))
	}
	for _, gate := range manifest.PilotGates {
		if strings.TrimSpace(gate.CurrentStatus) == "" || len(gate.Evidence) == 0 {
			t.Fatalf("pilot gate %s must expose current status and evidence", gate.ID)
		}
	}
	if manifest.Milestone != "v2u-full-source-implementation" || manifest.CurrentPageFamily != "administrative-unit" || manifest.MaximumDeliveryState != "crm-pilot" {
		t.Fatalf("unexpected V2.U milestone contract: %#v", manifest)
	}
	if len(manifest.PageFamilyCapabilities) != 7 {
		t.Fatalf("expected seven page-family capability records, got %d", len(manifest.PageFamilyCapabilities))
	}
	if manifest.PageFamilyCapabilities[0].ID != "administrative-unit" || manifest.PageFamilyCapabilities[0].ActivationState != "vertical-slice-active" {
		t.Fatalf("Administrative Unit capability is not activated: %#v", manifest.PageFamilyCapabilities[0])
	}
	for _, gate := range manifest.PilotGates {
		if !gate.BlocksPilot {
			t.Fatalf("all Phase B gates must block pilot: %#v", gate)
		}
	}
}

func TestSeoOperatingSystemV2UModuleSnapshotMergesStaticAndRuntimeEvidence(t *testing.T) {
	document := &dto.SeoPublicDocument{
		SchemaVersion: dto.SeoPublicDocumentSchemaVersion,
		ModuleID:      "administrative-unit",
		Identity: dto.SeoPublicIdentity{
			ID:            "ward:42",
			Slug:          "phuong-mau-00042",
			CanonicalPath: "/dia-ban/phuong-mau-00042",
		},
		SEO: dto.SeoPublicSEO{
			Title:         "Phường Mẫu | QHPro",
			Description:   "Tra cứu Phường Mẫu.",
			CanonicalPath: "/dia-ban/phuong-mau-00042",
		},
		Source: dto.SeoPublicSource{Authority: "official", UpdatedAt: "2026-07-24T00:00:00Z"},
	}

	modules := buildSeoOperatingModuleSnapshot(document)
	if len(modules) != 16 {
		t.Fatalf("expected 16 module snapshots, got %d", len(modules))
	}

	moduleByNumber := make(map[int]dto.SeoOperatingModuleCapability, len(modules))
	for _, module := range modules {
		moduleByNumber[module.Number] = module
	}

	m6 := moduleByNumber[6]
	if !containsString(m6.Evidence, "canonical-path") || !containsString(m6.Evidence, "route-registry") {
		t.Fatalf("M06 must merge static and runtime evidence: %#v", m6.Evidence)
	}
	m10 := moduleByNumber[10]
	if !containsString(m10.Evidence, "page-bootstrap-v2u") || !containsString(m10.Evidence, "next-ssr") {
		t.Fatalf("M10 must merge static and runtime evidence: %#v", m10.Evidence)
	}
}

func containsString(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func TestSeoOperatingSystemV2UOnlyActivatesAdministrativeUnitIsland(t *testing.T) {
	administrative := &dto.SeoPublicDocument{
		ModuleID: "administrative-unit",
		Identity: dto.SeoPublicIdentity{ID: "ward:42", Slug: "phuong-mau-00042", CanonicalPath: "/dia-ban/phuong-mau-00042"},
	}
	if bootstrap := buildSeoPageBootstrap(administrative); bootstrap == nil || bootstrap.MountPolicy != "client-island" {
		t.Fatalf("administrative unit must activate the V2.U client island: %#v", bootstrap)
	}

	project := &dto.SeoPublicDocument{
		ModuleID: "planning-project",
		Identity: dto.SeoPublicIdentity{ID: "42", Slug: "do-an-mau", CanonicalPath: "/do-an-quy-hoach/do-an-mau"},
	}
	if bootstrap := buildSeoPageBootstrap(project); bootstrap == nil || bootstrap.MountPolicy != "not-activated" {
		t.Fatalf("future page families must remain framework-only until activated: %#v", bootstrap)
	}
}
