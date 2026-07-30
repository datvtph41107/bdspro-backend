package usecase

import (
	_errors "common/errors"
	"context"
	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/interface/provider"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	bdspropb "pb/types/bdspro"
	"strconv"
	"strings"
)

type seoRefResolver struct {
	cfg            seo_domain.SeoRefTypeConfig
	bdsproProvider provider.BdsproProvider
	tqdProvider    provider.TqdProvider
}

type normalizedSeoRef struct {
	RefType        uint32
	RefTypeKey     string
	RefID          uint64
	Label          string
	Description    string
	ImageURL       string
	AdminURL       string
	PublicURL      string
	Source         string
	SourceKey      string
	SourceStatus   string
	Visibility     string
	Indexable      bool
	DisabledReason string
	SuggestedSlug  string
	Data           map[string]interface{}
}

func newSeoRefResolver(cfg seo_domain.SeoRefTypeConfig, bdsproProvider provider.BdsproProvider, tqdProvider provider.TqdProvider) *seoRefResolver {
	return &seoRefResolver{
		cfg:            seo_domain.NormalizeSeoRefTypeConfig(cfg),
		bdsproProvider: bdsproProvider,
		tqdProvider:    tqdProvider,
	}
}

func (r *seoRefResolver) Search(ctx context.Context, req *dto.SearchSeoRefValuesRequest) ([]dto.SeoRefValueResponse, int64, error) {
	if req == nil {
		return nil, 0, _errors.ReturnError(400, "request không hợp lệ")
	}

	if !seo_domain.IsSeoRefTypeSearchable(r.cfg) {
		return []dto.SeoRefValueResponse{}, 0, nil
	}
	if !seo_domain.IsSeoResolverReady(r.cfg.ResolverKey) {
		return nil, 0, resolverNotImplemented(r.cfg.ResolverKey)
	}

	page, size := normalizeRefPaging(req.Page, req.Size)

	switch r.cfg.ResolverKey {
	case seo_domain.SeoRefSourceBdsproProject:
		if r.bdsproProvider == nil {
			return nil, 0, resolverUnavailable(r.cfg.ResolverKey)
		}
		projects, total, err := r.bdsproProvider.SearchProjects(ctx, strings.TrimSpace(req.Text), page, size)
		if err != nil {
			return nil, 0, err
		}
		items := make([]dto.SeoRefValueResponse, 0, len(projects))
		for _, project := range projects {
			if project == nil {
				continue
			}
			items = append(items, r.normalizeProject(project).toValueResponse())
		}
		return items, total, nil

	case seo_domain.SeoRefSourceBdsproRegion:
		if r.bdsproProvider == nil {
			return nil, 0, resolverUnavailable(r.cfg.ResolverKey)
		}
		regions, total, err := r.bdsproProvider.SearchRegions(ctx, strings.TrimSpace(req.Text), page, size)
		if err != nil {
			return nil, 0, err
		}
		items := make([]dto.SeoRefValueResponse, 0, len(regions))
		for _, region := range regions {
			if region == nil {
				continue
			}
			items = append(items, r.normalizeRegion(region).toValueResponse())
		}
		return items, total, nil

	case seo_domain.SeoRefSourceBdsproAreaRegion:
		if r.bdsproProvider == nil {
			return nil, 0, resolverUnavailable(r.cfg.ResolverKey)
		}
		regions, err := r.bdsproProvider.ListAreaRegions(ctx)
		if err != nil {
			return nil, 0, err
		}
		filtered := filterAreaRegions(regions, req.Text)
		total := int64(len(filtered))
		filtered = paginateAreaRegions(filtered, page, size)

		items := make([]dto.SeoRefValueResponse, 0, len(filtered))
		for _, region := range filtered {
			if region == nil {
				continue
			}
			items = append(items, r.normalizeAreaRegion(region).toValueResponse())
		}
		return items, total, nil

	case seo_domain.SeoRefSourceTqdParcel:
		if r.tqdProvider == nil {
			return nil, 0, resolverUnavailable(r.cfg.ResolverKey)
		}
		items, total, err := r.searchParcels(ctx, req, size)
		if err != nil {
			return nil, 0, err
		}
		return items, total, nil

	default:
		return nil, 0, resolverNotImplemented(r.cfg.ResolverKey)
	}
}

