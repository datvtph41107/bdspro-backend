package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/enums"

	"gorm.io/datatypes"
)

const planningProjectSourceTimeout = 15 * time.Second

func (u *SeoRenderUsecase) hydratePlanningProjectSource(
	ctx context.Context,
	page *seo_domain.SeoDomain,
	now time.Time,
) error {
	if u == nil || u.tqdProvider == nil || page == nil {
		return nil
	}
	if normalizeSeoModuleID(firstNonEmptyRender(page.TemplateKey, page.GetSchema())) != "planning-project" {
		return nil
	}

	identity := planningProjectSourceIdentity(page)
	if identity == "" {
		return fmt.Errorf("planning project source identity is missing")
	}

	sourceCtx, cancel := context.WithTimeout(ctx, planningProjectSourceTimeout)
	defer cancel()

	projection, err := u.tqdProvider.GetPlanningProjectProjection(sourceCtx, identity)
	if err != nil {
		return fmt.Errorf("sync TQD planning project projection %q: %w", identity, err)
	}

	document := buildPlanningProjectPublicDocument(projection)
	warnings := make([]string, 0)
	normalizeSeoPublicDocument(document, &warnings)

	// Source hydration only refreshes content. It never changes publication,
	// index or sitemap flags configured on the SEO page.
	document.SEO = dto.SeoPublicSEO{}

	projectionJSON, err := json.Marshal(projection)
	if err != nil {
		return fmt.Errorf("encode TQD planning project snapshot: %w", err)
	}
	documentJSON, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("encode planning project public document: %w", err)
	}

	snapshotHashBytes := sha256.Sum256(projectionJSON)
	snapshotHash := hex.EncodeToString(snapshotHashBytes[:])
	sourceUpdatedAt := projection.Source.UpdatedAt

	updates := map[string]interface{}{
		"ref_type":           enums.ESEORefTypeProject,
		"ref_id":             projection.ID,
		"ref_source":         seo_domain.SeoRefSourceTqdPlanningProject,
		"ref_label":          projection.Name,
		"ref_url":            projection.CanonicalPath,
		"source_status":      seo_domain.SeoSourceStatusLinked,
		"ref_missing":        false,
		"ref_snapshot_json":  datatypes.JSON(projectionJSON),
		"ref_hash":           snapshotHash,
		"ref_last_synced_at": now,
		"metadata":           datatypes.JSON(documentJSON),
		"template_key":       "planning-project",
		"template_version":   dto.SeoPublicDocumentSchemaVersion,
	}
	if sourceUpdatedAt != nil {
		updates["source_updated_at"] = *sourceUpdatedAt
	}
	if strings.TrimSpace(projection.Summary) != "" {
		updates["summary"] = strings.TrimSpace(projection.Summary)
	}

	if err := u.seoDomainRepo.Update(ctx, page.ID, updates); err != nil {
		return fmt.Errorf("persist planning project source snapshot: %w", err)
	}

	page.RefType = enums.ESEORefTypeProject
	projectID := projection.ID
	page.RefID = &projectID
	page.RefSource = seo_domain.SeoRefSourceTqdPlanningProject
	page.RefLabel = projection.Name
	page.RefURL = projection.CanonicalPath
	page.SourceStatus = seo_domain.SeoSourceStatusLinked
	page.RefMissing = false
	page.RefSnapshotJSON = datatypes.JSON(projectionJSON)
	page.RefHash = snapshotHash
	page.RefLastSyncedAt = &now
	page.Metadata = datatypes.JSON(documentJSON)
	page.TemplateKey = "planning-project"
	page.TemplateVersion = dto.SeoPublicDocumentSchemaVersion
	page.SourceUpdatedAt = sourceUpdatedAt
	if strings.TrimSpace(projection.Summary) != "" {
		page.Summary = strings.TrimSpace(projection.Summary)
	}
	return nil
}

