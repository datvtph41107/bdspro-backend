package dto

import seo_domain "crm/internal/domain/seo"

type InternalSeoRenderDataResponse struct {
	SeoDomain     *seo_domain.SeoDomain
	Entity        interface{}
	InternalLinks []seo_domain.SeoInternalLink
	Relatives     []seo_domain.SeoRelative
}

type GetInternalSeoRenderDataBySlugRequest struct {
	Slug string `json:"slug" form:"slug"`
}

type GetInternalSeoStaticPathsRequest struct {
	Page uint32 `json:"page" form:"page"`
	Size uint32 `json:"size" form:"size"`
}

type SeoStaticPathItem struct {
	ID           uint64  `json:"id"`
	Slug         string  `json:"slug"`
	CanonicalURL string  `json:"canonicalUrl"`
	RefType      uint32  `json:"refType"`
	RefID        *uint64 `json:"refId"`
	UpdatedAt    string  `json:"updatedAt"`
}
