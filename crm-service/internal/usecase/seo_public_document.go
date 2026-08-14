package usecase

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
)

const (
	maxSeoFacts            = 24
	maxSeoSections         = 20
	maxSeoBlocksPerSection = 20
	maxSeoRelatedLinks     = 24
	maxSeoRelatedEntities  = 24
	maxSeoFAQItems         = 16
)

var legacyHTMLTagPattern = regexp.MustCompile(`(?s)<[^>]*>`)
var legacyHTMLBreakPattern = regexp.MustCompile(`(?i)</?(p|div|section|article|header|footer|h[1-6]|li|br)[^>]*>`)
var compactWhitespacePattern = regexp.MustCompile(`\s+`)

type publicDocumentMetadataEnvelope struct {
	PublicDocument *dto.SeoPublicDocument `json:"publicDocument"`
}

func composeSeoPublicDocument(page *seo_domain.SeoDomain, generatedAt time.Time) (*dto.SeoPublicDocument, []string, error) {
	if page == nil {
		return nil, nil, fmt.Errorf("seo page is nil")
	}

	page.NormalizeLifecycle()
	moduleID := normalizeSeoModuleID(firstNonEmptyRender(page.TemplateKey, page.GetSchema(), "default"))
	if moduleID == "" {
		return nil, nil, fmt.Errorf("unsupported SEO module %q", firstNonEmptyRender(page.TemplateKey, page.GetSchema()))
	}
	warnings := make([]string, 0, 8)
	document := buildDefaultSeoPublicDocument(page, generatedAt)

	metadataDocument, metadataWarnings := parseSeoPublicDocumentMetadata(json.RawMessage(page.Metadata))
	warnings = append(warnings, metadataWarnings...)
	if metadataDocument != nil {
		mergeSeoPublicDocument(document, metadataDocument)
	}

	applySeoDomainOwnedFields(document, page, generatedAt)
	mergeDatabaseRelationships(document, page)
	normalizeSeoPublicDocument(document, &warnings)
	return document, compactUniqueStrings(warnings), nil
}

func buildDefaultSeoPublicDocument(page *seo_domain.SeoDomain, generatedAt time.Time) *dto.SeoPublicDocument {
	moduleID := normalizeSeoModuleID(firstNonEmptyRender(page.TemplateKey, page.GetSchema(), "default"))
	resourceType := resourceTypeForModule(moduleID)
	canonicalPath := canonicalPathFromURL(page.CanonicalURL, page.Slug)
	updatedAt := chooseSeoUpdatedAt(page, generatedAt)
	summary := firstNonEmptyRender(page.Summary, page.Description)
	contentText := legacyContentToPlainText(page.Content)

	sections := []dto.SeoSection{
		{
			ID:    "overview",
			Title: "Tổng quan",
			Blocks: []dto.SeoContentBlock{{
				Type: "paragraph",
				Text: firstNonEmptyRender(summary, contentText, "Thông tin tổng quan đang được hoàn thiện."),
			}},
		},
	}
	if contentText != "" && contentText != summary {
		sections = append(sections, dto.SeoSection{
			ID:     "research-scope",
			Title:  "Phạm vi thông tin",
			Blocks: []dto.SeoContentBlock{{Type: "paragraph", Text: contentText}},
		})
	}

	document := &dto.SeoPublicDocument{
		SchemaVersion: dto.SeoPublicDocumentSchemaVersion,
		ModuleID:      moduleID,
		ResourceType:  resourceType,
		Identity: dto.SeoPublicIdentity{
			ID:            fmt.Sprintf("seo-%d", page.ID),
			Slug:          strings.Trim(strings.TrimSpace(page.Slug), "/"),
			CanonicalPath: canonicalPath,
		},
		Heading: dto.SeoPublicHeading{
			Eyebrow: moduleEyebrow(moduleID),
			Title:   strings.TrimSpace(page.Title),
			Summary: summary,
		},
		Breadcrumbs: defaultBreadcrumbs(moduleID, canonicalPath, page.Title),
		Facts: []dto.SeoFact{
			{ID: "planning-type", Label: "Loại quy hoạch", Value: defaultPlanningType(page, moduleID)},
			{ID: "status", Label: "Trạng thái", Value: pageStatusLabel(page.PageStatus)},
			{ID: "updated-at", Label: "Cập nhật", Value: updatedAt},
		},
		Sections:     sections,
		RelatedLinks: []dto.SeoRelatedLink{},
		SEO: dto.SeoPublicSEO{
			Title:         strings.TrimSpace(page.Title),
			Description:   strings.TrimSpace(page.Description),
			CanonicalPath: canonicalPath,
			Indexable:     page.IsPublicIndexable(),
			Follow:        true,
		},
		Source: defaultSeoSource(page, updatedAt),
	}

	if mapTarget := defaultMapTarget(page, moduleID); mapTarget != nil {
		document.MapTarget = mapTarget
	}

	return document
}

