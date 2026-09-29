package mapper

import (
	seo_domain "crm/internal/domain/seo"
	"strings"
	"time"

	crmpb "pb/types/crm"

	"gorm.io/datatypes"
)

type SeoDomainMapper struct{}

// @bind: crm/infra/mapper.SeoDomainMapper
func NewSeoDomainMapper() *SeoDomainMapper {
	return &SeoDomainMapper{}
}

func (m *SeoDomainMapper) SeoDomainToPb(seo *seo_domain.SeoDomain) *crmpb.SeoDomainResponse {
	if seo == nil {
		return nil
	}

	links := make([]*crmpb.SeoInternalLinkResponse, 0, len(seo.InternalLinks))
	for i := range seo.InternalLinks {
		links = append(links, m.SeoInternalLinkToPb(&seo.InternalLinks[i]))
	}

	relatives := make([]*crmpb.SeoRelativeResponse, 0, len(seo.Relatives))
	for i := range seo.Relatives {
		relatives = append(relatives, m.SeoRelativeToPb(&seo.Relatives[i]))
	}

	return &crmpb.SeoDomainResponse{
		Id:                seo.ID,
		Slug:              seo.Slug,
		OriginUrl:         seo.OriginURL,
		CanonicalUrl:      seo.CanonicalURL,
		RefType:           uint32(seo.RefType),
		RefId:             seo.RefID,
		RefSource:         seo.RefSource,
		RefLabel:          seo.RefLabel,
		RefUrl:            seo.RefURL,
		SourceStatus:      seo.SourceStatus,
		RefMissing:        seo.RefMissing,
		RefSnapshotJson:   jsonToOptionalString(seo.RefSnapshotJSON),
		RefHash:           seo.RefHash,
		RefLastSyncedAt:   formatTimePtr(seo.RefLastSyncedAt),
		Scope:             seo.Scope,
		Title:             seo.Title,
		Description:       seo.Description,
		Content:           seo.Content,
		Summary:           seo.Summary,
		Published:         seo.Published,
		PublishedAt:       formatTimePtr(seo.PublishedAt),
		IsSiteMap:         seo.IsSiteMap,
		IsIndex:           seo.IsIndex,
		IsRobot:           seo.IsRobot,
		SiteMapLastedAt:   formatTimePtr(seo.SiteMapLastedAt),
		SitemapPriority:   seo.SitemapPriority,
		SitemapChangeFreq: seo.SitemapChangeFreq,
		NeedGenerate:      seo.NeedGenerate,
		GeneratedAt:       formatTimePtr(seo.GeneratedAt),
		SourceUpdatedAt:   formatTimePtr(seo.SourceUpdatedAt),
		StaticHtmlPath:    seo.StaticHtmlPath,
		StaticHtmlHash:    seo.StaticHtmlHash,
		DeepLink:          seo.DeepLink,
		Metadata:          jsonToOptionalString(seo.Metadata),
		Note:              seo.Note,
		InternalLinks:     links,
		Relatives:         relatives,
		CreatedAt:         formatTime(seo.CreatedAt),
		UpdatedAt:         formatTime(seo.UpdatedAt),

		// V2 lifecycle/render fields.
		PageStatus:      seoPageStatusToPb(seo),
		RenderStatus:    seoRenderStatusToPb(seo),
		RenderedHtml:    seo.RenderedHTML,
		LastRenderError: seo.LastRenderError,
		TemplateKey:     seo.TemplateKey,
		TemplateVersion: seo.TemplateVersion,
	}
}

func (m *SeoDomainMapper) SeoDomainBriefToPb(seo *seo_domain.SeoDomain) *crmpb.SeoDomainBriefResponse {
	if seo == nil {
		return nil
	}
	return &crmpb.SeoDomainBriefResponse{
		Id:           seo.ID,
		Slug:         seo.Slug,
		OriginUrl:    seo.OriginURL,
		CanonicalUrl: seo.CanonicalURL,
		Title:        seo.Title,
		Scope:        seo.Scope,
		RefType:      uint32(seo.RefType),
		RefId:        seo.RefID,
		RefSource:    seo.RefSource,
		RefLabel:     seo.RefLabel,
		SourceStatus: seo.SourceStatus,

		// V2 lifecycle/render fields.
		PageStatus:   seoPageStatusToPb(seo),
		RenderStatus: seoRenderStatusToPb(seo),
		NeedGenerate: seo.NeedGenerate,
	}
}

