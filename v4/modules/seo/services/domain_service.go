package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	_errors "common/errors"
	"userv4/modules/seo/dto"
	"userv4/modules/seo/entity"
	"userv4/modules/seo/repository"
)

// SeoDomainService manages the full SEO domain lifecycle.
type SeoDomainService interface {
	List(ctx context.Context, filter dto.SeoDomainListFilter) (*dto.SeoDomainListResult, error)
	GetByID(ctx context.Context, id uint64) (*entity.SeoDomain, error)
	GetByRef(ctx context.Context, refType uint32, refID uint64) (*entity.SeoDomain, error)
	ListPublicByRef(ctx context.Context, refType uint32, refID uint64) ([]entity.SeoDomain, error)
	GetByCanonicalURL(ctx context.Context, canonicalURL string) (*entity.SeoDomain, error)
	GetBySlug(ctx context.Context, slug string) (*entity.SeoDomain, error)
	Create(ctx context.Context, req dto.CreateSeoDomainRequest) (*entity.SeoDomain, error)
	CreateDraft(ctx context.Context, req dto.CreateSeoDomainRequest) (*entity.SeoDomain, error)
	Update(ctx context.Context, id uint64, req dto.UpdateSeoDomainRequest) (*entity.SeoDomain, error)
	Delete(ctx context.Context, id uint64) error
	Publish(ctx context.Context, id uint64) (*entity.SeoDomain, error)
	Unpublish(ctx context.Context, id uint64) (*entity.SeoDomain, error)
	Archive(ctx context.Context, id uint64) error
	Restore(ctx context.Context, id uint64) error
	Preview(ctx context.Context, id uint64) (string, error)
	PreviewContent(ctx context.Context, req dto.CreateSeoDomainRequest) (string, error)
	TouchByRef(ctx context.Context, refType uint32, refID uint64) (int64, error)
	AssignCategory(ctx context.Context, id, categoryID uint64) (*entity.SeoDomain, error)
	ListRefTypes(ctx context.Context) []dto.SeoSourceItem
	ListSources(ctx context.Context) []dto.SeoSourceItem
	GetSource(ctx context.Context, sourceKey string) (*dto.SeoSourceItem, error)
	SearchRefValues(ctx context.Context, req dto.SearchSeoRefValuesRequest) ([]dto.SeoRefValueItem, int64, error)
	ResolveRefValue(ctx context.Context, req dto.ResolveSeoRefValuePath) (*dto.SeoRefSnapshotItem, error)
	SearchSourceValues(ctx context.Context, req dto.SearchSeoSourceValuesRequest) ([]dto.SeoRefValueItem, int64, error)
	ResolveSourceValue(ctx context.Context, req dto.ResolveSeoSourceValuePath) (*dto.SeoRefSnapshotItem, error)
	UpdateLink(ctx context.Context, id uint64, req dto.UpdateSeoLinkRequest) (*entity.SeoLink, error)
	GetClickStats(ctx context.Context, seoDomainID uint64) (*dto.ClickStatsItem, error)
	ListGenerationLogs(ctx context.Context, filter dto.GenerationLogFilter) ([]dto.GenerationLogItem, int64, error)
	ListForSitemap(ctx context.Context, publicBaseURL string) (string, error)
	// Links
	CreateLink(ctx context.Context, req dto.CreateSeoLinkRequest) (*entity.SeoLink, error)
	ListLinks(ctx context.Context, parentSeoID uint64) ([]entity.SeoLink, error)
	DeleteLink(ctx context.Context, id uint64) error
	// Relatives
	CreateRelative(ctx context.Context, parentSeoID uint64, rel dto.CreateSeoLinkRequest) (*entity.SeoRelative, error)
	ListRelatives(ctx context.Context, parentSeoID uint64) ([]entity.SeoRelative, error)
	DeleteRelative(ctx context.Context, id uint64) error
	// Generate from ref
	GenerateFromRef(ctx context.Context, req dto.GenerateFromRefRequest) (*entity.SeoDomain, bool, error)
	GenerateBatch(ctx context.Context, refType entity.SEORefType, refSource string, ids []uint64) (*dto.GenerateBatchResult, error)
	BatchUpdate(ctx context.Context, req dto.BatchUpdateSeoDomainRequest) error
	CountByRef(ctx context.Context, refType uint32, refID uint64) (int64, error)
	GetTopClicked(ctx context.Context, lastDates uint32) (*dto.GetSeoTopClickedResponse, error)
	GetClickDailyStats(ctx context.Context, lastDates uint32) (*dto.GetSeoClickDailyStatsResponse, error)
	CountRecentByRefType(ctx context.Context, days uint32) (*dto.CountRecentSeoByRefTypeResponse, error)
	// Public document
	GetPublicDocument(ctx context.Context, module, slug string) (*dto.PublicSEODocument, error)
	GetPageHTML(ctx context.Context, path string) (string, error)
}

