package usecase

import (
	"net/url"
	"strings"

	"crm/internal/dto"
)

type seoOperatingModuleDefinition struct {
	Number      int
	ID          string
	Name        string
	Owner       string
	RuntimeRole string
	State       string
	Evidence    []string
}

var seoOperatingModuleDefinitions = []seoOperatingModuleDefinition{
	{1, "objective-value", "Objective & Value", "crm", "strategy-policy", "domain-implemented", []string{"seo_strategy.objective", "objective-readiness", "review-due-at"}},
	{2, "audience-user-need", "Audience & User Need", "crm", "strategy-policy", "domain-implemented", []string{"seo_strategy.audience", "seo_strategy.user_need", "activation-readiness"}},
	{3, "query-keyword-intent", "Query, Keyword & Intent", "crm", "strategy-policy", "domain-implemented", []string{"query-evidence", "intent-hypothesis", "query-cluster"}},
	{4, "answer-specification", "Answer Specification", "crm", "document-policy", "policy-implemented", []string{"answer-specification-version", "required-facts", "required-sections"}},
	{5, "knowledge-data-trust", "Knowledge, Data & Trust", "crm+tqd", "source-contract", "contract-implemented", []string{"tqd-projection", "source-authority", "source-freshness", "limitations"}},
	{6, "page-strategy-ia", "Page Strategy & Information Architecture", "crm+next", "canonical-routing", "contract-implemented", []string{"canonical-path", "route-registry", "breadcrumbs", "related-links"}},
	{7, "ux-landing-conversion", "UX, Landing Page & Conversion", "next+react-native", "presentation-interaction", "administrative-unit-active", []string{"next-semantic-layout", "QHDetailScreen", "map-focus-action", "related-content-actions"}},
	{8, "public-search-document", "Public Search Document", "crm", "semantic-document", "operational-source", []string{"PublicSeoDocument", "GET+HEAD+ETag+304", "typed-provisioning", "rendered-html"}},
	{9, "discovery-crawling", "Discovery & Crawling", "crm+tqd", "discovery-sitemap", "operational-source", []string{"tqd-discovery", "canonical-detail-url", "crm-sitemap", "index-eligibility"}},
	{10, "rendering-javascript", "Rendering & JavaScript", "next+react-native", "ssr-client-island", "v2u-source-implemented", []string{"next-ssr-html", "page-bootstrap-v2u", "rn-client-island", "no-initial-double-fetch"}},
	{11, "indexing-canonicalization", "Indexing & Canonicalization", "crm", "index-policy", "operational-source", []string{"backend-index-decision", "canonical-slug", "fail-closed", "sitemap-policy"}},
	{12, "search-appearance", "Search Appearance", "crm+next", "metadata-structured-data", "operational-source", []string{"title-description", "robots", "BreadcrumbList", "Place", "FAQPage-conditional"}},
	{13, "relevance-quality-authority", "Relevance, Quality & Authority", "crm", "quality-readiness", "operational-source", []string{"answer-readiness", "quality-score", "source-trust", "freshness-blockers"}},
	{14, "generative-ai-search", "Generative AI Search", "crm", "answer-provenance", "v2u-foundation", []string{"citation-ready-blocks", "source-citations", "fact-section-provenance", "limitations"}},
	{15, "measurement-experimentation", "Measurement & Experimentation", "crm+next+react-native", "telemetry-feedback", "runtime-pipeline-implemented", []string{"measurement-plan-v2u", "next-sanitizer", "crm-runtime-event-store", "web-vitals", "handoff-events"}},
	{16, "governance-production", "Governance & Production", "crm+admin", "governance-operations", "v2u-source-implemented", []string{"operating-manifest", "admin-governance", "typed-provisioning", "internal-key-boundary", "B0-B7-gates"}},
}