func (m *SeoDomainMapper) SeoDomainListToPb(items []seo_domain.SeoDomain) []*crmpb.SeoDomainResponse {
	responses := make([]*crmpb.SeoDomainResponse, 0, len(items))
	for i := range items {
		responses = append(responses, m.SeoDomainToPb(&items[i]))
	}
	return responses
}

func (m *SeoDomainMapper) SeoSourceCapabilitiesToPb(cap seo_domain.SeoSourceCapabilities) *crmpb.SeoSourceCapabilitiesResponse {
	return &crmpb.SeoSourceCapabilitiesResponse{
		Searchable:       cap.Searchable,
		Resolvable:       cap.Resolvable,
		UniquePerRef:     cap.UniquePerRef,
		SupportsTemplate: cap.SupportsTemplate,
		SupportsPreview:  cap.SupportsPreview,
		SupportsGenerate: cap.SupportsGenerate,
		SupportsHealth:   cap.SupportsHealth,
	}
}

func (m *SeoDomainMapper) SeoRefTypeConfigToPb(cfg seo_domain.SeoRefTypeConfig) *crmpb.SeoRefTypeResponse {
	cfg = seo_domain.NormalizeSeoRefTypeConfig(cfg)

	return &crmpb.SeoRefTypeResponse{
		Value:            cfg.Value,
		Key:              cfg.Key,
		Label:            cfg.Label,
		Description:      cfg.Description,
		Family:           cfg.Family,
		SourceService:    cfg.SourceService,
		ResolverKey:      cfg.ResolverKey,
		Searchable:       cfg.Searchable,
		UniquePerRef:     cfg.UniquePerRef,
		SupportsTemplate: cfg.SupportsTemplate,
		DefaultPageType:  cfg.DefaultPageType,
		AdminUrlPattern:  cfg.AdminURLPattern,
		PublicUrlPattern: cfg.PublicURLPattern,

		// V2 compatibility fields.
		Status:         cfg.Status,
		Resolvable:     cfg.Resolvable,
		DisabledReason: cfg.DisabledReason,
		Capabilities:   m.SeoSourceCapabilitiesToPb(cfg.Capabilities),
	}
}

func (m *SeoDomainMapper) SeoRefTypeConfigsToPb(items []seo_domain.SeoRefTypeConfig) []*crmpb.SeoRefTypeResponse {
	responses := make([]*crmpb.SeoRefTypeResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, m.SeoRefTypeConfigToPb(item))
	}
	return responses
}

func (m *SeoDomainMapper) SeoSourceConfigToPb(cfg seo_domain.SeoRefTypeConfig) *crmpb.SeoSourceResponse {
	cfg = seo_domain.NormalizeSeoRefTypeConfig(cfg)

	return &crmpb.SeoSourceResponse{
		Value:            cfg.Value,
		Key:              cfg.Key,
		Label:            cfg.Label,
		Description:      cfg.Description,
		Family:           cfg.Family,
		SourceService:    cfg.SourceService,
		ResolverKey:      cfg.ResolverKey,
		Status:           cfg.Status,
		DisabledReason:   cfg.DisabledReason,
		Capabilities:     m.SeoSourceCapabilitiesToPb(cfg.Capabilities),
		DefaultPageType:  cfg.DefaultPageType,
		AdminUrlPattern:  cfg.AdminURLPattern,
		PublicUrlPattern: cfg.PublicURLPattern,

		// Compatibility aliases.
		Searchable:       cfg.Searchable,
		UniquePerRef:     cfg.UniquePerRef,
		SupportsTemplate: cfg.SupportsTemplate,
		Resolvable:       cfg.Resolvable,
	}
}

func (m *SeoDomainMapper) SeoSourceConfigsToPb(items []seo_domain.SeoRefTypeConfig) []*crmpb.SeoSourceResponse {
	responses := make([]*crmpb.SeoSourceResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, m.SeoSourceConfigToPb(item))
	}
	return responses
}

func (m *SeoDomainMapper) SeoGenerationLogToPb(log *seo_domain.SeoGenerationLog) *crmpb.SeoGenerationLogResponse {
	if log == nil {
		return nil
	}

	return &crmpb.SeoGenerationLogResponse{
		Id:             log.ID,
		SeoDomainId:    log.SeoDomainID,
		TriggerType:    log.TriggerType,
		Status:         log.Status,
		StartedAt:      formatTimePtr(log.StartedAt),
		FinishedAt:     formatTimePtr(log.FinishedAt),
		StaticHtmlPath: log.StaticHTMLPath,
		StaticHtmlHash: log.StaticHTMLHash,
		ErrorMessage:   log.ErrorMessage,
		Metadata:       jsonToString(log.Metadata),
		CreatedAt:      formatTime(log.CreatedAt),
		UpdatedAt:      formatTime(log.UpdatedAt),
	}
}