type LocalSeoDomainService struct {
	repo repository.SeoDomainRepository
}

func NewLocalSeoDomainService(repo repository.SeoDomainRepository) *LocalSeoDomainService {
	return &LocalSeoDomainService{repo: repo}
}

func (s *LocalSeoDomainService) List(ctx context.Context, filter dto.SeoDomainListFilter) (*dto.SeoDomainListResult, error) {
	return s.repo.List(ctx, filter)
}

func (s *LocalSeoDomainService) GetByID(ctx context.Context, id uint64) (*entity.SeoDomain, error) {
	if id == 0 {
		return nil, _errors.NotFoundException("seo domain not found")
	}
	seo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if seo == nil {
		return nil, _errors.NotFoundException("seo domain not found")
	}
	return seo, nil
}

func (s *LocalSeoDomainService) GetByRef(ctx context.Context, refType uint32, refID uint64) (*entity.SeoDomain, error) {
	seo, err := s.repo.GetByRef(ctx, refType, refID)
	if err != nil {
		return nil, err
	}
	if seo == nil {
		return nil, _errors.NotFoundException("seo domain not found for ref")
	}
	return seo, nil
}

func (s *LocalSeoDomainService) ListPublicByRef(ctx context.Context, refType uint32, refID uint64) ([]entity.SeoDomain, error) {
	return s.repo.ListByRef(ctx, refType, refID)
}

func (s *LocalSeoDomainService) CreateDraft(ctx context.Context, req dto.CreateSeoDomainRequest) (*entity.SeoDomain, error) {
	req.PageStatus = entity.PageStatusDraft
	return s.Create(ctx, req)
}

func (s *LocalSeoDomainService) Preview(ctx context.Context, id uint64) (string, error) {
	seo, err := s.GetByID(ctx, id)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(seo.RenderedHTML) != "" {
		return seo.RenderedHTML, nil
	}
	title := seo.Title
	if title == "" {
		title = seo.Slug
	}
	return fmt.Sprintf("<!DOCTYPE html><html><head><title>%s</title></head><body><h1>%s</h1><p>%s</p><div>%s</div></body></html>",
		title, title, seo.Description, seo.Content), nil
}

func (s *LocalSeoDomainService) PreviewContent(_ context.Context, req dto.CreateSeoDomainRequest) (string, error) {
	title := req.Title
	if title == "" {
		title = req.Slug
	}
	return fmt.Sprintf("<!DOCTYPE html><html><head><title>%s</title></head><body><h1>%s</h1><p>%s</p><div>%s</div></body></html>",
		title, title, req.Description, req.Content), nil
}

func (s *LocalSeoDomainService) TouchByRef(ctx context.Context, refType uint32, refID uint64) (int64, error) {
	if refType == 0 || refID == 0 {
		return 0, _errors.BadRequestException("refType and refId are required")
	}
	return s.repo.TouchByRef(ctx, refType, refID)
}

func (s *LocalSeoDomainService) AssignCategory(ctx context.Context, id, categoryID uint64) (*entity.SeoDomain, error) {
	if id == 0 {
		return nil, _errors.BadRequestException("id is required")
	}
	return s.repo.Update(ctx, id, map[string]any{"category_id": categoryID})
}

func (s *LocalSeoDomainService) ListRefTypes(_ context.Context) []dto.SeoSourceItem {
	return mapSourceConfigs(entity.SeoRefTypeConfigs())
}

func (s *LocalSeoDomainService) ListSources(_ context.Context) []dto.SeoSourceItem {
	return mapSourceConfigs(entity.SeoRefTypeConfigs())
}

func (s *LocalSeoDomainService) GetSource(_ context.Context, sourceKey string) (*dto.SeoSourceItem, error) {
	cfg, ok := entity.GetSeoRefTypeConfigByKey(sourceKey)
	if !ok {
		return nil, _errors.NotFoundException("seo source not found: " + strings.TrimSpace(sourceKey))
	}
	item := mapSourceConfig(cfg)
	return &item, nil
}

