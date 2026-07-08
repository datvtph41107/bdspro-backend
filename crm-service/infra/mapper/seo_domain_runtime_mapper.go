package mapper

import (
	seo_domain "crm/internal/domain/seo"

	crmpb "pb/types/crm"
)

func (m *SeoDomainMapper) SeoDomainToStaticPathPb(seo *seo_domain.SeoDomain) *crmpb.SeoStaticPathResponse {
	if seo == nil {
		return nil
	}
	return &crmpb.SeoStaticPathResponse{
		Id:           seo.ID,
		Slug:         seo.Slug,
		CanonicalUrl: seo.CanonicalURL,
		RefType:      uint32(seo.RefType),
		RefId:        seo.RefID,
		UpdatedAt:    formatTime(seo.UpdatedAt),

		// V2 lifecycle/render fields.
		PageStatus:   seoPageStatusToPb(seo),
		RenderStatus: seoRenderStatusToPb(seo),
		NeedGenerate: seo.NeedGenerate,
	}
}

func (m *SeoDomainMapper) SeoDomainStaticPathListToPb(items []seo_domain.SeoDomain) []*crmpb.SeoStaticPathResponse {
	responses := make([]*crmpb.SeoStaticPathResponse, 0, len(items))
	for i := range items {
		responses = append(responses, m.SeoDomainToStaticPathPb(&items[i]))
	}
	return responses
}