func buildSeoOperatingSystemManifest() *dto.SeoOperatingSystemManifest {
	modules := make([]dto.SeoOperatingModuleCapability, 0, len(seoOperatingModuleDefinitions))
	for _, module := range seoOperatingModuleDefinitions {
		modules = append(modules, dto.SeoOperatingModuleCapability{
			Number:      module.Number,
			ID:          module.ID,
			Name:        module.Name,
			Owner:       module.Owner,
			RuntimeRole: module.RuntimeRole,
			State:       module.State,
			Evidence:    append([]string(nil), module.Evidence...),
		})
	}

	return &dto.SeoOperatingSystemManifest{
		SchemaVersion:        dto.SeoOperatingSystemManifestVersion,
		Release:              dto.SeoOperatingSystemRelease,
		Milestone:            "v2u-full-source-implementation",
		ImplementationStatus: "source-implemented-awaiting-verification",
		CurrentPageFamily:    "administrative-unit",
		MaximumDeliveryState: "crm-pilot",
		Modules:              modules,
		PageFamilies: []string{
			"administrative-unit",
			"parcel",
			"planning-region",
			"planning-project",
			"planning-news",
			"planning-report",
			"legal-document",
		},
		PageFamilyCapabilities: []dto.SeoPageFamilyCapability{
			{ID: "administrative-unit", CanonicalRoute: "/tin-tuc/dia-ban/:slug", ActivationState: "vertical-slice-active", DeliveryOwner: "next", DeliveryState: "next-active", BackendContract: "projection+typed-provisioning+public-document+render", NextRendering: "ssr-semantic-html+rn-client-island", ReactNativeSurface: "QHDetailScreen-active", AdminGovernance: "read+provision+publish+generate", PilotReadiness: "blocked-by-b1-b7-verification"},
			{ID: "planning-project", CanonicalRoute: "/tin-tuc/do-an-quy-hoach/:slug", ActivationState: "existing-delivery-reference", DeliveryOwner: "crm", DeliveryState: "crm-active-existing", BackendContract: "projection+public-document+render", NextRendering: "external-rewrite+fallback", ReactNativeSurface: "detail-contract-not-v2u-activated", AdminGovernance: "generic", PilotReadiness: "outside-current-page-family"},
			{ID: "parcel", CanonicalRoute: "/thua-dat/:slug", ActivationState: "contract-foundation", DeliveryOwner: "next", DeliveryState: "next-active", BackendContract: "legacy+public-document-foundation", NextRendering: "ssr", ReactNativeSurface: "not-v2u-activated", AdminGovernance: "generic", PilotReadiness: "not-started"},
			{ID: "planning-region", CanonicalRoute: "/vung-quy-hoach/:slug", ActivationState: "contract-foundation", DeliveryOwner: "next", DeliveryState: "next-active", BackendContract: "public-document-foundation", NextRendering: "ssr", ReactNativeSurface: "not-v2u-activated", AdminGovernance: "generic", PilotReadiness: "not-started"},
			{ID: "planning-news", CanonicalRoute: "/tin-tuc/bai-viet/:slug", ActivationState: "contract-foundation", DeliveryOwner: "next", DeliveryState: "next-active", BackendContract: "catalog+public-document-foundation", NextRendering: "ssr", ReactNativeSurface: "not-v2u-activated", AdminGovernance: "generic", PilotReadiness: "not-started"},
			{ID: "planning-report", CanonicalRoute: "/tin-tuc/bao-cao/:slug", ActivationState: "contract-foundation", DeliveryOwner: "next", DeliveryState: "next-active", BackendContract: "public-document-foundation", NextRendering: "ssr", ReactNativeSurface: "not-v2u-activated", AdminGovernance: "generic", PilotReadiness: "not-started"},
			{ID: "legal-document", CanonicalRoute: "/van-ban/:slug", ActivationState: "registry-gap", DeliveryOwner: "unresolved", DeliveryState: "not-activated", BackendContract: "discovery-foundation-only", NextRendering: "not-activated", ReactNativeSurface: "not-v2u-activated", AdminGovernance: "generic", PilotReadiness: "blocked-by-owner-contract"},
		},
		RuntimeFlow: []dto.SeoOperatingRuntimeStage{
			{Order: 1, ID: "public-request", Name: "Public request enters Next", Owner: "next:proxy.ts", Input: "browser-or-googlebot-request", Output: "site-and-route-context"},
			{Order: 2, ID: "delivery-resolution", Name: "Resolve page-family delivery owner", Owner: "next:route-registry", Input: "site-and-route-context", Output: "delivery-owner-and-state"},
			{Order: 3, ID: "public-document", Name: "Compose authoritative PublicSeoDocument", Owner: "crm+tqd", Input: "canonical-slug", Output: "semantic-document-v1+v2u"},
			{Order: 4, ID: "server-render", Name: "Render crawlable semantic HTML", Owner: "next", Input: "public-seo-document", Output: "html-metadata-jsonld"},
			{Order: 5, ID: "initial-paint", Name: "Display usable content before JavaScript", Owner: "browser", Input: "server-html", Output: "visible-semantic-page"},
			{Order: 6, ID: "artifact-resolution", Name: "Resolve active RN Web artifact", Owner: "next:public/qhmap/assets.json", Input: "artifact-manifest", Output: "bridge-assets"},
			{Order: 7, ID: "client-island", Name: "Mount RN interactive detail surface", Owner: "next+react-native", Input: "page-bootstrap-v2u", Output: "qh-detail-screen"},
			{Order: 8, ID: "navigation-intent", Name: "Translate RN intent into browser navigation", Owner: "react-native->next", Input: "public-navigation-intent", Output: "next-router-transition"},
			{Order: 9, ID: "measurement", Name: "Collect correlated SEO and interaction evidence", Owner: "next+react-native+crm", Input: "runtime-events", Output: "measurement-evidence"},
			{Order: 10, ID: "governance", Name: "Expose backend decisions and release gates", Owner: "crm+admin", Input: "document-and-runtime-state", Output: "operator-read-model"},
		},
		PilotGates: []dto.SeoOperatingVerificationGate{
			{ID: "b1-backend", Name: "Backend compile, tests, vet and generation", RequiredStatus: "VERIFIED", CurrentStatus: "NOT STARTED", Evidence: []string{"v2u-backend-source-implemented", "go-ast-contract-check"}, BlockingReason: "rerun-go-1.25-tests-vet-protobuf-wire-after-v2u", BlocksPilot: true},
			{ID: "b2-frontends", Name: "Next, RN and Admin production builds", RequiredStatus: "BUILD PASS", CurrentStatus: "NOT STARTED", Evidence: []string{"frontend-source-contracts-implemented", "static-architecture-gates"}, BlockingReason: "install-locked-dependencies-and-run-production-builds", BlocksPilot: true},
			{ID: "b3-runtime", Name: "Database and service runtime smoke", RequiredStatus: "RUNTIME PASS", CurrentStatus: "BLOCKED", Evidence: []string{"runtime-smoke-harness-implemented", "typed-provisioning-implemented"}, BlockingReason: "requires-postgresql-tqd-crm-gateway-next-rn-runtime", BlocksPilot: true},
			{ID: "b4-parity", Name: "Next-versus-CRM semantic parity", RequiredStatus: "VERIFIED", CurrentStatus: "IMPLEMENTED", Evidence: []string{"semantic-parity-harness", "machine-readable-report-contract"}, BlockingReason: "requires-live-crm-and-next-output", BlocksPilot: true},
			{ID: "b5-seo-rn", Name: "SEO to RN page handoff E2E", RequiredStatus: "E2E PASS", CurrentStatus: "IMPLEMENTED", Evidence: []string{"playwright-page-handoff-spec", "rn-client-island-contract", "navigation-loop-guards"}, BlockingReason: "requires-production-build-and-browser-runtime", BlocksPilot: true},
			{ID: "b6-governance", Name: "Admin governance runtime", RequiredStatus: "VERIFIED", CurrentStatus: "STATIC PASS", Evidence: []string{"admin-operating-console", "governance-read-and-provision-operations", "admin-static-contract-gate"}, BlockingReason: "requires-authenticated-gateway-and-live-crm-runtime", BlocksPilot: true},
			{ID: "b7-pilot", Name: "Pilot checklist and rollback proof", RequiredStatus: "VERIFIED", CurrentStatus: "BLOCKED", Evidence: []string{"crm-pilot-no-public-rewrite-regression-test", "maximum-delivery-state-crm-pilot"}, BlockingReason: "blocked-until-b1-through-b6-reach-required-status", BlocksPilot: true},
		},
		Invariants: []string{
			"backend-owns-canonical-readiness-index-decision",
			"next-proxy-is-the-network-boundary",
			"next-route-registry-is-the-delivery-owner-registry",
			"next-owns-browser-url-when-react-native-is-embedded",
			"react-native-emits-navigation-intent-only",
			"crm-owns-qhpro-sitemap",
			"rn-web-assets-json-is-the-active-artifact-source",
			"production-fixtures-are-forbidden",
		},
	}
}