func (s *LocalSeoDomainService) SearchRefValues(ctx context.Context, req dto.SearchSeoRefValuesRequest) ([]dto.SeoRefValueItem, int64, error) {
	cfg, ok := entity.GetSeoRefTypeConfig(req.RefType, req.ResolverKey, req.SourceService)
	if !ok {
		return nil, 0, _errors.BadRequestException("refType not found")
	}
	cfg = entity.NormalizeSeoRefTypeConfig(cfg)
	if !entity.IsSeoRefTypeSearchable(cfg) {
		return []dto.SeoRefValueItem{}, 0, nil
	}
	return s.searchLocalRefValues(ctx, cfg, req.Text, req.Page, req.Size)
}

func (s *LocalSeoDomainService) ResolveRefValue(ctx context.Context, req dto.ResolveSeoRefValuePath) (*dto.SeoRefSnapshotItem, error) {
	if req.RefID == 0 {
		return nil, _errors.BadRequestException("refId is required")
	}
	cfg, ok := entity.GetSeoRefTypeConfig(req.RefType, req.ResolverKey, req.SourceService)
	if !ok {
		return nil, _errors.BadRequestException("refType not found")
	}
	cfg = entity.NormalizeSeoRefTypeConfig(cfg)
	if !entity.IsSeoRefTypeResolvable(cfg) {
		return nil, _errors.BadRequestException("seo source resolver not available for " + cfg.Key)
	}
	return s.resolveLocalRefValue(ctx, cfg, req.RefID)
}

func (s *LocalSeoDomainService) SearchSourceValues(ctx context.Context, req dto.SearchSeoSourceValuesRequest) ([]dto.SeoRefValueItem, int64, error) {
	cfg, ok := entity.GetSeoRefTypeConfigByKey(req.SourceKey)
	if !ok {
		return nil, 0, _errors.NotFoundException("seo source not found: " + strings.TrimSpace(req.SourceKey))
	}
	cfg = entity.NormalizeSeoRefTypeConfig(cfg)
	if !entity.IsSeoRefTypeSearchable(cfg) {
		return []dto.SeoRefValueItem{}, 0, nil
	}
	return s.searchLocalRefValues(ctx, cfg, req.Text, req.Page, req.Size)
}

func (s *LocalSeoDomainService) ResolveSourceValue(ctx context.Context, req dto.ResolveSeoSourceValuePath) (*dto.SeoRefSnapshotItem, error) {
	if req.RefID == 0 {
		return nil, _errors.BadRequestException("refId is required")
	}
	cfg, ok := entity.GetSeoRefTypeConfigByKey(req.SourceKey)
	if !ok {
		return nil, _errors.NotFoundException("seo source not found: " + strings.TrimSpace(req.SourceKey))
	}
	cfg = entity.NormalizeSeoRefTypeConfig(cfg)
	if !entity.IsSeoRefTypeResolvable(cfg) {
		return nil, _errors.BadRequestException("seo source resolver not available for " + cfg.Key)
	}
	return s.resolveLocalRefValue(ctx, cfg, req.RefID)
}

func (s *LocalSeoDomainService) searchLocalRefValues(ctx context.Context, cfg entity.SeoRefTypeConfig, text string, page, size int) ([]dto.SeoRefValueItem, int64, error) {
	result, err := s.repo.List(ctx, dto.SeoDomainListFilter{
		Page:    page,
		Size:    size,
		Keyword: text,
		RefType: cfg.Value,
	})
	if err != nil {
		return nil, 0, err
	}
	if result == nil {
		return []dto.SeoRefValueItem{}, 0, nil
	}
	out := make([]dto.SeoRefValueItem, 0, len(result.Data))
	for _, row := range result.Data {
		refID := uint64(0)
		if row.RefID != nil {
			refID = *row.RefID
		}
		label := row.RefLabel
		if label == "" {
			label = row.Title
		}
		out = append(out, dto.SeoRefValueItem{
			RefType:           cfg.Value,
			RefTypeKey:        cfg.Key,
			RefID:             refID,
			Label:             label,
			Description:       row.Description,
			ImageURL:          row.ImageURL,
			PublicURL:         row.CanonicalURL,
			AlreadyHasSeoPage: true,
			SeoPageID:         row.ID,
			SourceService:     cfg.SourceService,
			ResolverKey:       cfg.ResolverKey,
			SourceStatus:      row.SourceStatus,
			SourceKey:         cfg.Key,
		})
	}
	return out, result.Total, nil
}