func planningProjectSourceIdentity(page *seo_domain.SeoDomain) string {
	if page == nil {
		return ""
	}
	if (page.RefSource == "tqd" || page.RefSource == seo_domain.SeoRefSourceTqdPlanningProject) && page.RefID != nil && *page.RefID > 0 {
		return strconv.FormatUint(*page.RefID, 10)
	}
	path := strings.Trim(strings.TrimSpace(page.Slug), "/")
	if path == "" {
		path = strings.Trim(strings.TrimSpace(page.CanonicalURL), "/")
	}
	if index := strings.LastIndex(path, "/"); index >= 0 {
		path = path[index+1:]
	}
	return strings.TrimSpace(path)
}

func buildPlanningProjectPublicDocument(
	projection *dto.TqdPlanningProjectProjection,
) *dto.SeoPublicDocument {
	if projection == nil {
		return &dto.SeoPublicDocument{SchemaVersion: dto.SeoPublicDocumentSchemaVersion, ModuleID: "planning-project"}
	}

	canonicalPath := firstNonEmptyRender(
		projection.CanonicalPath,
		"/do-an-quy-hoach/"+firstNonEmptyRender(projection.Slug, projection.CanonicalKey, strconv.FormatUint(projection.ID, 10)),
	)
	updatedAt := "Chưa cập nhật"
	if projection.Source.UpdatedAt != nil {
		updatedAt = projection.Source.UpdatedAt.Format(time.RFC3339)
	}

	facts := []dto.SeoFact{
		{ID: "planning-type", Label: "Loại quy hoạch", Value: firstNonEmptyRender(projection.PlanningTypeName, strconv.FormatUint(uint64(projection.PlanningType), 10), "Chưa công bố")},
		{ID: "status", Label: "Trạng thái", Value: firstNonEmptyRender(projection.LegalStatusName, projection.ValidityStatus, strconv.FormatUint(uint64(projection.LegalStatus), 10), "Chưa công bố")},
		{ID: "administrative-unit", Label: "Địa bàn", Value: firstNonEmptyRender(projection.JurisdictionName, "Chưa công bố")},
		{ID: "updated-at", Label: "Cập nhật", Value: updatedAt},
	}
	if projection.TotalArea > 0 {
		facts = append(facts, dto.SeoFact{ID: "area", Label: "Diện tích", Value: fmt.Sprintf("%g ha", projection.TotalArea)})
	}
	if strings.TrimSpace(projection.PlanningLevelName) != "" {
		facts = append(facts, dto.SeoFact{ID: "planning-level", Label: "Cấp quy hoạch", Value: projection.PlanningLevelName})
	}
	if projection.ApprovalDate != nil {
		facts = append(facts, dto.SeoFact{ID: "approval-date", Label: "Ngày phê duyệt", Value: projection.ApprovalDate.Format(time.RFC3339)})
	}
	if strings.TrimSpace(projection.CurrentVersion) != "" {
		facts = append(facts, dto.SeoFact{ID: "current-version", Label: "Phiên bản", Value: projection.CurrentVersion})
	}

	timeline := make([]dto.SeoTimelineItem, 0, len(projection.Events)+2)
	for _, event := range projection.Events {
		title := strings.TrimSpace(event.Name)
		if eventType := strings.TrimSpace(event.EventType); eventType != "" {
			title = eventType + ": " + firstNonEmptyRender(title, "Mốc pháp lý")
		}
		item := dto.SeoTimelineItem{
			ID:          firstNonEmptyRender(strconv.FormatUint(event.ID, 10), normalizeSemanticID(event.Name)),
			Title:       firstNonEmptyRender(title, "Mốc pháp lý"),
			Description: firstNonEmptyRender(event.Description, event.SourceURL),
		}
		if event.EventDate != nil {
			item.Date = event.EventDate.Format(time.RFC3339)
		}
		timeline = append(timeline, item)
	}
	if len(timeline) == 0 && projection.ApprovalDate != nil {
		timeline = append(timeline, dto.SeoTimelineItem{ID: "approval-date", Title: "Phê duyệt đồ án", Date: projection.ApprovalDate.Format(time.RFC3339)})
	}
	if len(timeline) == 0 && projection.EffectiveDate != nil {
		timeline = append(timeline, dto.SeoTimelineItem{ID: "effective-date", Title: "Đồ án có hiệu lực", Date: projection.EffectiveDate.Format(time.RFC3339)})
	}

	sections := []dto.SeoSection{
		{ID: "overview", Title: "Tổng quan đồ án", Blocks: []dto.SeoContentBlock{{Type: "paragraph", Text: firstNonEmptyRender(projection.Summary, "Thông tin tổng quan đang được cập nhật từ hồ sơ đồ án.")}}},
		{ID: "research-scope", Title: "Phạm vi nghiên cứu", Blocks: []dto.SeoContentBlock{{Type: "paragraph", Text: firstNonEmptyRender(projection.ResearchScope, "Phạm vi nghiên cứu chưa được công bố đầy đủ trong nguồn hiện tại.")}}},
		{ID: "legal-timeline", Title: "Lịch sử pháp lý", Blocks: []dto.SeoContentBlock{{Type: "timeline", TimelineItems: timeline}}},
	}
	if indicators := planningIndicatorItems(projection.Indicators); len(indicators) > 0 {
		sections = append(sections, dto.SeoSection{ID: "indicators", Title: "Chỉ tiêu quy hoạch", Blocks: []dto.SeoContentBlock{{Type: "list", Items: indicators}}})
	}
	if documents := planningDocumentItems(projection.Documents); len(documents) > 0 {
		sections = append(sections, dto.SeoSection{ID: "documents", Title: "Văn bản liên quan", Blocks: []dto.SeoContentBlock{{Type: "list", Items: documents}}})
	}

	sourceName := firstNonEmptyRender(projection.Source.AuthorityName, projection.Source.Dataset, "TQD planning project projection")
	sections = append(sections, dto.SeoSection{
		ID: "source", Title: "Nguồn dữ liệu",
		Blocks: []dto.SeoContentBlock{{
			Type: "notice", Tone: "info", Title: firstNonEmptyRender(projection.Source.Dataset, "Nguồn dữ liệu"),
			Text: fmt.Sprintf("Nguồn: %s. Cập nhật: %s.", sourceName, updatedAt),
		}},
	})

	relatedLinks := []dto.SeoRelatedLink{
		{ID: "planning-map", Label: "Tra cứu trên bản đồ quy hoạch", Path: "/ban-do", Description: "Mở không gian bản đồ QHPro."},
		{ID: "planning-documents", Label: "Tài liệu quy hoạch", Path: "/quy-hoach/tai-lieu", Description: "Tra cứu tài liệu và hồ sơ công khai."},
	}
	relatedEntities := make([]dto.SeoRelatedEntity, 0, len(projection.RelatedProjects))
	for _, related := range projection.RelatedProjects {
		id := firstNonEmptyRender(strconv.FormatUint(related.ID, 10), related.Slug, related.Code, related.Name)
		if isSafePublicPath(related.PublicPath) {
			relatedLinks = append(relatedLinks, dto.SeoRelatedLink{ID: id, Label: firstNonEmptyRender(related.Name, "Đồ án liên quan"), Path: related.PublicPath})
		}
		relatedEntities = append(relatedEntities, dto.SeoRelatedEntity{
			ID: id, RelationType: "related", EntityType: "planning_project", EntityID: id,
			Title: firstNonEmptyRender(related.Name, "Đồ án liên quan"), PublicPath: related.PublicPath, MapPath: related.MapPath,
		})
	}

	document := &dto.SeoPublicDocument{
		SchemaVersion: dto.SeoPublicDocumentSchemaVersion,
		ModuleID:      "planning-project",
		ResourceType:  "planning_project",
		Identity:      dto.SeoPublicIdentity{ID: strconv.FormatUint(projection.ID, 10), Slug: projection.Slug, CanonicalPath: canonicalPath},
		Heading:       dto.SeoPublicHeading{Eyebrow: "Đồ án quy hoạch", Title: firstNonEmptyRender(projection.Name, "Đồ án quy hoạch"), Summary: firstNonEmptyRender(projection.Summary, "Thông tin công khai về "+firstNonEmptyRender(projection.Name, "đồ án quy hoạch")+".")},
		Breadcrumbs:   []dto.SeoBreadcrumb{{Label: "Trang chủ", Path: "/"}, {Label: "Đồ án quy hoạch", Path: "/do-an-quy-hoach"}, {Label: firstNonEmptyRender(projection.Name, "Đồ án quy hoạch"), Path: canonicalPath}},
		Facts:         facts, Sections: sections, RelatedLinks: relatedLinks, RelatedEntities: relatedEntities,
		FAQ: &dto.SeoFAQ{GroupKey: "planning-project-detail", Items: []dto.SeoFAQItem{
			{ID: "information-scope", Question: "Trang đồ án này cung cấp những thông tin gì?", Answer: "Trang tổng hợp loại quy hoạch, trạng thái, địa bàn, phạm vi nghiên cứu, các mốc pháp lý, chỉ tiêu và tài liệu liên quan theo dữ liệu nguồn hiện có."},
			{ID: "legal-reference", Question: "Có thể dùng nội dung trên trang làm căn cứ pháp lý cuối cùng không?", Answer: "Không. Nội dung có giá trị tra cứu và tham khảo; người dùng cần đối chiếu quyết định, bản vẽ và hồ sơ chính thức tại cơ quan có thẩm quyền."},
		}},
		MapTarget: &dto.SeoMapTarget{Path: "/ban-do?projectId=" + strconv.FormatUint(projection.ID, 10), Label: "Xem đồ án trên bản đồ"},
		Source:    dto.SeoPublicSource{SourceName: sourceName, UpdatedAt: updatedAt},
	}
	if projection.Preview.Bounds != nil {
		bounds := projection.Preview.Bounds
		document.Media = &dto.SeoMediaBundle{
			Primary: &dto.SeoMedia{ID: "planning-project-" + strconv.FormatUint(projection.ID, 10) + "-overview", Role: "overview", Alt: "Phạm vi " + firstNonEmptyRender(projection.Name, "đồ án quy hoạch"), Bounds: [4]float64{bounds.MinLongitude, bounds.MinLatitude, bounds.MaxLongitude, bounds.MaxLatitude}},
			Items:   []dto.SeoMedia{}, Completeness: firstNonEmptyRender(projection.Preview.Completeness, "complete"), FitMode: "contain", Note: "Khung xem trước sử dụng toàn bộ bounds của geometry để tránh cắt thiếu phạm vi.",
		}
	}
	for _, item := range projection.Documents {
		if strings.TrimSpace(item.Thumbnail) != "" {
			document.SEO.ImageURL = item.Thumbnail
			break
		}
	}
	return document
}