func (r *seoRefResolver) Resolve(ctx context.Context, req *dto.ResolveSeoRefValueRequest) (*dto.SeoRefSnapshotResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không hợp lệ")
	}
	if req.RefID == 0 {
		return nil, _errors.ReturnError(400, "refId không hợp lệ")
	}

	if !seo_domain.IsSeoRefTypeResolvable(r.cfg) {
		return nil, resolverNotImplemented(r.cfg.ResolverKey)
	}
	if !seo_domain.IsSeoResolverReady(r.cfg.ResolverKey) {
		return nil, resolverNotImplemented(r.cfg.ResolverKey)
	}

	var item normalizedSeoRef

	switch r.cfg.ResolverKey {
	case seo_domain.SeoRefSourceBdsproProject:
		project, err := r.resolveProject(ctx, req.RefID)
		if err != nil {
			return nil, err
		}
		item = r.normalizeProject(project)

	case seo_domain.SeoRefSourceBdsproRegion:
		if r.bdsproProvider == nil {
			return nil, resolverUnavailable(r.cfg.ResolverKey)
		}
		region, err := r.bdsproProvider.GetRegionByID(ctx, req.RefID)
		if err != nil {
			return nil, err
		}
		if region == nil || region.GetId() == 0 {
			return nil, _errors.ReturnError(404, "source không tồn tại")
		}
		item = r.normalizeRegion(region)

	case seo_domain.SeoRefSourceBdsproAreaRegion:
		if r.bdsproProvider == nil {
			return nil, resolverUnavailable(r.cfg.ResolverKey)
		}
		region, err := r.bdsproProvider.GetAreaRegionByID(ctx, req.RefID)
		if err != nil {
			return nil, err
		}
		if region == nil || region.GetId() == 0 {
			return nil, _errors.ReturnError(404, "source không tồn tại")
		}
		item = r.normalizeAreaRegion(region)

	case seo_domain.SeoRefSourceTqdParcel:
		if r.tqdProvider == nil {
			return nil, resolverUnavailable(r.cfg.ResolverKey)
		}
		source, err := r.tqdProvider.GetParcelSeoSource(ctx, req.RefID)
		if err != nil {
			return nil, err
		}
		if source == nil || source.ParcelID == 0 {
			return nil, _errors.ReturnError(404, "source parcel không tồn tại")
		}
		item = r.normalizeParcel(*source)

	case seo_domain.SeoRefSourceTqdPlanningProject:
		if r.tqdProvider == nil {
			return nil, resolverUnavailable(r.cfg.ResolverKey)
		}

		source, err := r.tqdProvider.GetPlanningProjectProjection(
			ctx,
			strconv.FormatUint(req.RefID, 10),
		)
		if err != nil {
			return nil, err
		}
		if source == nil || source.ID == 0 {
			return nil, _errors.ReturnError(
				404,
				"source đồ án quy hoạch không tồn tại",
			)
		}

		item = r.normalizePlanningProject(*source)

	default:
		return nil, resolverNotImplemented(r.cfg.ResolverKey)
	}

	snapshot := item.toSnapshotResponse()
	return &snapshot, nil
}

func (r *seoRefResolver) resolveProject(ctx context.Context, refID uint64) (*bdspropb.Project, error) {
	if r.bdsproProvider == nil {
		return nil, resolverUnavailable(r.cfg.ResolverKey)
	}

	text := strconv.FormatUint(refID, 10)
	candidates, _, err := r.bdsproProvider.SearchProjects(ctx, text, 0, 50)
	if err != nil {
		return nil, err
	}
	if project := findProjectByID(candidates, refID); project != nil {
		return project, nil
	}

	candidates, _, err = r.bdsproProvider.SearchProjects(ctx, "", 0, 500)
	if err != nil {
		return nil, err
	}
	if project := findProjectByID(candidates, refID); project != nil {
		return project, nil
	}

	return nil, _errors.ReturnError(404, "source project không tồn tại hoặc AdminProjectService chưa có API detail theo id")
}