func applySeoOperatingSystemV2U(document *dto.SeoPublicDocument) {
	if document == nil {
		return
	}

	document.ClientHandoff = buildSeoPageBootstrap(document)
	document.Measurement = buildSeoMeasurementPlan(document)
	document.GenerativeAI = buildSeoGenerativeAnswer(document)
	document.OperatingSystem = &dto.SeoOperatingSystemSnapshot{
		Release:       dto.SeoOperatingSystemRelease,
		PageFamily:    document.ModuleID,
		DeliveryOwner: "next-route-registry",
		DeliveryState: "resolved-by-consumer",
		Modules:       buildSeoOperatingModuleSnapshot(document),
	}
}

func buildSeoOperatingModuleSnapshot(document *dto.SeoPublicDocument) []dto.SeoOperatingModuleCapability {
	modules := make([]dto.SeoOperatingModuleCapability, 0, len(seoOperatingModuleDefinitions))
	for _, definition := range seoOperatingModuleDefinitions {
		state, runtimeEvidence := seoOperatingModuleState(definition.Number, document)
		evidence := append([]string(nil), definition.Evidence...)
		evidence = appendUniqueStrings(evidence, runtimeEvidence...)
		modules = append(modules, dto.SeoOperatingModuleCapability{
			Number:      definition.Number,
			ID:          definition.ID,
			Name:        definition.Name,
			Owner:       definition.Owner,
			RuntimeRole: definition.RuntimeRole,
			State:       state,
			Evidence:    evidence,
		})
	}
	return modules
}

