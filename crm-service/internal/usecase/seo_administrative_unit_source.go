package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/enums"

	"gorm.io/datatypes"
)

const administrativeUnitSourceTimeout = 15 * time.Second

func (u *SeoRenderUsecase) hydrateSeoSource(ctx context.Context, page *seo_domain.SeoDomain, now time.Time) error {
	if page == nil {
		return nil
	}
	switch normalizeSeoModuleID(firstNonEmptyRender(page.TemplateKey, page.GetSchema())) {
	case "administrative-unit":
		return u.hydrateAdministrativeUnitSource(ctx, page, now)
	case "planning-project":
		return u.hydratePlanningProjectSource(ctx, page, now)
	default:
		return nil
	}
}

func (u *SeoRenderUsecase) hydrateAdministrativeUnitSource(
	ctx context.Context,
	page *seo_domain.SeoDomain,
	now time.Time,
) error {
	if u == nil || u.tqdProvider == nil || page == nil {
		return nil
	}
	if normalizeSeoModuleID(firstNonEmptyRender(page.TemplateKey, page.GetSchema())) != "administrative-unit" {
		return nil
	}

	identity := administrativeUnitSourceIdentity(page)
	if identity == "" {
		return fmt.Errorf("administrative unit source identity is missing")
	}

	sourceCtx, cancel := context.WithTimeout(ctx, administrativeUnitSourceTimeout)
	defer cancel()

	projection, err := u.tqdProvider.GetAdministrativeUnitProjection(sourceCtx, identity)
	if err != nil {
		return fmt.Errorf("sync TQD administrative unit projection %q: %w", identity, err)
	}

	pageCanonicalPath := canonicalPathFromURL(page.CanonicalURL, page.Slug)
	projectionCanonicalPath := strings.TrimSpace(projection.CanonicalPath)
	if pageCanonicalPath != projectionCanonicalPath {
		return fmt.Errorf(
			"administrative unit canonical mismatch: page=%q projection=%q",
			pageCanonicalPath,
			projectionCanonicalPath,
		)
	}

	document := buildAdministrativeUnitPublicDocument(projection)
	warnings := make([]string, 0)
	normalizeSeoPublicDocument(document, &warnings)

	// Source hydration only refreshes content. It never changes publication,
	// index or sitemap flags configured on the SEO page.
	document.SEO = dto.SeoPublicSEO{}

	projectionJSON, err := json.Marshal(projection)
	if err != nil {
		return fmt.Errorf("encode TQD administrative unit snapshot: %w", err)
	}
	documentJSON, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("encode administrative unit public document: %w", err)
	}

	snapshotHashBytes := sha256.Sum256(projectionJSON)
	snapshotHash := hex.EncodeToString(snapshotHashBytes[:])
	sourceUpdatedAt := projection.Source.UpdatedAt
	refID := administrativeUnitReferenceID(projection)

	updates := map[string]interface{}{
		"ref_type":           enums.ESEORefTypeAdmUnit,
		"ref_source":         seo_domain.SeoRefSourceTqdAdministrativeUnit,
		"ref_label":          projection.Name,
		"ref_url":            projection.EntityID,
		"source_status":      seo_domain.SeoSourceStatusLinked,
		"ref_missing":        false,
		"ref_snapshot_json":  datatypes.JSON(projectionJSON),
		"ref_hash":           snapshotHash,
		"ref_last_synced_at": now,
		"metadata":           datatypes.JSON(documentJSON),
		"template_key":       "administrative-unit",
		"template_version":   dto.SeoPublicDocumentSchemaVersion,
		"summary":            strings.TrimSpace(projection.PlanningContext.Summary),
	}
	if refID > 0 {
		updates["ref_id"] = refID
	}
	if sourceUpdatedAt != nil {
		updates["source_updated_at"] = *sourceUpdatedAt
	}

	if err := u.seoDomainRepo.Update(ctx, page.ID, updates); err != nil {
		return fmt.Errorf("persist administrative unit source snapshot: %w", err)
	}

	page.RefType = enums.ESEORefTypeAdmUnit
	if refID > 0 {
		page.RefID = &refID
	}
	page.RefSource = seo_domain.SeoRefSourceTqdAdministrativeUnit
	page.RefLabel = projection.Name
	page.RefURL = projection.EntityID
	page.SourceStatus = seo_domain.SeoSourceStatusLinked
	page.RefMissing = false
	page.RefSnapshotJSON = datatypes.JSON(projectionJSON)
	page.RefHash = snapshotHash
	page.RefLastSyncedAt = &now
	page.Metadata = datatypes.JSON(documentJSON)
	page.TemplateKey = "administrative-unit"
	page.TemplateVersion = dto.SeoPublicDocumentSchemaVersion
	page.SourceUpdatedAt = sourceUpdatedAt
	page.Summary = strings.TrimSpace(projection.PlanningContext.Summary)
	return nil
}