func (r *seoRefResolver) normalizeProject(project *bdspropb.Project) normalizedSeoRef {
	label := strings.TrimSpace(project.GetName())
	if label == "" {
		label = fmt.Sprintf("Dự án #%d", project.GetId())
	}
	slug := slugify(label)

	data := map[string]interface{}{
		"id":          project.GetId(),
		"name":        project.GetName(),
		"developerId": project.GetDeveloperId(),
		"description": project.GetDescription(),
		"active":      project.GetActive(),
		"createdAt":   project.GetCreatedAt(),
		"updatedAt":   project.GetUpdatedAt(),
		"source":      r.cfg.ResolverKey,
		"sourceKey":   r.cfg.Key,
	}

	return r.newNormalizedSeoRef(normalizedSeoRef{
		RefType:      r.cfg.Value,
		RefTypeKey:   enums.SEORefTypeKey(r.cfg.Value),
		RefID:        project.GetId(),
		Label:        label,
		Description:  strings.TrimSpace(project.GetDescription()),
		AdminURL:     formatRefURL(r.cfg.AdminURLPattern, project.GetId(), slug),
		PublicURL:    formatRefURL(r.cfg.PublicURLPattern, project.GetId(), slug),
		Source:       r.cfg.ResolverKey,
		SourceStatus: seo_domain.SeoSourceStatusLinked,
		Data:         data,
	})
}

func (r *seoRefResolver) normalizeRegion(region *bdspropb.Region) normalizedSeoRef {
	label := strings.TrimSpace(region.GetName())
	if label == "" {
		label = fmt.Sprintf("Đơn vị hành chính #%d", region.GetId())
	}
	slug := slugify(label)

	description := strings.TrimSpace(fmt.Sprintf("%s cấp %d, mã %d", region.GetUnit(), region.GetLevel(), region.GetCode()))
	data := map[string]interface{}{
		"id":        region.GetId(),
		"name":      region.GetName(),
		"code":      region.GetCode(),
		"codeName":  region.GetCodeName(),
		"level":     region.GetLevel(),
		"parentId":  region.GetParentId(),
		"unit":      region.GetUnit(),
		"createdAt": region.GetCreatedAt(),
		"updatedAt": region.GetUpdatedAt(),
		"source":    r.cfg.ResolverKey,
		"sourceKey": r.cfg.Key,
	}

	return r.newNormalizedSeoRef(normalizedSeoRef{
		RefType:      r.cfg.Value,
		RefTypeKey:   enums.SEORefTypeKey(r.cfg.Value),
		RefID:        region.GetId(),
		Label:        label,
		Description:  description,
		AdminURL:     formatRefURL(r.cfg.AdminURLPattern, region.GetId(), slug),
		PublicURL:    formatRefURL(r.cfg.PublicURLPattern, region.GetId(), slug),
		Source:       r.cfg.ResolverKey,
		SourceStatus: seo_domain.SeoSourceStatusLinked,
		Data:         data,
	})
}

func (r *seoRefResolver) normalizeAreaRegion(region *bdspropb.AreaRegion) normalizedSeoRef {
	label := strings.TrimSpace(region.GetName())
	if label == "" {
		label = fmt.Sprintf("Khu vực #%d", region.GetId())
	}
	slug := slugify(label)

	description := "Khu vực phục vụ SEO địa phương"
	if !region.GetActive() {
		description = "Khu vực đang tắt hoạt động"
	}

	data := map[string]interface{}{
		"id":        region.GetId(),
		"name":      region.GetName(),
		"code":      region.GetCode(),
		"active":    region.GetActive(),
		"source":    r.cfg.ResolverKey,
		"sourceKey": r.cfg.Key,
	}

	return r.newNormalizedSeoRef(normalizedSeoRef{
		RefType:      r.cfg.Value,
		RefTypeKey:   enums.SEORefTypeKey(r.cfg.Value),
		RefID:        region.GetId(),
		Label:        label,
		Description:  description,
		AdminURL:     formatRefURL(r.cfg.AdminURLPattern, region.GetId(), slug),
		PublicURL:    formatRefURL(r.cfg.PublicURLPattern, region.GetId(), slug),
		Source:       r.cfg.ResolverKey,
		SourceStatus: seo_domain.SeoSourceStatusLinked,
		Data:         data,
	})
}