func (m *SeoDomainMapper) SeoGenerationLogListToPb(items []seo_domain.SeoGenerationLog) []*crmpb.SeoGenerationLogResponse {
	responses := make([]*crmpb.SeoGenerationLogResponse, 0, len(items))
	for i := range items {
		responses = append(responses, m.SeoGenerationLogToPb(&items[i]))
	}
	return responses
}

func (m *SeoDomainMapper) SeoInternalLinkToPb(link *seo_domain.SeoInternalLink) *crmpb.SeoInternalLinkResponse {
	if link == nil {
		return nil
	}
	return &crmpb.SeoInternalLinkResponse{
		Id:          link.ID,
		Link:        link.Link,
		Title:       link.Title,
		LinkType:    link.LinkType,
		Priority:    link.Priority,
		ParentSeoId: link.ParentSeoID,
		ChildSeoId:  link.ChildSeoID,
		ParentSeo:   m.SeoDomainBriefToPb(link.ParentSeo),
		ChildSeo:    m.SeoDomainBriefToPb(link.ChildSeo),
		CreatedAt:   formatTime(link.CreatedAt),
		UpdatedAt:   formatTime(link.UpdatedAt),
	}
}

func (m *SeoDomainMapper) SeoInternalLinkListToPb(items []seo_domain.SeoInternalLink) []*crmpb.SeoInternalLinkResponse {
	responses := make([]*crmpb.SeoInternalLinkResponse, 0, len(items))
	for i := range items {
		responses = append(responses, m.SeoInternalLinkToPb(&items[i]))
	}
	return responses
}

func (m *SeoDomainMapper) SeoRelativeToPb(relative *seo_domain.SeoRelative) *crmpb.SeoRelativeResponse {
	if relative == nil {
		return nil
	}
	return &crmpb.SeoRelativeResponse{
		Id:           relative.ID,
		ParentSeoId:  relative.ParentSeoID,
		ChildSeoId:   relative.ChildSeoID,
		RelationType: relative.RelationType,
		Priority:     relative.Priority,
		ParentSeo:    m.SeoDomainBriefToPb(relative.ParentSeo),
		ChildSeo:     m.SeoDomainBriefToPb(relative.ChildSeo),
		CreatedAt:    formatTime(relative.CreatedAt),
		UpdatedAt:    formatTime(relative.UpdatedAt),
	}
}

func (m *SeoDomainMapper) SeoRelativeListToPb(items []seo_domain.SeoRelative) []*crmpb.SeoRelativeResponse {
	responses := make([]*crmpb.SeoRelativeResponse, 0, len(items))
	for i := range items {
		responses = append(responses, m.SeoRelativeToPb(&items[i]))
	}
	return responses
}

func jsonToString(value datatypes.JSON) string {
	if len(value) == 0 {
		return ""
	}

	text := strings.TrimSpace(string(value))
	if text == "" || text == "null" {
		return ""
	}

	return text
}

func jsonToOptionalString(data datatypes.JSON) *string {
	if len(data) == 0 {
		return nil
	}
	s := string(data)
	if s == "" || s == "null" {
		return nil
	}
	return &s
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func formatTimePtr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func seoPageStatusToPb(seo *seo_domain.SeoDomain) string {
	if seo == nil {
		return ""
	}

	switch seo.PageStatus {
	case seo_domain.SeoPageStatusDraft,
		seo_domain.SeoPageStatusPublished,
		seo_domain.SeoPageStatusArchived:
		return seo.PageStatus
	}

	if seo.DeletedAt.Valid {
		return seo_domain.SeoPageStatusArchived
	}
	if seo.Published {
		return seo_domain.SeoPageStatusPublished
	}
	return seo_domain.SeoPageStatusDraft
}

func seoRenderStatusToPb(seo *seo_domain.SeoDomain) string {
	if seo == nil {
		return ""
	}

	switch seo.RenderStatus {
	case seo_domain.SeoRenderStatusNone,
		seo_domain.SeoRenderStatusPending,
		seo_domain.SeoRenderStatusRendering,
		seo_domain.SeoRenderStatusSuccess,
		seo_domain.SeoRenderStatusFailed:
		return seo.RenderStatus
	}

	if seo.NeedGenerate {
		return seo_domain.SeoRenderStatusPending
	}
	if seo.StaticHtmlHash != "" || seo.GeneratedAt != nil {
		return seo_domain.SeoRenderStatusSuccess
	}
	return seo_domain.SeoRenderStatusNone
}
