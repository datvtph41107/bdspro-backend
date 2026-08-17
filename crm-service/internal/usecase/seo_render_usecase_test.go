package usecase

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/enums"

	"gorm.io/datatypes"
)

func planningProjectMetadata(t *testing.T) datatypes.JSON {
	t.Helper()
	payload := map[string]any{
		"schemaVersion": "qhpro-public-seo-document/v1",
		"moduleId":      "planning-project",
		"resourceType":  "planning_project",
		"heading": map[string]any{
			"eyebrow": "Đồ án quy hoạch",
		},
		"breadcrumbs": []map[string]any{
			{"label": "Trang chủ", "path": "/"},
			{"label": "Đồ án quy hoạch", "path": "/tin-tuc/do-an-quy-hoach"},
			{"label": "Quy hoạch phân khu đô thị H2-1", "path": "/tin-tuc/do-an-quy-hoach/qh2"},
		},
		"facts": []map[string]any{
			{"id": "planning-type", "label": "Loại quy hoạch", "value": "Quy hoạch phân khu đô thị"},
			{"id": "status", "label": "Trạng thái", "value": "Đã phê duyệt"},
			{"id": "administrative-unit", "label": "Địa bàn", "value": "Khu vực H2-1, Hà Nội"},
			{"id": "updated-at", "label": "Cập nhật", "value": "2026-07-23"},
		},
		"sections": []map[string]any{
			{
				"id": "overview", "title": "Tổng quan",
				"blocks": []map[string]any{{"type": "paragraph", "text": "Đồ án H2-1 cung cấp thông tin tổng quan về định hướng phát triển không gian đô thị và hạ tầng."}},
			},
			{
				"id": "research-scope", "title": "Phạm vi nghiên cứu",
				"blocks": []map[string]any{{"type": "list", "items": []string{"Phạm vi ranh giới nghiên cứu.", "Định hướng sử dụng đất và hạ tầng kỹ thuật."}}},
			},
			{
				"id": "legal-timeline", "title": "Lịch sử pháp lý",
				"blocks": []map[string]any{{"type": "timeline", "items": []map[string]any{{"id": "approval", "title": "Phê duyệt đồ án", "date": "2024-01-15", "description": "Mốc pháp lý được công bố trong bộ dữ liệu nguồn."}}}},
			},
			{
				"id": "source", "title": "Nguồn dữ liệu",
				"blocks": []map[string]any{{"type": "notice", "tone": "info", "title": "Nguồn tham khảo", "text": "Nội dung được tổng hợp từ dữ liệu quy hoạch đã được đồng bộ và cần đối chiếu hồ sơ chính thức."}},
			},
		},
		"faq": map[string]any{
			"groupKey": "planning-project-detail",
			"items": []map[string]any{
				{"id": "scope", "question": "Đồ án H2-1 cung cấp thông tin gì?", "answer": "Trang tổng hợp phạm vi nghiên cứu, định hướng sử dụng đất, hạ tầng và các mốc pháp lý liên quan."},
				{"id": "legal", "question": "Có thể dùng trang này làm căn cứ pháp lý không?", "answer": "Không. Người dùng cần đối chiếu quyết định, bản vẽ và hồ sơ tại cơ quan có thẩm quyền."},
			},
		},
		"mapTarget":    map[string]any{"path": "/ban-do?planningId=68", "label": "Xem đồ án trên bản đồ"},
		"relatedLinks": []map[string]any{{"id": "map", "label": "Tra cứu bản đồ quy hoạch", "path": "/ban-do"}},
		"source": map[string]any{
			"environment": "production",
			"authority":   "official",
			"sourceName":  "Cơ sở dữ liệu quy hoạch QHPro",
			"sourceUrl":   "https://qhpro.vn/quy-hoach/tai-lieu",
			"quality":     "complete",
			"updatedAt":   "2026-07-23T09:38:48+07:00",
			"limitations": []string{"Thông tin có giá trị tham khảo; hồ sơ pháp lý gốc là căn cứ cuối cùng."},
		},
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal metadata: %v", err)
	}
	return datatypes.JSON(raw)
}