func (s *LocalSeoDomainService) resolveLocalRefValue(ctx context.Context, cfg entity.SeoRefTypeConfig, refID uint64) (*dto.SeoRefSnapshotItem, error) {
	seo, err := s.repo.GetByRef(ctx, cfg.Value, refID)
	if err != nil {
		return nil, err
	}
	if seo == nil {
		return nil, _errors.NotFoundException("seo source value not found")
	}
	label := seo.RefLabel
	if label == "" {
		label = seo.Title
	}
	return &dto.SeoRefSnapshotItem{
		RefType:                  cfg.Value,
		RefTypeKey:               cfg.Key,
		RefID:                    refID,
		Label:                    label,
		Description:              seo.Description,
		ImageURL:                 seo.ImageURL,
		PublicURL:                seo.CanonicalURL,
		SuggestedSlug:            seo.Slug,
		SuggestedTitle:           seo.Title,
		SuggestedMetaTitle:       seo.Title,
		SuggestedMetaDescription: seo.Description,
		SourceService:            cfg.SourceService,
		ResolverKey:              cfg.ResolverKey,
		SourceStatus:             seo.SourceStatus,
		Hash:                     seo.RefHash,
		SourceKey:                cfg.Key,
		Visibility:               "public",
		Indexable:                true,
		DisabledReason:           cfg.DisabledReason,
	}, nil
}

func mapSourceConfigs(cfgs []entity.SeoRefTypeConfig) []dto.SeoSourceItem {
	out := make([]dto.SeoSourceItem, 0, len(cfgs))
	for _, cfg := range cfgs {
		out = append(out, mapSourceConfig(cfg))
	}
	return out
}

func mapSourceConfig(cfg entity.SeoRefTypeConfig) dto.SeoSourceItem {
	cfg = entity.NormalizeSeoRefTypeConfig(cfg)
	return dto.SeoSourceItem{
		Value:            cfg.Value,
		Key:              cfg.Key,
		Label:            cfg.Label,
		Description:      cfg.Description,
		Family:           cfg.Family,
		SourceService:    cfg.SourceService,
		ResolverKey:      cfg.ResolverKey,
		Status:           cfg.Status,
		DisabledReason:   cfg.DisabledReason,
		Capabilities: dto.SeoSourceCapabilities{
			Searchable:       cfg.Capabilities.Searchable,
			Resolvable:       cfg.Capabilities.Resolvable,
			UniquePerRef:     cfg.Capabilities.UniquePerRef,
			SupportsTemplate: cfg.Capabilities.SupportsTemplate,
			SupportsPreview:  cfg.Capabilities.SupportsPreview,
			SupportsGenerate: cfg.Capabilities.SupportsGenerate,
			SupportsHealth:   cfg.Capabilities.SupportsHealth,
		},
		DefaultPageType:  cfg.DefaultPageType,
		AdminURLPattern:  cfg.AdminURLPattern,
		PublicURLPattern: cfg.PublicURLPattern,
		Searchable:       cfg.Searchable,
		UniquePerRef:     cfg.UniquePerRef,
		SupportsTemplate: cfg.SupportsTemplate,
		Resolvable:       cfg.Resolvable,
	}
}

