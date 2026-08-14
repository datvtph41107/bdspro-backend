package dto

import (
	_dto "common/domain/dto"
	"crm/internal/enums"
)

// GetSeoRefsRequest chứa bộ lọc tham chiếu nguồn của một trang SEO.
type GetSeoRefsRequest struct {
	_dto.Pagable
	RefType enums.ESEORefType `json:"refType"`
	RefID   uint64            `json:"refId"`
}

// SearchSeoRefValuesRequest tìm dữ liệu nguồn theo bộ khóa tương thích cũ.
type SearchSeoRefValuesRequest struct {
	_dto.Pagable

	RefType       uint32 `json:"refType" form:"refType"`
	ResolverKey   string `json:"resolverKey" form:"resolverKey"`
	SourceService string `json:"sourceService" form:"sourceService"`
}

// SearchSeoSourceValuesRequest tìm dữ liệu nguồn theo một sourceKey.
type SearchSeoSourceValuesRequest struct {
	_dto.Pagable

	SourceKey string `json:"sourceKey" form:"sourceKey"`
}

// SeoRefValueResponse là một kết quả tìm kiếm từ dữ liệu nguồn.
type SeoRefValueResponse struct {
	RefType           uint32 `json:"refType"`
	RefTypeKey        string `json:"refTypeKey"`
	RefID             uint64 `json:"refId"`
	Label             string `json:"label"`
	Description       string `json:"description"`
	ImageURL          string `json:"imageUrl"`
	AdminURL          string `json:"adminUrl"`
	PublicURL         string `json:"publicUrl"`
	AlreadyHasSeoPage bool   `json:"alreadyHasSeoPage"`
	SeoPageID         uint64 `json:"seoPageId"`
	SourceService     string `json:"sourceService"`
	ResolverKey       string `json:"resolverKey"`
	SourceStatus      string `json:"sourceStatus"`

	// SourceKey là khóa nguồn dùng cho API hiện tại.
	SourceKey string `json:"sourceKey"`
}

// ResolveSeoRefValueRequest đọc một bản ghi nguồn theo bộ khóa tương thích cũ.
type ResolveSeoRefValueRequest struct {
	RefType       uint32 `json:"refType" form:"refType"`
	RefID         uint64 `json:"refId" form:"refId"`
	ResolverKey   string `json:"resolverKey" form:"resolverKey"`
	SourceService string `json:"sourceService" form:"sourceService"`
}

// ResolveSeoSourceValueRequest đọc một bản ghi nguồn theo sourceKey.
type ResolveSeoSourceValueRequest struct {
	SourceKey string `json:"sourceKey" form:"sourceKey"`
	RefID     uint64 `json:"refId" form:"refId"`
}

// SeoRefSnapshotResponse là bản chụp nguồn dùng để tạo hoặc cập nhật trang SEO.
type SeoRefSnapshotResponse struct {
	RefType                  uint32 `json:"refType"`
	RefTypeKey               string `json:"refTypeKey"`
	RefID                    uint64 `json:"refId"`
	Label                    string `json:"label"`
	Description              string `json:"description"`
	ImageURL                 string `json:"imageUrl"`
	AdminURL                 string `json:"adminUrl"`
	PublicURL                string `json:"publicUrl"`
	SuggestedSlug            string `json:"suggestedSlug"`
	SuggestedTitle           string `json:"suggestedTitle"`
	SuggestedMetaTitle       string `json:"suggestedMetaTitle"`
	SuggestedMetaDescription string `json:"suggestedMetaDescription"`
	SourceService            string `json:"sourceService"`
	ResolverKey              string `json:"resolverKey"`
	SourceStatus             string `json:"sourceStatus"`
	Hash                     string `json:"hash"`
	DataJSON                 string `json:"dataJson"`

	// Các trường kỹ thuật dùng bởi API theo sourceKey.
	SourceKey      string `json:"sourceKey"`
	Visibility     string `json:"visibility"`
	Indexable      bool   `json:"indexable"`
	DisabledReason string `json:"disabledReason"`
}