func parseSeoPublicDocumentMetadata(raw json.RawMessage) (*dto.SeoPublicDocument, []string) {
	text := strings.TrimSpace(string(raw))
	if text == "" || text == "{}" || text == "null" {
		return nil, nil
	}

	var root map[string]json.RawMessage
	if err := json.Unmarshal(raw, &root); err != nil {
		return nil, []string{"metadata_invalid_json"}
	}

	if nestedRaw, ok := root["publicDocument"]; ok {
		var document dto.SeoPublicDocument
		if err := json.Unmarshal(nestedRaw, &document); err != nil {
			return nil, []string{"public_document_invalid"}
		}
		return &document, nil
	}

	if schemaRaw, ok := root["schemaVersion"]; ok {
		var schemaVersion string
		_ = json.Unmarshal(schemaRaw, &schemaVersion)
		if strings.HasPrefix(schemaVersion, "qhpro-public-seo-document/") {
			var document dto.SeoPublicDocument
			if err := json.Unmarshal(raw, &document); err != nil {
				return nil, []string{"public_document_invalid"}
			}
			return &document, nil
		}
	}

	// Legacy metadata may be a standalone JSON-LD object. It is intentionally
	// ignored as document input because Gate 2.2 generates JSON-LD from the
	// versioned semantic document, avoiding two structured-data sources.
	if _, hasContext := root["@context"]; hasContext {
		return nil, []string{"legacy_jsonld_metadata_ignored"}
	}

	var envelope publicDocumentMetadataEnvelope
	if err := json.Unmarshal(raw, &envelope); err == nil && envelope.PublicDocument != nil {
		return envelope.PublicDocument, nil
	}

	return nil, []string{"public_document_metadata_missing"}
}

func mergeSeoPublicDocument(target *dto.SeoPublicDocument, source *dto.SeoPublicDocument) {
	if target == nil || source == nil {
		return
	}

	if strings.TrimSpace(source.SchemaVersion) != "" {
		target.SchemaVersion = source.SchemaVersion
	}
	if strings.TrimSpace(source.ModuleID) != "" {
		target.ModuleID = source.ModuleID
	}
	if strings.TrimSpace(source.ResourceType) != "" {
		target.ResourceType = source.ResourceType
	}
	if strings.TrimSpace(source.Heading.Eyebrow) != "" {
		target.Heading.Eyebrow = source.Heading.Eyebrow
	}
	if strings.TrimSpace(source.Heading.Title) != "" {
		target.Heading.Title = source.Heading.Title
	}
	if strings.TrimSpace(source.Heading.Summary) != "" {
		target.Heading.Summary = source.Heading.Summary
	}
	if len(source.Breadcrumbs) > 0 {
		target.Breadcrumbs = source.Breadcrumbs
	}
	if len(source.Facts) > 0 {
		target.Facts = source.Facts
	}
	if len(source.Sections) > 0 {
		target.Sections = source.Sections
	}
	if len(source.RelatedLinks) > 0 {
		target.RelatedLinks = source.RelatedLinks
	}
	if len(source.RelatedEntities) > 0 {
		target.RelatedEntities = source.RelatedEntities
	}
	if source.Media != nil {
		target.Media = source.Media
	}
	if source.FAQ != nil {
		target.FAQ = source.FAQ
	}
	if source.MapTarget != nil {
		target.MapTarget = source.MapTarget
	}
	if strings.TrimSpace(source.SEO.ImageURL) != "" {
		target.SEO.ImageURL = source.SEO.ImageURL
	}
	mergeSeoSource(&target.Source, source.Source)
}