func appendUniqueStrings(values []string, additions ...string) []string {
	seen := make(map[string]struct{}, len(values)+len(additions))
	result := make([]string, 0, len(values)+len(additions))

	for _, value := range append(values, additions...) {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}

	return result
}

func seoOperatingModuleState(number int, document *dto.SeoPublicDocument) (string, []string) {
	switch number {
	case 1, 2, 3:
		return "governed", []string{"seo-strategy-domain"}
	case 4:
		if document.AnswerReadiness != nil && document.AnswerReadiness.SpecificationVersion != "" {
			return "active", []string{document.AnswerReadiness.SpecificationVersion}
		}
		return "blocked", []string{"answer-specification-missing"}
	case 5:
		if strings.TrimSpace(document.Source.UpdatedAt) != "" && strings.TrimSpace(document.Source.Authority) != "" {
			return "active", []string{"source", "trust", "freshness"}
		}
		return "blocked", []string{"source-contract-incomplete"}
	case 6:
		if isSafePublicPath(document.Identity.CanonicalPath) {
			return "active", []string{"canonical-path", "breadcrumbs", "related-links"}
		}
		return "blocked", []string{"canonical-path-invalid"}
	case 7:
		return "active", []string{"server-semantic-layout", "react-native-client-handoff"}
	case 8:
		return "active", []string{document.SchemaVersion}
	case 9:
		return "implemented", []string{"discovery-contract", "crm-sitemap-owner"}
	case 10:
		return "active", []string{"next-ssr", "rn-client-island"}
	case 11:
		if document.IndexDecision == nil {
			return "blocked", []string{"index-decision-missing"}
		}
		if document.IndexDecision.Eligible {
			return "eligible", append([]string{document.IndexDecision.PolicyVersion}, document.IndexDecision.ReasonCodes...)
		}
		return "noindex", append([]string{document.IndexDecision.PolicyVersion}, document.IndexDecision.ReasonCodes...)
	case 12:
		if document.SEO.Title != "" && document.SEO.Description != "" {
			return "active", []string{"metadata", "breadcrumbs", "json-ld"}
		}
		return "blocked", []string{"search-appearance-incomplete"}
	case 13:
		if document.AnswerReadiness != nil && document.AnswerReadiness.Ready {
			return "ready", []string{"answer-readiness", "source-quality"}
		}
		return "blocked", []string{"answer-readiness-failed"}
	case 14:
		if document.GenerativeAI != nil && document.GenerativeAI.Ready {
			return "ready", []string{document.GenerativeAI.SchemaVersion, "citation-provenance"}
		}
		return "blocked", []string{"generative-answer-not-ready"}
	case 15:
		if document.Measurement != nil && len(document.Measurement.RequiredEvents) > 0 {
			return "instrumented", append([]string{document.Measurement.SchemaVersion}, document.Measurement.RequiredEvents...)
		}
		return "blocked", []string{"measurement-plan-missing"}
	case 16:
		return "governed", []string{"backend-decision", "admin-read-model", "pilot-gates"}
	default:
		return "not-started", nil
	}
}