func planningProjectPage(t *testing.T) *seo_domain.SeoDomain {
	t.Helper()
	refID := uint64(68)
	publishedAt := time.Date(2026, 7, 23, 2, 38, 48, 0, time.UTC)
	return &seo_domain.SeoDomain{
		ID:              68,
		Slug:            "/tin-tuc/do-an-quy-hoach/qh2",
		CanonicalURL:    "https://qhpro.vn/tin-tuc/do-an-quy-hoach/qh2",
		RefID:           &refID,
		RefSource:       "tqd",
		RefLabel:        "Quy hoạch phân khu đô thị H2-1",
		SourceStatus:    seo_domain.SeoSourceStatusLinked,
		QualityScore:    96,
		Scope:           seo_domain.SeoScopePublic,
		Classify:        enums.ESeoClassifyA,
		PageStatus:      seo_domain.SeoPageStatusPublished,
		Published:       true,
		PublishedAt:     &publishedAt,
		Title:           "Quy hoạch phân khu đô thị H2-1",
		Description:     "Thông tin quy hoạch phân khu đô thị H2-1.",
		Summary:         "Trang SEO kiểm thử cho đồ án quy hoạch H2-1.",
		Content:         "<p>Legacy content must become safe plain text.</p>",
		IsIndex:         true,
		IsSiteMap:       true,
		TemplateKey:     "planning-project",
		TemplateVersion: "qhpro-public-seo-document/v1",
		Metadata:        planningProjectMetadata(t),
		UpdatedAt:       publishedAt,
	}
}

func planningProjectPageWithoutLegalEvents(t *testing.T) *seo_domain.SeoDomain {
	t.Helper()
	page := planningProjectPage(t)

	var metadata map[string]any
	if err := json.Unmarshal(page.Metadata, &metadata); err != nil {
		t.Fatalf("unmarshal planning-project metadata: %v", err)
	}

	sections, ok := metadata["sections"].([]any)
	if !ok {
		t.Fatal("planning-project sections are missing")
	}

	for _, rawSection := range sections {
		section, ok := rawSection.(map[string]any)
		if !ok || section["id"] != "legal-timeline" {
			continue
		}

		section["blocks"] = []any{
			map[string]any{
				"type":  "timeline",
				"items": []any{},
			},
		}
	}

	source, ok := metadata["source"].(map[string]any)
	if !ok {
		t.Fatal("planning-project source is missing")
	}
	source["authority"] = "aggregated"
	source["quality"] = "partial"
	source["limitations"] = []string{
		"Chưa có lịch sử pháp lý được công bố.",
	}

	raw, err := json.Marshal(metadata)
	if err != nil {
		t.Fatalf("marshal partial planning-project metadata: %v", err)
	}
	page.Metadata = datatypes.JSON(raw)

	return page
}

func TestComposePlanningProjectDocumentReady(t *testing.T) {
	document, warnings, err := composeSeoPublicDocument(planningProjectPage(t), time.Now().UTC())
	if err != nil {
		t.Fatalf("compose failed: %v", err)
	}
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %#v", warnings)
	}
	if document.AnswerReadiness == nil || !document.AnswerReadiness.Ready {
		t.Fatalf("expected answer-ready document: %#v", document.AnswerReadiness)
	}
	if !document.SEO.Indexable {
		t.Fatal("expected trusted answer-ready document to remain indexable")
	}
	if document.Identity.CanonicalPath != "/tin-tuc/do-an-quy-hoach/qh2" {
		t.Fatalf("unexpected canonical path: %s", document.Identity.CanonicalPath)
	}
	if len(document.Sections) != 4 {
		t.Fatalf("expected 4 rich sections, got %d", len(document.Sections))
	}
}

func TestComposeWithoutRequiredMetadataBecomesNoindex(t *testing.T) {
	page := planningProjectPage(t)
	page.Metadata = datatypes.JSON([]byte("{}"))
	document, _, err := composeSeoPublicDocument(page, time.Now().UTC())
	if err != nil {
		t.Fatalf("compose failed: %v", err)
	}
	if document.AnswerReadiness == nil || document.AnswerReadiness.Ready {
		t.Fatalf("expected document not ready: %#v", document.AnswerReadiness)
	}
	if document.SEO.Indexable {
		t.Fatal("page must be noindex until required facts and sections exist")
	}
	missing := strings.Join(document.AnswerReadiness.MissingRequiredFacts, ",")
	if !strings.Contains(missing, "administrative-unit") {
		t.Fatalf("expected administrative-unit missing, got %s", missing)
	}
}