func mergeSeoSource(target *dto.SeoPublicSource, source dto.SeoPublicSource) {
	if target == nil {
		return
	}
	if strings.TrimSpace(source.SourceName) != "" {
		target.SourceName = source.SourceName
	}
	if strings.TrimSpace(source.SourceURL) != "" {
		target.SourceURL = source.SourceURL
	}
	if strings.TrimSpace(source.UpdatedAt) != "" {
		target.UpdatedAt = source.UpdatedAt
	}
	if strings.TrimSpace(source.PublishedAt) != "" {
		target.PublishedAt = source.PublishedAt
	}
}

func applySeoDomainOwnedFields(document *dto.SeoPublicDocument, page *seo_domain.SeoDomain, generatedAt time.Time) {
	if document == nil || page == nil {
		return
	}

	moduleID := normalizeSeoModuleID(firstNonEmptyRender(page.TemplateKey, document.ModuleID, page.GetSchema(), "default"))
	canonicalPath := canonicalPathFromURL(page.CanonicalURL, page.Slug)
	summary := firstNonEmptyRender(page.Summary, page.Description, document.Heading.Summary)

	document.SchemaVersion = dto.SeoPublicDocumentSchemaVersion
	document.DocumentVersion = seoPublicDocumentVersion(page)
	document.ModuleID = moduleID
	document.ResourceType = resourceTypeForModule(moduleID)
	document.Identity.ID = fmt.Sprintf("seo-%d", page.ID)
	document.Identity.Slug = strings.Trim(strings.TrimSpace(page.Slug), "/")
	document.Identity.CanonicalPath = canonicalPath
	document.Heading.Title = strings.TrimSpace(page.Title)
	document.Heading.Summary = summary
	if strings.TrimSpace(document.Heading.Eyebrow) == "" {
		document.Heading.Eyebrow = moduleEyebrow(moduleID)
	}

	document.SEO.Title = strings.TrimSpace(page.Title)
	document.SEO.Description = strings.TrimSpace(page.Description)
	document.SEO.CanonicalPath = canonicalPath
	document.SEO.Indexable = page.IsPublicIndexable()
	document.SEO.Follow = true
	if strings.TrimSpace(document.Source.UpdatedAt) == "" {
		document.Source.UpdatedAt = chooseSeoUpdatedAt(page, generatedAt)
	}
	if page.PublishedAt != nil && strings.TrimSpace(document.Source.PublishedAt) == "" {
		document.Source.PublishedAt = page.PublishedAt.Format(time.RFC3339)
	}
}

func seoPublicDocumentVersion(page *seo_domain.SeoDomain) string {
	if page == nil {
		return ""
	}

	sourceUpdatedAt := ""
	if page.SourceUpdatedAt != nil {
		sourceUpdatedAt = page.SourceUpdatedAt.UTC().Format(time.RFC3339Nano)
	}
	publishedAt := ""
	if page.PublishedAt != nil {
		publishedAt = page.PublishedAt.UTC().Format(time.RFC3339Nano)
	}

	// The public-document version deliberately excludes render output fields
	// (RenderedHTML, StaticHtmlHash, GeneratedAt and UpdatedAt). Those fields are
	// written after composition and would make the JSON document version differ
	// from the version embedded in the rendered HTML generated from that same
	// semantic document. Only source, content, template and publication inputs
	// participate in the version fingerprint.
	material := strings.Join([]string{
		dto.SeoPublicDocumentSchemaVersion,
		strconv.FormatUint(page.ID, 10),
		strings.TrimSpace(page.Slug),
		strings.TrimSpace(page.CanonicalURL),
		strings.TrimSpace(page.RefHash),
		strings.TrimSpace(page.TemplateKey),
		strings.TrimSpace(page.TemplateVersion),
		strings.TrimSpace(page.Title),
		strings.TrimSpace(page.Description),
		strings.TrimSpace(page.Summary),
		strings.TrimSpace(page.Content),
		strings.TrimSpace(page.PageStatus),
		strconv.FormatBool(page.Published),
		strconv.FormatBool(page.IsIndex),
		strconv.FormatBool(page.IsSiteMap),
		sourceUpdatedAt,
		publishedAt,
	}, "|")
	digest := sha256.Sum256([]byte(material))
	return fmt.Sprintf("qhpro-public-document/%x", digest[:16])
}