func (r *seoRefResolver) searchParcels(ctx context.Context, req *dto.SearchSeoRefValuesRequest, size uint32) ([]dto.SeoRefValueResponse, int64, error) {
	text := strings.TrimSpace(req.Text)
	if text != "" {
		if id, err := strconv.ParseUint(text, 10, 64); err == nil && id > 0 {
			source, err := r.tqdProvider.GetParcelSeoSource(ctx, id)
			if err != nil {
				return nil, 0, err
			}
			if source == nil || source.ParcelID == 0 {
				return []dto.SeoRefValueResponse{}, 0, nil
			}
			item := r.normalizeParcel(*source).toValueResponse()
			return []dto.SeoRefValueResponse{item}, 1, nil
		}
	}

	fetchLimit := size
	if fetchLimit == 0 {
		fetchLimit = 20
	}
	if req.Page > 0 {
		fetchLimit = (req.Page + 1) * fetchLimit
	}

	sources, err := r.tqdProvider.GetParcelSeoSourcesForGenerate(ctx, fetchLimit)
	if err != nil {
		return nil, 0, err
	}

	filtered := filterParcelSeoSources(sources, text)
	total := int64(len(filtered))
	filtered = paginateParcelSeoSources(filtered, req.Page, size)

	items := make([]dto.SeoRefValueResponse, 0, len(filtered))
	for _, source := range filtered {
		items = append(items, r.normalizeParcel(source).toValueResponse())
	}
	return items, total, nil
}

func (r *seoRefResolver) normalizeParcel(source dto.ParcelSeoSource) normalizedSeoRef {
	label := strings.TrimSpace(source.AdrSearch)
	if label == "" {
		label = fmt.Sprintf("Thửa đất #%d", source.ParcelID)
	}
	slug := slugify(label)

	data := map[string]interface{}{
		"parcelId":  source.ParcelID,
		"adrSearch": source.AdrSearch,
		"seoId":     source.SeoID,
		"source":    r.cfg.ResolverKey,
		"sourceKey": r.cfg.Key,
	}

	return r.newNormalizedSeoRef(normalizedSeoRef{
		RefType:      r.cfg.Value,
		RefTypeKey:   enums.SEORefTypeKey(r.cfg.Value),
		RefID:        source.ParcelID,
		Label:        label,
		Description:  fmt.Sprintf("Dữ liệu thửa đất TQD #%d", source.ParcelID),
		AdminURL:     formatRefURL(r.cfg.AdminURLPattern, source.ParcelID, slug),
		PublicURL:    formatRefURL(r.cfg.PublicURLPattern, source.ParcelID, slug),
		Source:       r.cfg.ResolverKey,
		SourceStatus: seo_domain.SeoSourceStatusLinked,
		Data:         data,
	})
}

