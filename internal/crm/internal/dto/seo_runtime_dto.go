package dto

import seo_domain "crm/internal/domain/seo"

// InternalSeoRenderDataResponse gom dữ liệu cần thiết để renderer tạo HTML.
type InternalSeoRenderDataResponse struct {
	SeoDomain     *seo_domain.SeoDomain
	Entity        interface{}
	InternalLinks []seo_domain.SeoInternalLink
	Relatives     []seo_domain.SeoRelative
}

// GetInternalSeoRenderDataBySlugRequest tìm dữ liệu render theo slug.
type GetInternalSeoRenderDataBySlugRequest struct {
	Slug string `json:"slug" form:"slug"`
}

// GetInternalSeoStaticPathsRequest phân trang danh sách đường dẫn tĩnh.
type GetInternalSeoStaticPathsRequest struct {
	Page uint32 `json:"page" form:"page"`
	Size uint32 `json:"size" form:"size"`
}

// SeoStaticPathItem là một đường dẫn cần tạo hoặc cập nhật HTML tĩnh.
type SeoStaticPathItem struct {
	ID           uint64  `json:"id"`
	Slug         string  `json:"slug"`
	CanonicalURL string  `json:"canonicalUrl"`
	RefType      uint32  `json:"refType"`
	RefID        *uint64 `json:"refId"`
	UpdatedAt    string  `json:"updatedAt"`
}