func mergeDatabaseRelationships(document *dto.SeoPublicDocument, page *seo_domain.SeoDomain) {
	if document == nil || page == nil {
		return
	}

	links := append([]dto.SeoRelatedLink{}, document.RelatedLinks...)
	for _, item := range page.InternalLinks {
		path := firstNonEmptyRender(item.Link)
		if path == "" && item.ChildSeo != nil {
			path = canonicalPathFromURL(item.ChildSeo.CanonicalURL, item.ChildSeo.Slug)
		}
		if !isSafePublicPath(path) {
			continue
		}
		label := firstNonEmptyRender(item.Title)
		if label == "" && item.ChildSeo != nil {
			label = item.ChildSeo.Title
		}
		if label == "" {
			continue
		}
		links = append(links, dto.SeoRelatedLink{
			ID:    fmt.Sprintf("internal-link-%d", item.ID),
			Label: label,
			Path:  path,
		})
	}
	document.RelatedLinks = dedupeRelatedLinks(links)

	entities := append([]dto.SeoRelatedEntity{}, document.RelatedEntities...)
	for _, item := range page.Relatives {
		if item.ChildSeo == nil {
			continue
		}
		child := item.ChildSeo
		publicPath := canonicalPathFromURL(child.CanonicalURL, child.Slug)
		if !isSafePublicPath(publicPath) {
			publicPath = ""
		}
		entities = append(entities, dto.SeoRelatedEntity{
			ID:           fmt.Sprintf("relative-%d", item.ID),
			RelationType: normalizeRelationType(item.RelationType),
			EntityType:   resourceTypeForModule(normalizeSeoModuleID(firstNonEmptyRender(child.TemplateKey, child.GetSchema()))),
			EntityID:     strconv.FormatUint(child.ID, 10),
			Title:        firstNonEmptyRender(child.Title, child.RefLabel, child.Slug),
			Description:  child.Summary,
			PublicPath:   publicPath,
		})
	}
	document.RelatedEntities = dedupeRelatedEntities(entities)
}

func normalizeSeoPublicDocument(document *dto.SeoPublicDocument, warnings *[]string) {
	if document == nil {
		return
	}

	document.SchemaVersion = dto.SeoPublicDocumentSchemaVersion
	document.ModuleID = normalizeSeoModuleID(document.ModuleID)
	document.ResourceType = resourceTypeForModule(document.ModuleID)
	document.Heading.Eyebrow = trimLimit(document.Heading.Eyebrow, 120)
	document.Heading.Title = trimLimit(document.Heading.Title, 240)
	document.Heading.Summary = trimLimit(document.Heading.Summary, 1000)
	document.Breadcrumbs = normalizeBreadcrumbs(document.Breadcrumbs, document.Identity.CanonicalPath, document.Heading.Title)
	document.Facts = normalizeFacts(document.Facts)
	document.Sections = normalizeSections(document.Sections)
	document.RelatedLinks = dedupeRelatedLinks(document.RelatedLinks)
	document.RelatedEntities = dedupeRelatedEntities(document.RelatedEntities)
	document.FAQ = normalizeFAQ(document.FAQ)
	document.MapTarget = normalizeMapTarget(document.MapTarget)
	document.Media = normalizeMediaBundle(document.Media)
	document.Source = normalizeSource(document.Source)

	if len(document.Facts) == 0 {
		*warnings = append(*warnings, "public_document_has_no_facts")
	}
	if len(document.Sections) == 0 {
		*warnings = append(*warnings, "public_document_has_no_sections")
	}
}

func defaultSeoSource(page *seo_domain.SeoDomain, updatedAt string) dto.SeoPublicSource {
	sourceName := ""
	sourceURL := ""
	if page != nil {
		sourceName = firstNonEmptyRender(page.RefLabel, page.RefSource)
		sourceURL = page.RefURL
	}

	return dto.SeoPublicSource{
		SourceName: sourceName,
		SourceURL:  sourceURL,
		UpdatedAt:  updatedAt,
	}
}