func administrativeUnitSourceIdentity(page *seo_domain.SeoDomain) string {
	if page == nil {
		return ""
	}
	for _, value := range []string{page.RefURL, page.CanonicalURL, page.Slug} {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if parsed, err := url.Parse(value); err == nil && strings.TrimSpace(parsed.Path) != "" {
			value = parsed.Path
		}
		value = strings.Trim(value, "/")
		if index := strings.LastIndex(value, "/"); index >= 0 {
			value = value[index+1:]
		}
		if value != "" {
			return value
		}
	}
	return ""
}

func administrativeUnitReferenceID(projection *dto.TqdAdministrativeUnitProjection) uint64 {
	if projection == nil {
		return 0
	}
	value, _ := strconv.ParseUint(strings.TrimSpace(projection.AdministrativeCode), 10, 64)
	return value
}

func buildAdministrativeUnitPublicDocument(
	projection *dto.TqdAdministrativeUnitProjection,
) *dto.SeoPublicDocument {
	if projection == nil {
		return &dto.SeoPublicDocument{SchemaVersion: dto.SeoPublicDocumentSchemaVersion, ModuleID: "administrative-unit"}
	}

	canonicalPath := strings.TrimSpace(projection.CanonicalPath)
	updatedAt := "Chưa cập nhật"
	if projection.Source.UpdatedAt != nil {
		updatedAt = projection.Source.UpdatedAt.Format(time.RFC3339)
	}

	hierarchyNames := make([]string, 0, len(projection.Hierarchy))
	breadcrumbs := []dto.SeoBreadcrumb{{Label: "Trang chủ", Path: "/"}, {Label: "Địa bàn", Path: "/dia-ban"}}
	for _, item := range projection.Hierarchy {
		if strings.TrimSpace(item.Name) == "" {
			continue
		}
		hierarchyNames = append(hierarchyNames, item.Name)
		if isSafePublicPath(item.CanonicalPath) {
			breadcrumbs = append(breadcrumbs, dto.SeoBreadcrumb{Label: item.Name, Path: item.CanonicalPath})
		}
	}
	if len(breadcrumbs) == 2 {
		breadcrumbs = append(breadcrumbs, dto.SeoBreadcrumb{Label: projection.Name, Path: canonicalPath})
	}

	facts := []dto.SeoFact{
		{ID: "unit-type", Label: "Loại đơn vị", Value: firstNonEmptyRender(projection.UnitTypeName, projection.UnitType)},
		{ID: "unit-type-code", Label: "Mã loại đơn vị", Value: projection.UnitType},
		{ID: "administrative-hierarchy", Label: "Phân cấp hành chính", Value: strings.Join(hierarchyNames, " › ")},
		{ID: "administrative-code", Label: "Mã hành chính", Value: projection.AdministrativeCode},
		{ID: "updated-at", Label: "Cập nhật", Value: updatedAt},
	}
	if projection.PlanningContext.PlanningProjectCount != nil {
		facts = append(facts, dto.SeoFact{
			ID: "planning-project-count", Label: "Đồ án liên kết",
			Value: strconv.FormatInt(*projection.PlanningContext.PlanningProjectCount, 10),
		})
	}

	planningItems := []string{firstNonEmptyRender(projection.PlanningContext.Summary, "Thông tin quy hoạch theo địa bàn đang được tổng hợp.")}
	if projection.PlanningContext.PlanningProjectCount != nil {
		planningItems = append(planningItems, fmt.Sprintf("Nguồn hiện tại liên kết %d đồ án quy hoạch với địa bàn này.", *projection.PlanningContext.PlanningProjectCount))
	}
	sections := []dto.SeoSection{
		{ID: "overview", Title: "Tổng quan địa bàn", Blocks: []dto.SeoContentBlock{{Type: "paragraph", Text: "Thông tin hành chính và bối cảnh quy hoạch công khai của " + projection.Name + "."}}},
		{ID: "planning-context", Title: "Bối cảnh quy hoạch", Blocks: []dto.SeoContentBlock{{Type: "list", Items: planningItems}}},
	}

	sourceName := firstNonEmptyRender(projection.Source.AuthorityName, projection.Source.Dataset, "TQD administrative unit projection")
	sections = append(sections, dto.SeoSection{
		ID: "source", Title: "Nguồn dữ liệu",
		Blocks: []dto.SeoContentBlock{{
			Type: "notice", Tone: "info", Title: firstNonEmptyRender(projection.Source.Dataset, "Nguồn dữ liệu"),
			Text: fmt.Sprintf("Nguồn: %s. Cập nhật: %s.", sourceName, updatedAt),
		}},
	})

	mapPath := "/ban-do?focus=" + url.QueryEscape("administrative_unit:"+projection.EntityID)
	relatedLinks := []dto.SeoRelatedLink{
		{ID: "administrative-map", Label: "Xem địa bàn trên bản đồ", Path: mapPath, Description: "Mở vị trí trung tâm và các lớp quy hoạch liên quan trong QHPro."},
	}
	relatedEntities := make([]dto.SeoRelatedEntity, 0, len(projection.PlanningContext.RelatedProjects))
	for _, project := range projection.PlanningContext.RelatedProjects {
		entityID := strings.TrimSpace(project.EntityID)
		publicPath := strings.TrimSpace(project.PublicPath)
		name := strings.TrimSpace(project.Name)
		if entityID == "" || name == "" || !isSafePublicPath(publicPath) {
			continue
		}
		relatedEntities = append(relatedEntities, dto.SeoRelatedEntity{
			ID:           "administrative-project-" + strings.ReplaceAll(entityID, ":", "-"),
			RelationType: "references",
			EntityType:   "planning_project",
			EntityID:     entityID,
			Title:        name,
			Description:  strings.TrimSpace(project.Summary),
			PublicPath:   publicPath,
			MapPath:      strings.TrimSpace(project.MapPath),
		})
		relatedLinks = append(relatedLinks, dto.SeoRelatedLink{
			ID:          "administrative-project-link-" + strings.ReplaceAll(entityID, ":", "-"),
			Label:       name,
			Path:        publicPath,
			Description: strings.TrimSpace(project.Summary),
		})
	}
	document := &dto.SeoPublicDocument{
		SchemaVersion: dto.SeoPublicDocumentSchemaVersion,
		ModuleID:      "administrative-unit",
		ResourceType:  "administrative_area",
		Identity: dto.SeoPublicIdentity{
			ID: projection.EntityID, Slug: projection.Slug, CanonicalPath: canonicalPath,
		},
		Heading: dto.SeoPublicHeading{
			Eyebrow: firstNonEmptyRender(projection.UnitTypeName, "Đơn vị hành chính"),
			Title:   projection.Name,
			Summary: firstNonEmptyRender(projection.PlanningContext.Summary, "Thông tin địa bàn và quy hoạch liên quan."),
		},
		Breadcrumbs:     breadcrumbs,
		Facts:           facts,
		Sections:        sections,
		RelatedLinks:    relatedLinks,
		RelatedEntities: relatedEntities,
		MapTarget:       &dto.SeoMapTarget{Path: mapPath, Label: "Xem địa bàn trên bản đồ"},
		FAQ: &dto.SeoFAQ{GroupKey: "administrative-unit-detail", Items: []dto.SeoFAQItem{
			{ID: "data-scope", Question: "Trang địa bàn này cung cấp thông tin gì?", Answer: "Trang tổng hợp định danh hành chính, phân cấp, vị trí trung tâm, bối cảnh quy hoạch và nguồn dữ liệu hiện có."},
			{ID: "boundary-limit", Question: "Vị trí trên bản đồ có phải ranh giới pháp lý không?", Answer: "Không. Nguồn hiện tại chỉ có điểm trung tâm hành chính; cần đối chiếu polygon địa giới và hồ sơ chính thức khi cần xác định ranh giới."},
		}},
		Source: dto.SeoPublicSource{
			SourceName: sourceName, UpdatedAt: updatedAt,
		},
	}
	if projection.GeometrySummary.Centroid != nil {
		point := projection.GeometrySummary.Centroid
		document.Media = &dto.SeoMediaBundle{
			Primary: &dto.SeoMedia{
				ID: "administrative-unit-" + projection.AdministrativeCode + "-centroid", Role: "overview",
				Alt:    "Vị trí trung tâm " + projection.Name,
				Bounds: [4]float64{point.Longitude, point.Latitude, point.Longitude, point.Latitude},
			},
			Items: []dto.SeoMedia{}, Completeness: firstNonEmptyRender(projection.GeometrySummary.Completeness, "centroid_only"),
			FitMode: "center", Note: "Điểm trung tâm chỉ phục vụ định vị; không đại diện cho polygon địa giới pháp lý.",
		}
	}
	return document
}