func (r *seoRefResolver) normalizePlanningProject(
	source dto.TqdPlanningProjectProjection,
) normalizedSeoRef {
	label := strings.TrimSpace(source.Name)
	if label == "" {
		label = fmt.Sprintf("Đồ án quy hoạch #%d", source.ID)
	}

	slug := strings.Trim(strings.TrimSpace(source.Slug), "/")
	if slug == "" {
		slug = slugify(label)
	}

	publicURL := strings.TrimSpace(source.CanonicalPath)
	if publicURL == "" {
		publicURL = formatRefURL(
			r.cfg.PublicURLPattern,
			source.ID,
			slug,
		)
	}

	description := strings.TrimSpace(source.Summary)
	if description == "" {
		parts := make([]string, 0, 3)

		if value := strings.TrimSpace(source.PlanningTypeName); value != "" {
			parts = append(parts, value)
		}
		if value := strings.TrimSpace(source.LegalStatusName); value != "" {
			parts = append(parts, value)
		}
		if value := strings.TrimSpace(source.JurisdictionName); value != "" {
			parts = append(parts, value)
		}

		description = strings.Join(parts, " · ")
	}
	if description == "" {
		description = fmt.Sprintf(
			"Dữ liệu đồ án quy hoạch TQD #%d",
			source.ID,
		)
	}

	data := map[string]interface{}{
		"id":                source.ID,
		"canonicalKey":      source.CanonicalKey,
		"slug":              slug,
		"canonicalPath":     publicURL,
		"code":              source.Code,
		"name":              source.Name,
		"planningType":      source.PlanningType,
		"planningTypeName":  source.PlanningTypeName,
		"planningLevel":     source.PlanningLevel,
		"planningLevelName": source.PlanningLevelName,
		"totalArea":         source.TotalArea,
		"legalStatus":       source.LegalStatus,
		"legalStatusName":   source.LegalStatusName,
		"validityStatus":    source.ValidityStatus,
		"authority":         source.Authority,
		"researchScope":     source.ResearchScope,
		"approvalDate":      source.ApprovalDate,
		"effectiveDate":     source.EffectiveDate,
		"expiryDate":        source.ExpiryDate,
		"jurisdictionId":    source.JurisdictionID,
		"jurisdictionName":  source.JurisdictionName,
		"indicators":        source.Indicators,
		"events":            source.Events,
		"documents":         source.Documents,
		"relatedProjects":   source.RelatedProjects,
		"source":            r.cfg.ResolverKey,
		"sourceKey":         r.cfg.Key,
		"projection":        source,
	}

	return r.newNormalizedSeoRef(normalizedSeoRef{
		RefType:       r.cfg.Value,
		RefTypeKey:    enums.SEORefTypeKey(r.cfg.Value),
		RefID:         source.ID,
		Label:         label,
		Description:   description,
		PublicURL:     publicURL,
		Source:        r.cfg.ResolverKey,
		SourceStatus:  seo_domain.SeoSourceStatusLinked,
		SuggestedSlug: slug,
		Data:          data,
	})
}

func (r *seoRefResolver) newNormalizedSeoRef(item normalizedSeoRef) normalizedSeoRef {
	if item.Source == "" {
		item.Source = r.cfg.ResolverKey
	}
	if item.SourceKey == "" {
		item.SourceKey = r.cfg.Key
	}
	if item.SourceStatus == "" {
		item.SourceStatus = seo_domain.SeoSourceStatusLinked
	}
	if item.Visibility == "" {
		item.Visibility = "public"
	}
	item.Indexable = item.Visibility == "public"
	item.DisabledReason = r.cfg.DisabledReason
	return item
}

func (n normalizedSeoRef) toValueResponse() dto.SeoRefValueResponse {
	return dto.SeoRefValueResponse{
		RefType:       n.RefType,
		RefTypeKey:    n.RefTypeKey,
		RefID:         n.RefID,
		Label:         n.Label,
		Description:   n.Description,
		ImageURL:      n.ImageURL,
		AdminURL:      n.AdminURL,
		PublicURL:     n.PublicURL,
		ResolverKey:   n.Source,
		SourceStatus:  n.SourceStatus,
		SourceService: sourceServiceFromResolver(n.Source),
		SourceKey:     n.SourceKey,
	}
}