func defaultBreadcrumbs(moduleID string, canonicalPath string, title string) []dto.SeoBreadcrumb {
	label, collectionPath := moduleCollection(moduleID)
	return []dto.SeoBreadcrumb{
		{Label: "Trang chủ", Path: "/"},
		{Label: label, Path: collectionPath},
		{Label: firstNonEmptyRender(title, "Chi tiết"), Path: canonicalPath},
	}
}

func moduleCollection(moduleID string) (string, string) {
	switch moduleID {
	case "administrative-unit":
		return "Địa bàn", "/dia-ban"
	case "parcel":
		return "Thửa đất", "/thua-dat"
	case "planning-region":
		return "Vùng quy hoạch", "/vung-quy-hoach"
	case "planning-project":
		return "Đồ án quy hoạch", "/do-an-quy-hoach"
	case "planning-news":
		return "Tin quy hoạch", "/quy-hoach/tin-tuc"
	case "planning-report":
		return "Báo cáo quy hoạch", "/quy-hoach/bao-cao"
	default:
		return "Quy hoạch", "/quy-hoach"
	}
}

func moduleEyebrow(moduleID string) string {
	label, _ := moduleCollection(moduleID)
	return label
}

func defaultPlanningType(page *seo_domain.SeoDomain, moduleID string) string {
	if page != nil && strings.TrimSpace(page.RefLabel) != "" {
		return strings.TrimSpace(page.RefLabel)
	}
	label, _ := moduleCollection(moduleID)
	return label
}

func defaultMapTarget(page *seo_domain.SeoDomain, moduleID string) *dto.SeoMapTarget {
	if page == nil {
		return nil
	}
	if isSafePublicPath(page.DeepLink) {
		return &dto.SeoMapTarget{Path: page.DeepLink, Label: "Xem trên bản đồ"}
	}
	if page.RefID == nil {
		return nil
	}
	queryKey := "planningId"
	switch moduleID {
	case "administrative-unit":
		queryKey = "administrativeId"
	case "parcel":
		queryKey = "parcelId"
	case "planning-region":
		queryKey = "regionId"
	}
	return &dto.SeoMapTarget{
		Path:  "/ban-do?" + queryKey + "=" + strconv.FormatUint(*page.RefID, 10),
		Label: "Xem trên bản đồ",
	}
}

func chooseSeoUpdatedAt(page *seo_domain.SeoDomain, generatedAt time.Time) string {
	if page != nil {
		if page.SourceUpdatedAt != nil {
			return page.SourceUpdatedAt.Format(time.RFC3339)
		}
		if page.RefLastSyncedAt != nil {
			return page.RefLastSyncedAt.Format(time.RFC3339)
		}
		if !page.UpdatedAt.IsZero() {
			return page.UpdatedAt.Format(time.RFC3339)
		}
	}
	return generatedAt.Format(time.RFC3339)
}

func canonicalPathFromURL(canonicalURL string, slug string) string {
	canonicalURL = strings.TrimSpace(canonicalURL)
	if parsed, err := url.Parse(canonicalURL); err == nil && parsed.Path != "" {
		path := parsed.EscapedPath()
		if parsed.RawQuery != "" {
			path += "?" + parsed.RawQuery
		}
		if isSafePublicPath(path) {
			return path
		}
	}
	path := "/" + strings.Trim(strings.TrimSpace(slug), "/")
	if isSafePublicPath(path) {
		return path
	}
	return "/"
}

func legacyContentToPlainText(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	value = legacyHTMLBreakPattern.ReplaceAllString(value, "\n")
	value = legacyHTMLTagPattern.ReplaceAllString(value, " ")
	value = html.UnescapeString(value)
	lines := strings.Split(value, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = compactWhitespacePattern.ReplaceAllString(strings.TrimSpace(line), " ")
		if line != "" {
			out = append(out, line)
		}
	}
	return trimLimit(strings.Join(out, "\n\n"), 6000)
}

func normalizeSeoModuleID(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "_", "-")
	switch value {
	case "administrative-unit", "administrative-area":
		return "administrative-unit"
	case "parcel", "planning-parcel":
		return "parcel"
	case "planning-region":
		return "planning-region"
	case "planning-project":
		return "planning-project"
	case "planning-news", "article", "news-article":
		return "planning-news"
	case "planning-report", "report", "analysis-report":
		return "planning-report"
	case "default", "generic", "generic-page":
		return "default"
	default:
		return ""
	}
}