func (s *LocalSeoDomainService) UpdateLink(ctx context.Context, id uint64, req dto.UpdateSeoLinkRequest) (*entity.SeoLink, error) {
	updates := map[string]any{}
	if req.Title != nil {
		updates["anchor_text"] = *req.Title
	}
	if req.Link != nil {
		updates["href"] = *req.Link
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.ChildSeoID != nil {
		updates["child_seo_id"] = *req.ChildSeoID
	}
	if len(updates) == 0 {
		return nil, _errors.BadRequestException("no fields to update")
	}
	return s.repo.UpdateLink(ctx, id, updates)
}

func (s *LocalSeoDomainService) GetByCanonicalURL(ctx context.Context, canonicalURL string) (*entity.SeoDomain, error) {
	seo, err := s.repo.GetByCanonicalURL(ctx, canonicalURL)
	if err != nil {
		return nil, err
	}
	if seo == nil {
		return nil, _errors.NotFoundException("seo page not found")
	}
	return seo, nil
}

func (s *LocalSeoDomainService) GetBySlug(ctx context.Context, slug string) (*entity.SeoDomain, error) {
	seo, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if seo == nil {
		return nil, _errors.NotFoundException("seo page not found")
	}
	return seo, nil
}

func (s *LocalSeoDomainService) Create(ctx context.Context, req dto.CreateSeoDomainRequest) (*entity.SeoDomain, error) {
	if strings.TrimSpace(req.CanonicalURL) == "" {
		return nil, _errors.BadRequestException("canonicalUrl is required")
	}
	seo := &entity.SeoDomain{
		Slug:              req.Slug,
		CanonicalURL:      req.CanonicalURL,
		OriginURL:         req.OriginURL,
		Title:             req.Title,
		Description:       req.Description,
		Summary:           req.Summary,
		Content:           req.Content,
		ImageURL:          req.ImageURL,
		Scope:             firstNonEmpty(req.Scope, entity.ScopePublic),
		PageStatus:        firstNonEmpty(req.PageStatus, entity.PageStatusDraft),
		RefType:           entity.SEORefType(req.RefType),
		RefID:             req.RefID,
		RefSource:         firstNonEmpty(req.RefSource, entity.SourceStatusManual),
		RefLabel:          req.RefLabel,
		RefURL:            req.RefURL,
		CategoryID:        req.CategoryID,
		ProvinceID:        req.ProvinceID,
		VisibleOn:         req.VisibleOn,
		IsSiteMap:         req.IsSiteMap,
		IsIndex:           req.IsIndex,
		TemplateKey:       req.TemplateKey,
		TemplateVersion:   req.TemplateVersion,
		DeepLink:          req.DeepLink,
		Note:              req.Note,
		SitemapPriority:   req.SitemapPriority,
		SitemapChangeFreq: req.SitemapChangeFreq,
		NeedGenerate:      true,
		RenderStatus:      entity.RenderStatusNone,
	}
	if seo.SitemapPriority == 0 {
		seo.SitemapPriority = 0.5
	}
	if seo.VisibleOn == 0 {
		seo.VisibleOn = 30
	}
	return s.repo.Create(ctx, seo)
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func (s *LocalSeoDomainService) Update(ctx context.Context, id uint64, req dto.UpdateSeoDomainRequest) (*entity.SeoDomain, error) {
	if id == 0 {
		return nil, _errors.NotFoundException("seo domain not found")
	}
	updates := map[string]any{}
	applyStr := func(key string, val *string) {
		if val != nil {
			updates[key] = *val
		}
	}
	applyUint64 := func(key string, val *uint64) {
		if val != nil {
			updates[key] = *val
		}
	}
	applyBool := func(key string, val *bool) {
		if val != nil {
			updates[key] = *val
		}
	}
	applyStr("slug", req.Slug)
	applyStr("canonical_url", req.CanonicalURL)
	applyStr("origin_url", req.OriginURL)
	applyStr("title", req.Title)
	applyStr("description", req.Description)
	applyStr("summary", req.Summary)
	applyStr("content", req.Content)
	applyStr("image_url", req.ImageURL)
	applyStr("scope", req.Scope)
	applyStr("page_status", req.PageStatus)
	applyStr("ref_source", req.RefSource)
	applyStr("ref_label", req.RefLabel)
	applyStr("ref_url", req.RefURL)
	applyStr("template_key", req.TemplateKey)
	applyStr("template_version", req.TemplateVersion)
	applyStr("deep_link", req.DeepLink)
	applyStr("note", req.Note)
	if req.RefType != nil {
		updates["ref_type"] = *req.RefType
	}
	applyUint64("ref_id", req.RefID)
	applyUint64("category_id", req.CategoryID)
	applyUint64("province_id", req.ProvinceID)
	applyBool("is_site_map", req.IsSiteMap)
	applyBool("is_index", req.IsIndex)
	if req.VisibleOn != nil {
		updates["visible_on"] = *req.VisibleOn
	}
	if req.SitemapPriority != nil {
		updates["sitemap_priority"] = *req.SitemapPriority
	}
	applyStr("sitemap_change_freq", req.SitemapChangeFreq)
	if len(updates) == 0 {
		return s.GetByID(ctx, id)
	}
	return s.repo.Update(ctx, id, updates)
}

func (s *LocalSeoDomainService) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return _errors.NotFoundException("seo domain not found")
	}
	return s.repo.Delete(ctx, id)
}

func (s *LocalSeoDomainService) Publish(ctx context.Context, id uint64) (*entity.SeoDomain, error) {
	if id == 0 {
		return nil, _errors.NotFoundException("seo domain not found")
	}
	return s.repo.Publish(ctx, id)
}