func buildSeoPageBootstrap(document *dto.SeoPublicDocument) *dto.SeoPageBootstrap {
	entityKind, presentation := seoClientPresentation(document.ModuleID)
	if entityKind == "" || presentation == "" || !isSafePublicPath(document.Identity.CanonicalPath) {
		return nil
	}

	mountPolicy := "not-activated"
	if document.ModuleID == "administrative-unit" {
		mountPolicy = "client-island"
	}

	bootstrap := &dto.SeoPageBootstrap{
		ContractVersion: dto.SeoPageBootstrapContractVersion,
		Presentation:    presentation,
		Surface:         "seo-detail",
		EntityKind:      entityKind,
		EntityID:        strings.TrimSpace(document.Identity.ID),
		CanonicalSlug:   strings.TrimSpace(document.Identity.Slug),
		CanonicalPath:   document.Identity.CanonicalPath,
		InitialScreen:   "QHDetailScreen",
		MountPolicy:     mountPolicy,
		DataStrategy:    "server-bootstrap-revalidate",
		DocumentVersion: strings.TrimSpace(document.DocumentVersion),
	}
	if kind, id := parseSeoMapFocus(document.MapTarget); kind != "" && id != "" {
		bootstrap.MapFocus = &dto.SeoPageBootstrapMapFocus{Kind: kind, ID: id}
	}
	return bootstrap
}

func seoClientPresentation(moduleID string) (string, string) {
	switch moduleID {
	case "administrative-unit":
		return "administrative_unit", "administrative-unit-detail"
	case "parcel":
		return "parcel", "parcel-detail"
	case "planning-region":
		return "planning_region", "planning-region-detail"
	case "planning-project":
		return "planning_project", "planning-project-detail"
	case "planning-report":
		return "planning_report", "planning-report-detail"
	case "planning-news":
		return "planning_news", "planning-news-detail"
	case "legal-document":
		return "legal_document", "legal-document-detail"
	default:
		return "", ""
	}
}