func resourceTypeForModule(moduleID string) string {
	switch normalizeSeoModuleID(moduleID) {
	case "administrative-unit":
		return "administrative_area"
	case "parcel":
		return "parcel"
	case "planning-region":
		return "planning_region"
	case "planning-news":
		return "article"
	case "planning-report":
		return "report"
	case "planning-project":
		return "planning_project"
	case "default":
		return "generic_page"
	default:
		return ""
	}
}

func pageStatusLabel(value string) string {
	switch strings.TrimSpace(value) {
	case seo_domain.SeoPageStatusPublished:
		return "Đã xuất bản"
	case seo_domain.SeoPageStatusArchived:
		return "Đã lưu trữ"
	default:
		return "Bản nháp"
	}
}

func normalizeBreadcrumbs(items []dto.SeoBreadcrumb, canonicalPath string, title string) []dto.SeoBreadcrumb {
	out := make([]dto.SeoBreadcrumb, 0, len(items)+1)
	seen := map[string]struct{}{}
	for _, item := range items {
		label := trimLimit(item.Label, 180)
		path := strings.TrimSpace(item.Path)
		if label == "" || !isSafePublicPath(path) {
			continue
		}
		key := path + "\x00" + label
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, dto.SeoBreadcrumb{Label: label, Path: path})
		if len(out) >= 8 {
			break
		}
	}
	if len(out) == 0 {
		out = append(out, dto.SeoBreadcrumb{Label: "Trang chủ", Path: "/"})
	}
	if canonicalPath != "" {
		lastPath := out[len(out)-1].Path
		if lastPath != canonicalPath {
			out = append(out, dto.SeoBreadcrumb{Label: firstNonEmptyRender(title, "Chi tiết"), Path: canonicalPath})
		}
	}
	return out
}

func normalizeFacts(items []dto.SeoFact) []dto.SeoFact {
	out := make([]dto.SeoFact, 0, len(items))
	seen := map[string]struct{}{}
	for _, item := range items {
		id := normalizeSemanticID(item.ID)
		label := trimLimit(item.Label, 160)
		value := trimLimit(item.Value, 500)
		if id == "" || label == "" || value == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, dto.SeoFact{ID: id, Label: label, Value: value, Description: trimLimit(item.Description, 800)})
		if len(out) >= maxSeoFacts {
			break
		}
	}
	return out
}

func normalizeSections(items []dto.SeoSection) []dto.SeoSection {
	out := make([]dto.SeoSection, 0, len(items))
	seen := map[string]struct{}{}
	for _, item := range items {
		id := normalizeSemanticID(item.ID)
		title := trimLimit(item.Title, 220)
		if id == "" || title == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		blocks := normalizeContentBlocks(item.Blocks)
		if len(blocks) == 0 {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, dto.SeoSection{ID: id, Title: title, Blocks: blocks})
		if len(out) >= maxSeoSections {
			break
		}
	}
	return out
}

func normalizeContentBlocks(items []dto.SeoContentBlock) []dto.SeoContentBlock {
	out := make([]dto.SeoContentBlock, 0, len(items))
	for _, item := range items {
		typeName := strings.ToLower(strings.TrimSpace(item.Type))
		switch typeName {
		case "paragraph":
			text := trimLimit(item.Text, 6000)
			if text != "" {
				out = append(out, dto.SeoContentBlock{Type: "paragraph", Text: text})
			}
		case "list":
			listItems := compactUniqueStrings(item.Items)
			if len(listItems) > 20 {
				listItems = listItems[:20]
			}
			for index := range listItems {
				listItems[index] = trimLimit(listItems[index], 1000)
			}
			if len(listItems) > 0 {
				out = append(out, dto.SeoContentBlock{Type: "list", Items: listItems})
			}
		case "notice":
			text := trimLimit(item.Text, 3000)
			if text != "" {
				tone := strings.ToLower(strings.TrimSpace(item.Tone))
				if tone != "warning" {
					tone = "info"
				}
				out = append(out, dto.SeoContentBlock{Type: "notice", Tone: tone, Title: trimLimit(item.Title, 180), Text: text})
			}
		case "timeline":
			timeline := make([]dto.SeoTimelineItem, 0, len(item.TimelineItems))
			for _, timelineItem := range item.TimelineItems {
				title := trimLimit(timelineItem.Title, 300)
				if title == "" {
					continue
				}
				timeline = append(timeline, dto.SeoTimelineItem{
					ID:          firstNonEmptyRender(normalizeSemanticID(timelineItem.ID), fmt.Sprintf("event-%d", len(timeline)+1)),
					Title:       title,
					Date:        trimLimit(timelineItem.Date, 80),
					Description: trimLimit(timelineItem.Description, 1500),
				})
				if len(timeline) >= 20 {
					break
				}
			}
			if len(timeline) > 0 {
				out = append(out, dto.SeoContentBlock{Type: "timeline", TimelineItems: timeline})
			}
		}
		if len(out) >= maxSeoBlocksPerSection {
			break
		}
	}
	return out
}