func (s *LocalSeoDomainService) Unpublish(ctx context.Context, id uint64) (*entity.SeoDomain, error) {
	if id == 0 {
		return nil, _errors.NotFoundException("seo domain not found")
	}
	return s.repo.Unpublish(ctx, id)
}

func (s *LocalSeoDomainService) Archive(ctx context.Context, id uint64) error {
	return s.repo.Archive(ctx, id)
}

func (s *LocalSeoDomainService) Restore(ctx context.Context, id uint64) error {
	return s.repo.Restore(ctx, id)
}

func (s *LocalSeoDomainService) GetClickStats(ctx context.Context, seoDomainID uint64) (*dto.ClickStatsItem, error) {
	return s.repo.GetClickStats(ctx, seoDomainID)
}

func (s *LocalSeoDomainService) ListGenerationLogs(ctx context.Context, filter dto.GenerationLogFilter) ([]dto.GenerationLogItem, int64, error) {
	return s.repo.ListGenerationLogs(ctx, filter)
}

// ListForSitemap builds a simple sitemap XML from eligible pages.
func (s *LocalSeoDomainService) ListForSitemap(ctx context.Context, publicBaseURL string) (string, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(publicBaseURL), "/")
	if baseURL == "" {
		baseURL = "https://qhpro.vn"
	}
	const batchSize = 5000
	var buf strings.Builder
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	buf.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	offset := 0
	for {
		items, err := s.repo.ListForSitemap(ctx, batchSize, offset)
		if err != nil {
			return "", err
		}
		for _, item := range items {
			url := item.CanonicalURL
			if !strings.HasPrefix(url, "http") {
				url = baseURL + url
			}
			lastmod := ""
			if item.SiteMapLastedAt != nil {
				lastmod = item.SiteMapLastedAt.Format(time.RFC3339)
			} else if item.PublishedAt != nil {
				lastmod = item.PublishedAt.Format(time.RFC3339)
			}
			freq := item.SitemapChangeFreq
			if freq == "" {
				freq = "daily"
			}
			priority := fmt.Sprintf("%.1f", item.SitemapPriority)
			if item.SitemapPriority == 0 {
				priority = "0.5"
			}
			buf.WriteString("<url>")
			buf.WriteString("<loc>" + url + "</loc>")
			if lastmod != "" {
				buf.WriteString("<lastmod>" + lastmod + "</lastmod>")
			}
			buf.WriteString("<changefreq>" + freq + "</changefreq>")
			buf.WriteString("<priority>" + priority + "</priority>")
			buf.WriteString("</url>")
		}
		if len(items) < batchSize {
			break
		}
		offset += batchSize
	}
	buf.WriteString("</urlset>")
	return buf.String(), nil
}

// CreateLink adds an internal link from a parent SEO page.
func (s *LocalSeoDomainService) CreateLink(ctx context.Context, req dto.CreateSeoLinkRequest) (*entity.SeoLink, error) {
	if req.ParentSeoID == 0 {
		return nil, _errors.BadRequestException("parentSeoId is required")
	}
	link := &entity.SeoLink{
		ParentSeoID: req.ParentSeoID,
		ChildSeoID:  req.ChildSeoID,
		Href:        req.Href,
		AnchorText:  req.AnchorText,
		Priority:    req.Priority,
	}
	return s.repo.CreateLink(ctx, link)
}

func (s *LocalSeoDomainService) ListLinks(ctx context.Context, parentSeoID uint64) ([]entity.SeoLink, error) {
	return s.repo.ListLinks(ctx, parentSeoID)
}

func (s *LocalSeoDomainService) DeleteLink(ctx context.Context, id uint64) error {
	return s.repo.DeleteLink(ctx, id)
}

func (s *LocalSeoDomainService) CreateRelative(ctx context.Context, parentSeoID uint64, req dto.CreateSeoLinkRequest) (*entity.SeoRelative, error) {
	rel := &entity.SeoRelative{
		ParentSeoID:  parentSeoID,
		ChildSeoID:   req.ChildSeoID,
		RelationType: req.RelationType,
		Priority:     req.Priority,
	}
	return s.repo.CreateRelative(ctx, rel)
}

func (s *LocalSeoDomainService) ListRelatives(ctx context.Context, parentSeoID uint64) ([]entity.SeoRelative, error) {
	return s.repo.ListRelatives(ctx, parentSeoID)
}

func (s *LocalSeoDomainService) DeleteRelative(ctx context.Context, id uint64) error {
	return s.repo.DeleteRelative(ctx, id)
}