func (n normalizedSeoRef) toSnapshotResponse() dto.SeoRefSnapshotResponse {
	dataJSON, hash := buildRefDataJSONAndHash(n.Data)

	title := strings.TrimSpace(n.Label)
	if title == "" {
		title = fmt.Sprintf("%s #%d", n.RefTypeKey, n.RefID)
	}

	description := strings.TrimSpace(n.Description)
	if description == "" {
		description = "Nguồn dữ liệu đã được CRM SEO resolver chuẩn hóa."
	}

	suggestedSlug := strings.Trim(strings.TrimSpace(n.SuggestedSlug), "/")
	if suggestedSlug == "" {
		suggestedSlug = slugify(title)
	}

	return dto.SeoRefSnapshotResponse{
		RefType:                  n.RefType,
		RefTypeKey:               n.RefTypeKey,
		RefID:                    n.RefID,
		Label:                    n.Label,
		Description:              n.Description,
		ImageURL:                 n.ImageURL,
		AdminURL:                 n.AdminURL,
		PublicURL:                n.PublicURL,
		SuggestedSlug:            suggestedSlug,
		SuggestedTitle:           title,
		SuggestedMetaTitle:       title,
		SuggestedMetaDescription: description,
		SourceService:            sourceServiceFromResolver(n.Source),
		ResolverKey:              n.Source,
		SourceStatus:             n.SourceStatus,
		Hash:                     hash,
		DataJSON:                 dataJSON,
		SourceKey:                n.SourceKey,
		Visibility:               n.Visibility,
		Indexable:                n.Indexable,
		DisabledReason:           n.DisabledReason,
	}
}

func buildRefDataJSONAndHash(data map[string]interface{}) (string, string) {
	if data == nil {
		data = map[string]interface{}{}
	}
	payload, err := json.Marshal(data)
	if err != nil {
		payload = []byte("{}")
	}
	sum := sha256.Sum256(payload)
	return string(payload), hex.EncodeToString(sum[:])
}

func normalizeRefPaging(page uint32, size uint32) (uint32, uint32) {
	if size == 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return page, size
}

func resolverUnavailable(resolverKey string) error {
	return _errors.ReturnError(503, "resolver SEO chưa cấu hình client nguồn: "+resolverKey)
}

func resolverNotImplemented(resolverKey string) error {
	return _errors.ReturnError(501, "resolver SEO chưa có adapter production: "+resolverKey)
}

func sourceServiceFromResolver(resolverKey string) string {
	if strings.HasPrefix(resolverKey, "bdspro.") {
		return "bdspro"
	}
	if strings.HasPrefix(resolverKey, "tqd.") {
		return "tqd"
	}
	if strings.HasPrefix(resolverKey, "crm.") || resolverKey == seo_domain.SeoRefSourceManual {
		return "crm"
	}
	return ""
}

func formatRefURL(pattern string, id uint64, slug string) string {
	if strings.TrimSpace(pattern) == "" {
		return ""
	}
	out := strings.ReplaceAll(pattern, "{id}", strconv.FormatUint(id, 10))
	out = strings.ReplaceAll(out, "{slug}", slug)
	return out
}

func findProjectByID(items []*bdspropb.Project, id uint64) *bdspropb.Project {
	for _, item := range items {
		if item != nil && item.GetId() == id {
			return item
		}
	}
	return nil
}

func filterAreaRegions(items []*bdspropb.AreaRegion, text string) []*bdspropb.AreaRegion {
	text = strings.ToLower(strings.TrimSpace(text))
	if text == "" {
		return items
	}

	filtered := make([]*bdspropb.AreaRegion, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		haystack := strings.ToLower(fmt.Sprintf("%d %s %d", item.GetId(), item.GetName(), item.GetCode()))
		if strings.Contains(haystack, text) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func paginateAreaRegions(items []*bdspropb.AreaRegion, page uint32, size uint32) []*bdspropb.AreaRegion {
	start := int(page * size)
	if start >= len(items) {
		return []*bdspropb.AreaRegion{}
	}
	end := start + int(size)
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

func filterParcelSeoSources(items []dto.ParcelSeoSource, text string) []dto.ParcelSeoSource {
	text = strings.ToLower(strings.TrimSpace(text))
	if text == "" {
		return items
	}

	filtered := make([]dto.ParcelSeoSource, 0, len(items))
	for _, item := range items {
		haystack := strings.ToLower(fmt.Sprintf("%d %s", item.ParcelID, item.AdrSearch))
		if strings.Contains(haystack, text) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func paginateParcelSeoSources(items []dto.ParcelSeoSource, page uint32, size uint32) []dto.ParcelSeoSource {
	start := int(page * size)
	if start >= len(items) {
		return []dto.ParcelSeoSource{}
	}
	end := start + int(size)
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}