func planningIndicatorItems(items []dto.TqdPlanningIndicator) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		name, value := strings.TrimSpace(item.Name), strings.TrimSpace(item.Value)
		if name == "" || value == "" {
			continue
		}
		text := name + ": " + value
		if unit := strings.TrimSpace(item.Unit); unit != "" {
			text += " " + unit
		}
		out = append(out, text)
	}
	return out
}

func planningDocumentItems(items []dto.TqdPlanningDocument) []string {
	labels := make([]string, len(items))
	labelCounts := make(map[string]int, len(items))

	for index, item := range items {
		title := firstNonEmptyRender(
			item.Title,
			item.Code,
			"Văn bản quy hoạch",
		)

		if kind := strings.TrimSpace(item.DocumentTypeName); kind != "" {
			title = kind + ": " + title
		}

		labels[index] = title
		labelCounts[strings.ToLower(strings.TrimSpace(title))]++
	}

	out := make([]string, 0, len(items))

	for index, item := range items {
		title := labels[index]
		key := strings.ToLower(strings.TrimSpace(title))

		if labelCounts[key] > 1 {
			switch {
			case item.ID > 0:
				title += " · Hồ sơ #" +
					strconv.FormatUint(item.ID, 10)

			case strings.TrimSpace(item.Code) != "":
				title += " · " + strings.TrimSpace(item.Code)

			default:
				title += fmt.Sprintf(" · Hồ sơ %d", index+1)
			}
		}

		out = append(out, title)
	}

	return out
}