// GenerateFromRef creates or updates a SEO page from a reference entity.
func (s *LocalSeoDomainService) GenerateFromRef(ctx context.Context, req dto.GenerateFromRefRequest) (*entity.SeoDomain, bool, error) {
	if req.RefType == 0 {
		return nil, false, _errors.BadRequestException("refType is required")
	}
	if req.RefID == 0 {
		return nil, false, _errors.BadRequestException("refId is required")
	}
	refType := entity.SEORefType(req.RefType)
	slug := fmt.Sprintf("%s-%d", strings.ToLower(refType.Name()), req.RefID)
	slug = strings.ReplaceAll(slug, " ", "-")
	canonicalURL := fmt.Sprintf("/seo/%s/%d", firstNonEmpty(strings.ToLower(req.RefSource), "ref"), req.RefID)
	title := fmt.Sprintf("%s #%d", refType.Name(), req.RefID)
	refID := req.RefID
	seo := &entity.SeoDomain{
		Slug:         slug,
		CanonicalURL: canonicalURL,
		OriginURL:    canonicalURL,
		Title:        title,
		Scope:        entity.ScopePublic,
		PageStatus:   entity.PageStatusDraft,
		RefType:      refType,
		RefID:        &refID,
		RefSource:    firstNonEmpty(req.RefSource, entity.SourceStatusManual),
		RefLabel:     title,
		RefURL:       canonicalURL,
		SourceStatus: entity.SourceStatusLinked,
		NeedGenerate: true,
		RenderStatus: entity.RenderStatusPending,
		VisibleOn:    30,
		IsIndex:      true,
		IsSiteMap:    true,
	}
	return s.repo.UpsertFromRef(ctx, seo)
}

func (s *LocalSeoDomainService) GenerateBatch(ctx context.Context, refType entity.SEORefType, refSource string, ids []uint64) (*dto.GenerateBatchResult, error) {
	result := &dto.GenerateBatchResult{
		Total: int64(len(ids)),
		Items: make([]dto.GenerateBatchItem, 0, len(ids)),
	}
	for _, id := range ids {
		if id == 0 {
			result.Skipped++
			continue
		}
		seo, created, err := s.GenerateFromRef(ctx, dto.GenerateFromRefRequest{
			RefType:   uint32(refType),
			RefID:     id,
			RefSource: refSource,
		})
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("refId=%d: %v", id, err))
			continue
		}
		if seo == nil {
			result.Skipped++
			continue
		}
		if created {
			result.Created++
		}
		result.Items = append(result.Items, dto.GenerateBatchItem{
			RefID:        id,
			SeoID:        seo.ID,
			CanonicalURL: seo.CanonicalURL,
			Slug:         seo.Slug,
			Created:      created,
		})
	}
	return result, nil
}

func (s *LocalSeoDomainService) BatchUpdate(ctx context.Context, req dto.BatchUpdateSeoDomainRequest) error {
	if len(req.IDs) == 0 {
		return _errors.BadRequestException("ids is required")
	}
	if len(req.Fields) == 0 {
		return _errors.BadRequestException("fields is required")
	}
	ids := make([]uint64, 0, len(req.IDs))
	seen := map[uint64]struct{}{}
	for _, id := range req.IDs {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return _errors.BadRequestException("ids is invalid")
	}
	updates := map[string]any{}
	for _, raw := range req.Fields {
		switch strings.ToLower(strings.TrimSpace(raw)) {
		case "issitemap", "is_site_map":
			if req.IsSiteMap == nil {
				return _errors.BadRequestException("isSiteMap is required when fields contains isSiteMap")
			}
			updates["is_site_map"] = *req.IsSiteMap
		default:
			return _errors.BadRequestException("unsupported field: " + raw)
		}
	}
	if len(updates) == 0 {
		return _errors.BadRequestException("fields is invalid")
	}
	return s.repo.UpdateByIDs(ctx, ids, updates)
}

func (s *LocalSeoDomainService) CountByRef(ctx context.Context, refType uint32, refID uint64) (int64, error) {
	return s.repo.CountByRef(ctx, refType, refID)
}

func normalizeLastDates(lastDates uint32) uint32 {
	if lastDates == 0 {
		lastDates = 3
	}
	if lastDates > 30 {
		lastDates = 30
	}
	return lastDates
}