func TestComposePlanningProjectAddsLegalTimelineEmptyState(t *testing.T) {
	document, warnings, err := composeSeoPublicDocument(
		planningProjectPageWithoutLegalEvents(t),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("compose failed: %v", err)
	}
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %#v", warnings)
	}

	var legalTimeline *struct {
		Title  string
		Blocks int
		Type   string
		Text   string
	}
	for _, section := range document.Sections {
		if section.ID != "legal-timeline" {
			continue
		}
		value := &struct {
			Title  string
			Blocks int
			Type   string
			Text   string
		}{Title: section.Title, Blocks: len(section.Blocks)}
		if len(section.Blocks) > 0 {
			value.Type = section.Blocks[0].Type
			value.Text = section.Blocks[0].Text
		}
		legalTimeline = value
		break
	}

	if legalTimeline == nil {
		t.Fatal("legal-timeline section is missing")
	}
	if legalTimeline.Blocks != 1 || legalTimeline.Type != "notice" {
		t.Fatalf("expected one notice empty-state block, got %#v", legalTimeline)
	}
	if !strings.Contains(legalTimeline.Text, "chưa có các mốc phê duyệt") {
		t.Fatalf("unexpected legal timeline empty state: %q", legalTimeline.Text)
	}
	if document.AnswerReadiness == nil || document.AnswerReadiness.Ready {
		t.Fatalf("empty-state timeline must not make document answer-ready: %#v", document.AnswerReadiness)
	}
	if !strings.Contains(
		strings.Join(document.AnswerReadiness.RequiredSections, ","),
		"legal-timeline",
	) {
		t.Fatalf("legal-timeline must remain required: %#v", document.AnswerReadiness)
	}
	if !strings.Contains(
		strings.Join(document.AnswerReadiness.RenderedSections, ","),
		"legal-timeline",
	) {
		t.Fatalf("legal-timeline must be visibly rendered: %#v", document.AnswerReadiness)
	}
	if document.SEO.Indexable {
		t.Fatal("partial source with legal timeline empty state must remain noindex")
	}
}

func TestRenderPartialPlanningProjectIncludesLegalTimelineEmptyState(t *testing.T) {
	page := planningProjectPageWithoutLegalEvents(t)
	htmlBody, err := renderSeoHTML(renderSeoInput{
		Page:            page,
		CanonicalURL:    page.CanonicalURL,
		TemplateKey:     page.TemplateKey,
		ModuleID:        "planning-project",
		Environment:     "crm-rendered",
		TelemetryScript: "/static/qhpro/seo-m8-client.js",
		GeneratedAt:     time.Unix(0, 0).UTC(),
	})
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	checks := []string{
		`<meta name="robots" content="noindex,follow">`,
		`data-qhpro-answer-ready="false"`,
		`data-seo-section-id="legal-timeline"`,
		`Lịch sử pháp lý chưa được công bố`,
		`Nguồn dữ liệu hiện tại chưa có các mốc phê duyệt`,
	}
	for _, check := range checks {
		if !strings.Contains(htmlBody, check) {
			t.Fatalf("missing %q in partial rendered HTML", check)
		}
	}
	if strings.Contains(htmlBody, `"@type":"FAQPage"`) {
		t.Fatal("partial answer-not-ready page must not emit FAQPage JSON-LD")
	}
}

func TestRenderSeoHTMLIncludesRichParityTelemetryAndJSONLD(t *testing.T) {
	page := planningProjectPage(t)
	htmlBody, err := renderSeoHTML(renderSeoInput{
		Page:            page,
		CanonicalURL:    page.CanonicalURL,
		TemplateKey:     page.TemplateKey,
		ModuleID:        "planning-project",
		Environment:     "crm-rendered",
		TelemetryScript: "/static/qhpro/seo-m8-client.js",
		GeneratedAt:     time.Unix(0, 0).UTC(),
	})
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	checks := []string{
		`<meta name="robots" content="index,follow">`,
		`data-qhpro-seo-module="planning-project"`,
		`data-qhpro-answer-ready="true"`,
		`src="/static/qhpro/seo-m8-client.js"`,
		`aria-label="Breadcrumb"`,
		`data-seo-fact-id="administrative-unit"`,
		`data-seo-section-id="legal-timeline"`,
		`data-seo-faq-id="scope"`,
		`data-seo-measurement-action="seo_map_cta_click"`,
		`data-seo-measurement-action="seo_related_link_click"`,
		`"@type":"BreadcrumbList"`,
		`"@type":"FAQPage"`,
		`"@type":"Project"`,
	}
	for _, check := range checks {
		if !strings.Contains(htmlBody, check) {
			t.Fatalf("missing %q in rendered HTML", check)
		}
	}
}

