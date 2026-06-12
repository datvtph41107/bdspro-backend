package mapper

import (
	seo_domain "crm/internal/domain/seo"
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
		Scope:             seo.Scope,
		Classify:          uint32(seo.Classify),
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
		Classify:     uint32(seo.Classify),
		RefType:      uint32(seo.RefType),
		RefId:        seo.RefID,
	}
}

func (m *SeoDomainMapper) SeoDomainListToPb(items []seo_domain.SeoDomain) []*crmpb.SeoDomainResponse {
	responses := make([]*crmpb.SeoDomainResponse, 0, len(items))
	for i := range items {
		responses = append(responses, m.SeoDomainToPb(&items[i]))
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
