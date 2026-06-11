package mapper

import (
	seo_domain "crm/internal/domain/seo"
	"time"

	crmpb "pb/types/crm"
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
		Id:              seo.ID,
		OriginUrl:       seo.OriginURL,
		CanonicalUrl:    seo.CanonicalURL,
		Classify:        uint32(seo.Classify),
		SiteMapLastedAt: formatTimePtr(seo.SiteMapLastedAt),
		IsSiteMap:       seo.IsSiteMap,
		IsIndex:         seo.IsIndex,
		IsRobot:         seo.IsRobot,
		Metadata:        seo.Metadata,
		Note:            seo.Note,
		RefType:         uint32(seo.RefType),
		RefId:           seo.RefID,
		InternalLinks:   links,
		Relatives:       relatives,
		CreatedAt:       formatTime(seo.CreatedAt),
		UpdatedAt:       formatTime(seo.UpdatedAt),
	}
}

func (m *SeoDomainMapper) SeoDomainBriefToPb(seo *seo_domain.SeoDomain) *crmpb.SeoDomainBriefResponse {
	if seo == nil {
		return nil
	}
	return &crmpb.SeoDomainBriefResponse{
		Id:           seo.ID,
		OriginUrl:    seo.OriginURL,
		CanonicalUrl: seo.CanonicalURL,
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
		Id:          relative.ID,
		ParentSeoId: relative.ParentSeoID,
		ChildSeoId:  relative.ChildSeoID,
		ParentSeo:   m.SeoDomainBriefToPb(relative.ParentSeo),
		ChildSeo:    m.SeoDomainBriefToPb(relative.ChildSeo),
		CreatedAt:   formatTime(relative.CreatedAt),
		UpdatedAt:   formatTime(relative.UpdatedAt),
	}
}

func (m *SeoDomainMapper) SeoRelativeListToPb(items []seo_domain.SeoRelative) []*crmpb.SeoRelativeResponse {
	responses := make([]*crmpb.SeoRelativeResponse, 0, len(items))
	for i := range items {
		responses = append(responses, m.SeoRelativeToPb(&items[i]))
	}
	return responses
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