func parseSeoMapFocus(target *dto.SeoMapTarget) (string, string) {
	if target == nil || !isSafePublicPath(target.Path) {
		return "", ""
	}
	parsed, err := url.Parse(target.Path)
	if err != nil {
		return "", ""
	}
	focus := strings.TrimSpace(parsed.Query().Get("focus"))
	separator := strings.Index(focus, ":")
	if separator <= 0 || separator >= len(focus)-1 {
		return "", ""
	}
	return focus[:separator], focus[separator+1:]
}

func buildSeoMeasurementPlan(document *dto.SeoPublicDocument) *dto.SeoMeasurementPlan {
	events := []string{
		"seo_page_view",
		"rn_activation",
		"entity_detail_open",
		"canonical_page_navigation",
		"navigation_error",
		"invalid_canonical_url",
	}
	if document.MapTarget != nil {
		events = append(events, "seo_map_cta_click", "map_entity_focus")
	}
	if len(document.RelatedLinks) > 0 || len(document.RelatedEntities) > 0 {
		events = append(events, "related_content_open")
	}
	return &dto.SeoMeasurementPlan{
		SchemaVersion:  dto.SeoMeasurementPlanSchemaVersion,
		RequiredEvents: compactUniqueStrings(events),
		CorrelationKeys: []string{
			"moduleId",
			"entityId",
			"canonicalPath",
			"documentVersion",
			"sessionId",
		},
	}
}

func buildSeoGenerativeAnswer(document *dto.SeoPublicDocument) *dto.SeoGenerativeAnswer {
	citationID := "source-primary"
	citation := dto.SeoGenerativeCitation{
		ID:         citationID,
		SourceName: firstNonEmptyRender(document.Source.SourceName, "Nguồn dữ liệu công khai"),
		SourceURL:  strings.TrimSpace(document.Source.SourceURL),
		UpdatedAt:  strings.TrimSpace(document.Source.UpdatedAt),
		Authority:  strings.TrimSpace(document.Source.Authority),
	}

	blocks := []dto.SeoGenerativeAnswerBlock{{
		ID:          "answer-summary",
		Kind:        "summary",
		Title:       document.Heading.Title,
		Text:        document.Heading.Summary,
		CitationIDs: []string{citationID},
	}}
	if len(document.Facts) > 0 {
		factIDs := make([]string, 0, len(document.Facts))
		facts := make([]string, 0, len(document.Facts))
		for _, fact := range document.Facts {
			if strings.TrimSpace(fact.ID) == "" || !isMeaningfulSeoValue(fact.Value) {
				continue
			}
			factIDs = append(factIDs, fact.ID)
			facts = append(facts, fact.Label+": "+fact.Value)
		}
		if len(facts) > 0 {
			blocks = append(blocks, dto.SeoGenerativeAnswerBlock{
				ID:          "answer-facts",
				Kind:        "facts",
				Title:       "Thông tin chính",
				Text:        strings.Join(facts, "; "),
				FactIDs:     factIDs,
				CitationIDs: []string{citationID},
			})
		}
	}

	sectionIDs := make([]string, 0, len(document.Sections))
	for _, section := range document.Sections {
		if hasMeaningfulSeoSection(section) {
			sectionIDs = append(sectionIDs, section.ID)
		}
	}
	if len(sectionIDs) > 0 {
		blocks = append(blocks, dto.SeoGenerativeAnswerBlock{
			ID:          "answer-sections",
			Kind:        "coverage",
			Title:       "Phạm vi câu trả lời",
			Text:        "Nội dung công khai được cấu trúc theo các phần đã kiểm chứng trong PublicSeoDocument.",
			SectionIDs:  sectionIDs,
			CitationIDs: []string{citationID},
		})
	}

	ready := document.AnswerReadiness != nil && document.AnswerReadiness.Ready &&
		citation.SourceName != "" && citation.UpdatedAt != ""
	return &dto.SeoGenerativeAnswer{
		SchemaVersion: dto.SeoGenerativeAnswerSchemaVersion,
		Ready:         ready,
		Blocks:        blocks,
		Citations:     []dto.SeoGenerativeCitation{citation},
		Limitations:   append([]string{}, document.Source.Limitations...),
	}
}