func TestRichRendererEscapesDocumentContentAndScriptBreakout(t *testing.T) {
	page := planningProjectPage(t)
	metadata := map[string]any{}
	if err := json.Unmarshal(page.Metadata, &metadata); err != nil {
		t.Fatal(err)
	}
	faq := metadata["faq"].(map[string]any)
	items := faq["items"].([]any)
	items[0].(map[string]any)["answer"] = `Không </script><script>alert("x")</script>`
	raw, err := json.Marshal(metadata)
	if err != nil {
		t.Fatal(err)
	}
	page.Metadata = datatypes.JSON(raw)

	htmlBody, err := renderSeoHTML(renderSeoInput{Page: page, CanonicalURL: page.CanonicalURL, TelemetryScript: "/static/qhpro/seo-m8-client.js", GeneratedAt: time.Now().UTC()})
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	if strings.Contains(htmlBody, `<script>alert`) {
		t.Fatal("script breakout was not escaped")
	}
	if !strings.Contains(htmlBody, `&lt;/script&gt;&lt;script&gt;alert`) {
		t.Fatal("visible FAQ answer must be HTML escaped")
	}
	if !strings.Contains(htmlBody, `\u003c/script\u003e`) {
		t.Fatal("JSON-LD must use script-safe JSON escaping")
	}
}

func TestBuildMetadataScriptProducesValidJSONLD(t *testing.T) {
	got := buildMetadataScript(`{"@context":"https://schema.org","name":"A < B"}`)
	if got == "" {
		t.Fatal("expected JSON-LD script")
	}
	if strings.Contains(got, "&#34;") {
		t.Fatalf("JSON-LD must not be HTML entity encoded: %s", got)
	}
	if !strings.Contains(got, `"@context":"https://schema.org"`) {
		t.Fatalf("missing JSON-LD payload: %s", got)
	}
	if !strings.Contains(got, `A \u003c B`) {
		t.Fatalf("expected script-safe JSON escaping: %s", got)
	}
}

func TestBuildPlanningProjectPublicDocumentPartialKeepsVisibleTimelineButNoindex(t *testing.T) {
	projection := &dto.TqdPlanningProjectProjection{
		ID:                2,
		Slug:              "qh2",
		CanonicalPath:     "/tin-tuc/do-an-quy-hoach/qh2",
		Name:              "Quy hoạch H2-1",
		PlanningTypeName:  "QH sử dụng đất",
		LegalStatusName:   "Dự thảo",
		JurisdictionName:  "UBND Hà Nội",
		PlanningLevelName: "Cấp đô thị",
		Preview: dto.TqdPlanningPreview{
			Bounds:       &dto.TqdPlanningBounds{MinLongitude: 105.1, MinLatitude: 20.1, MaxLongitude: 106.1, MaxLatitude: 21.1},
			Completeness: "complete",
		},
		Source: dto.TqdPlanningProjectionSource{
			Dataset:        "qh_planning_projects",
			AuthorityClass: "aggregated",
			DataQuality:    "partial",
		},
	}

	document := buildPlanningProjectPublicDocument(projection)
	warnings := []string{}
	normalizeSeoPublicDocument(document, &warnings)
	document.SEO.Indexable = true
	applySeoAnswerReadiness(document)

	if document.AnswerReadiness == nil || document.AnswerReadiness.Ready {
		t.Fatalf("partial projection must not be answer-ready: %#v", document.AnswerReadiness)
	}
	if document.SEO.Indexable {
		t.Fatal("partial projection must be noindex")
	}
	found := false
	for _, section := range document.Sections {
		if section.ID != "legal-timeline" {
			continue
		}
		found = true
		if len(section.Blocks) != 1 || section.Blocks[0].Type != "notice" {
			t.Fatalf("expected legal timeline empty-state notice: %#v", section.Blocks)
		}
	}
	if !found {
		t.Fatal("legal-timeline section must remain visible")
	}
}