func normalizeFAQ(faq *dto.SeoFAQ) *dto.SeoFAQ {
	if faq == nil {
		return nil
	}
	out := &dto.SeoFAQ{GroupKey: trimLimit(faq.GroupKey, 120), Items: make([]dto.SeoFAQItem, 0, len(faq.Items))}
	seen := map[string]struct{}{}
	for _, item := range faq.Items {
		question := trimLimit(item.Question, 500)
		answer := trimLimit(item.Answer, 3000)
		if question == "" || answer == "" {
			continue
		}
		key := strings.ToLower(question)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out.Items = append(out.Items, dto.SeoFAQItem{ID: firstNonEmptyRender(normalizeSemanticID(item.ID), fmt.Sprintf("faq-%d", len(out.Items)+1)), Question: question, Answer: answer})
		if len(out.Items) >= maxSeoFAQItems {
			break
		}
	}
	if len(out.Items) == 0 {
		return nil
	}
	return out
}

func normalizeMapTarget(target *dto.SeoMapTarget) *dto.SeoMapTarget {
	if target == nil || !isSafePublicPath(target.Path) {
		return nil
	}
	return &dto.SeoMapTarget{Path: target.Path, Label: firstNonEmptyRender(trimLimit(target.Label, 160), "Xem trên bản đồ")}
}

func normalizeMediaBundle(bundle *dto.SeoMediaBundle) *dto.SeoMediaBundle {
	if bundle == nil {
		return nil
	}
	out := &dto.SeoMediaBundle{Completeness: bundle.Completeness, FitMode: bundle.FitMode, Note: trimLimit(bundle.Note, 1000)}
	if out.Completeness != "complete" && out.Completeness != "partial" && out.Completeness != "centroid_only" {
		out.Completeness = "missing"
	}
	if out.FitMode != "cover" && out.FitMode != "center" {
		out.FitMode = "contain"
	}
	for _, item := range bundle.Items {
		if media := normalizeMedia(item); media != nil {
			out.Items = append(out.Items, *media)
		}
		if len(out.Items) >= 20 {
			break
		}
	}
	if bundle.Primary != nil {
		out.Primary = normalizeMedia(*bundle.Primary)
	}
	if out.Primary == nil && len(out.Items) > 0 {
		primary := out.Items[0]
		out.Primary = &primary
	}
	if out.Primary == nil && len(out.Items) == 0 {
		return nil
	}
	sort.SliceStable(out.Items, func(i, j int) bool { return out.Items[i].Order < out.Items[j].Order })
	return out
}

func normalizeMedia(item dto.SeoMedia) *dto.SeoMedia {
	alt := trimLimit(item.Alt, 500)
	if alt == "" {
		return nil
	}
	item.ID = firstNonEmptyRender(normalizeSemanticID(item.ID), "media")
	item.Role = firstNonEmptyRender(normalizeSemanticID(item.Role), "overview")
	item.Alt = alt
	item.URL = strings.TrimSpace(item.URL)
	if item.URL != "" && !isSafeMediaURL(item.URL) {
		item.URL = ""
	}
	if item.Width < 0 {
		item.Width = 0
	}
	if item.Height < 0 {
		item.Height = 0
	}
	return &item
}

func normalizeSource(source dto.SeoPublicSource) dto.SeoPublicSource {
	source.SourceName = trimLimit(source.SourceName, 300)
	source.SourceURL = strings.TrimSpace(source.SourceURL)
	if source.SourceURL != "" && !isSafeMediaURL(source.SourceURL) {
		source.SourceURL = ""
	}
	source.UpdatedAt = trimLimit(source.UpdatedAt, 80)
	source.PublishedAt = trimLimit(source.PublishedAt, 80)
	return source
}