func seoClickStatsFromTime(lastDates uint32) time.Time {
	now := time.Now()
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		loc = now.Location()
	}
	today := time.Date(now.In(loc).Year(), now.In(loc).Month(), now.In(loc).Day(), 0, 0, 0, 0, loc)
	return today.AddDate(0, 0, -int(lastDates)+1)
}

func (s *LocalSeoDomainService) GetTopClicked(ctx context.Context, lastDates uint32) (*dto.GetSeoTopClickedResponse, error) {
	lastDates = normalizeLastDates(lastDates)
	fromTime := seoClickStatsFromTime(lastDates)
	total, err := s.repo.CountClicksFromTime(ctx, fromTime)
	if err != nil {
		return nil, err
	}
	items, err := s.repo.GetTopClicked(ctx, fromTime, 15)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []dto.SeoTopClickedItem{}
	}
	return &dto.GetSeoTopClickedResponse{
		TotalClicks: total,
		Data:        items,
		LastDates:   lastDates,
	}, nil
}

func (s *LocalSeoDomainService) GetClickDailyStats(ctx context.Context, lastDates uint32) (*dto.GetSeoClickDailyStatsResponse, error) {
	lastDates = normalizeLastDates(lastDates)
	fromTime := seoClickStatsFromTime(lastDates)
	rows, err := s.repo.GetClickDailyStats(ctx, fromTime)
	if err != nil {
		return nil, err
	}
	byDate := map[string]uint64{}
	for _, row := range rows {
		byDate[row.Date] = row.Count
	}
	data := make([]dto.SeoClickDailyItem, 0, lastDates)
	for i := 0; i < int(lastDates); i++ {
		day := fromTime.AddDate(0, 0, i)
		date := day.Format("2006-01-02")
		data = append(data, dto.SeoClickDailyItem{Date: date, Count: byDate[date]})
	}
	return &dto.GetSeoClickDailyStatsResponse{Data: data, LastDates: lastDates}, nil
}

func (s *LocalSeoDomainService) CountRecentByRefType(ctx context.Context, days uint32) (*dto.CountRecentSeoByRefTypeResponse, error) {
	if days == 0 {
		days = 5
	}
	if days > 90 {
		days = 90
	}
	items, err := s.repo.CountRecentByRefType(ctx, days)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []dto.SeoRefTypeCountItem{}
	}
	return &dto.CountRecentSeoByRefTypeResponse{Data: items, Days: days}, nil
}

// GetPublicDocument serves a published SEO page by module+slug.
func (s *LocalSeoDomainService) GetPublicDocument(ctx context.Context, module, slug string) (*dto.PublicSEODocument, error) {
	canonicalURL := "/" + strings.TrimLeft(strings.Join([]string{module, slug}, "/"), "/")
	seo, err := s.repo.GetByCanonicalURL(ctx, canonicalURL)
	if err != nil || seo == nil {
		seo, err = s.repo.GetBySlug(ctx, slug)
		if err != nil {
			return nil, err
		}
	}
	if seo == nil || seo.PageStatus != entity.PageStatusPublished {
		return nil, _errors.NotFoundException("document not found")
	}
	return &dto.PublicSEODocument{
		ID:           seo.ID,
		Slug:         seo.Slug,
		Title:        seo.Title,
		Description:  seo.Description,
		CanonicalURL: seo.CanonicalURL,
		Content:      seo.Content,
		RenderedHTML: seo.RenderedHTML,
		ImageURL:     seo.ImageURL,
		PublishedAt:  seo.PublishedAt,
	}, nil
}

func (s *LocalSeoDomainService) GetPageHTML(ctx context.Context, path string) (string, error) {
	path = strings.Trim(strings.TrimSpace(path), "/")
	seo, err := s.repo.GetBySlug(ctx, path)
	if err != nil {
		return "", err
	}
	if seo == nil {
		seo, err = s.repo.GetByCanonicalURL(ctx, "/"+path)
		if err != nil {
			return "", err
		}
	}
	if seo == nil || seo.PageStatus != entity.PageStatusPublished {
		return "", _errors.NotFoundException("page not found")
	}
	if seo.RenderedHTML != "" {
		return seo.RenderedHTML, nil
	}
	// Minimal HTML fallback when static render is not yet available.
	title := seo.Title
	if title == "" {
		title = seo.Slug
	}
	body := seo.Content
	if body == "" {
		body = seo.Description
	}
	return fmt.Sprintf("<!DOCTYPE html><html><head><meta charset=\"utf-8\"><title>%s</title></head><body><h1>%s</h1><div>%s</div></body></html>",
		title, title, body), nil
}