func TestBuildPlanningProjectPublicDocumentCompleteBecomesReady(t *testing.T) {
	updatedAt := time.Date(2026, 7, 23, 7, 0, 0, 0, time.UTC)
	projection := &dto.TqdPlanningProjectProjection{
		ID:               2,
		Slug:             "qh2",
		CanonicalPath:    "/tin-tuc/do-an-quy-hoach/qh2",
		Name:             "Quy hoạch H2-1",
		PlanningTypeName: "Quy hoạch phân khu",
		LegalStatusName:  "Đã phê duyệt",
		JurisdictionName: "UBND Hà Nội",
		Summary:          "Đồ án định hướng tổ chức không gian, sử dụng đất và hệ thống hạ tầng kỹ thuật trong khu vực nghiên cứu.",
		ResearchScope:    "Phạm vi nghiên cứu bao gồm ranh giới đồ án, các khu chức năng và mạng lưới hạ tầng liên quan.",
		Events:           []dto.TqdPlanningEvent{{ID: 1, Name: "Phê duyệt đồ án", EventDate: &updatedAt}},
		Indicators:       []dto.TqdPlanningIndicator{{Name: "Mật độ xây dựng", Value: "40", Unit: "%"}},
		Preview: dto.TqdPlanningPreview{
			Bounds:       &dto.TqdPlanningBounds{MinLongitude: 105.1, MinLatitude: 20.1, MaxLongitude: 106.1, MaxLatitude: 21.1},
			Completeness: "complete",
		},
		Source: dto.TqdPlanningProjectionSource{
			Dataset:        "qh_planning_projects",
			AuthorityClass: "aggregated",
			DataQuality:    "complete",
			UpdatedAt:      &updatedAt,
		},
	}

	document := buildPlanningProjectPublicDocument(projection)
	warnings := []string{}
	normalizeSeoPublicDocument(document, &warnings)
	document.SEO.Indexable = true
	applySeoAnswerReadiness(document)

	if document.AnswerReadiness == nil || !document.AnswerReadiness.Ready {
		t.Fatalf("complete projection must be answer-ready: %#v", document.AnswerReadiness)
	}
	if !document.SEO.Indexable {
		t.Fatal("complete trusted projection must remain indexable")
	}
}

func TestSeoPublicDocumentVersionIsStableAcrossRenderOutputWrites(t *testing.T) {
	page := planningProjectPage(t)
	before, _, err := composeSeoPublicDocument(page, time.Now().UTC())
	if err != nil {
		t.Fatalf("compose before render failed: %v", err)
	}
	if strings.TrimSpace(before.DocumentVersion) == "" {
		t.Fatal("document version must be backend-owned and non-empty")
	}

	generatedAt := time.Now().UTC().Add(time.Minute)
	page.RenderedHTML = "<html>rendered output</html>"
	page.StaticHtmlHash = "new-render-hash"
	page.GeneratedAt = &generatedAt
	page.UpdatedAt = generatedAt

	after, _, err := composeSeoPublicDocument(page, generatedAt)
	if err != nil {
		t.Fatalf("compose after render failed: %v", err)
	}
	if after.DocumentVersion != before.DocumentVersion {
		t.Fatalf("render output must not change semantic document version: before=%q after=%q", before.DocumentVersion, after.DocumentVersion)
	}
}

func TestSeoPublicDocumentVersionChangesWithSourceContent(t *testing.T) {
	page := planningProjectPage(t)
	before, _, err := composeSeoPublicDocument(page, time.Now().UTC())
	if err != nil {
		t.Fatalf("compose before source update failed: %v", err)
	}

	page.RefHash = "source-version-2"
	after, _, err := composeSeoPublicDocument(page, time.Now().UTC())
	if err != nil {
		t.Fatalf("compose after source update failed: %v", err)
	}
	if after.DocumentVersion == before.DocumentVersion {
		t.Fatalf("source/content changes must update semantic document version: %q", before.DocumentVersion)
	}
}