func dedupeRelatedLinks(items []dto.SeoRelatedLink) []dto.SeoRelatedLink {
	out := make([]dto.SeoRelatedLink, 0, len(items))
	seen := map[string]struct{}{}
	for _, item := range items {
		path := strings.TrimSpace(item.Path)
		label := trimLimit(item.Label, 240)
		if !isSafePublicPath(path) || label == "" {
			continue
		}
		if _, ok := seen[path]; ok {
			continue
		}
		seen[path] = struct{}{}
		out = append(out, dto.SeoRelatedLink{ID: firstNonEmptyRender(normalizeSemanticID(item.ID), fmt.Sprintf("related-%d", len(out)+1)), Label: label, Path: path, Description: trimLimit(item.Description, 800)})
		if len(out) >= maxSeoRelatedLinks {
			break
		}
	}
	return out
}

func dedupeRelatedEntities(items []dto.SeoRelatedEntity) []dto.SeoRelatedEntity {
	out := make([]dto.SeoRelatedEntity, 0, len(items))
	seen := map[string]struct{}{}
	for _, item := range items {
		title := trimLimit(item.Title, 240)
		if title == "" {
			continue
		}
		key := firstNonEmptyRender(item.ID, item.EntityType+":"+item.EntityID, title)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		publicPath := strings.TrimSpace(item.PublicPath)
		if publicPath != "" && !isSafePublicPath(publicPath) {
			publicPath = ""
		}
		mapPath := strings.TrimSpace(item.MapPath)
		if mapPath != "" && !isSafePublicPath(mapPath) {
			mapPath = ""
		}
		out = append(out, dto.SeoRelatedEntity{
			ID:           firstNonEmptyRender(normalizeSemanticID(item.ID), fmt.Sprintf("entity-%d", len(out)+1)),
			RelationType: normalizeRelationType(item.RelationType),
			EntityType:   trimLimit(item.EntityType, 100),
			EntityID:     trimLimit(item.EntityID, 120),
			Title:        title,
			Description:  trimLimit(item.Description, 1000),
			PublicPath:   publicPath,
			MapPath:      mapPath,
		})
		if len(out) >= maxSeoRelatedEntities {
			break
		}
	}
	return out
}

func normalizeRelationType(value string) string {
	return normalizeEnum(value, []string{"intersects", "contains", "belongs_to", "references", "related"}, "related")
}

func normalizeEnum(value string, allowed []string, fallback string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	for _, item := range allowed {
		if value == item {
			return value
		}
	}
	return fallback
}

func normalizeSemanticID(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "_", "-")
	var out strings.Builder
	lastDash := false
	for _, r := range value {
		valid := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		if valid {
			out.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash && out.Len() > 0 {
			out.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(out.String(), "-")
}

func isSafePublicPath(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" || !strings.HasPrefix(value, "/") || strings.HasPrefix(value, "//") {
		return false
	}
	lower := strings.ToLower(value)
	return !strings.Contains(lower, "undefined") && !strings.Contains(lower, "null") && !strings.Contains(lower, "[object object]") && !strings.ContainsAny(value, "\r\n\x00")
}

func isSafeMediaURL(value string) bool {
	value = strings.TrimSpace(value)
	if isSafePublicPath(value) {
		return true
	}
	parsed, err := url.Parse(value)
	if err != nil {
		return false
	}
	return (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Host != ""
}

func trimLimit(value string, limit int) string {
	value = strings.TrimSpace(value)
	if limit <= 0 {
		return ""
	}
	runes := []rune(value)
	if len(runes) > limit {
		return string(runes[:limit])
	}
	return value
}

func compactUniqueStrings(items []string) []string {
	out := make([]string, 0, len(items))
	seen := map[string]struct{}{}
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}

func moduleRequiresMapTarget(moduleID string) bool {
	switch moduleID {
	case "planning-news", "planning-report":
		return false
	default:
		return true
	}
}

func moduleRequiresSourceDisclosure(moduleID string) bool {
	return moduleID != ""
}
