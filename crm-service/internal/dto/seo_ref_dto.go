package dto

import _dto "common/domain/dto"

// SearchSeoRefValuesRequest is the legacy search request.
// It is still used by /v2/crm/admin/seo/ref-values.
type SearchSeoRefValuesRequest struct {
	_dto.Pagable

	RefType       uint32 `json:"refType" form:"refType"`
	ResolverKey   string `json:"resolverKey" form:"resolverKey"`
	SourceService string `json:"sourceService" form:"sourceService"`
}

// SearchSeoSourceValuesRequest is the V2 source-key based search request.
// New APIs should use this instead of asking frontend to compose
// refType + resolverKey + sourceService.
type SearchSeoSourceValuesRequest struct {
	_dto.Pagable

	SourceKey string `json:"sourceKey" form:"sourceKey"`
}

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

	// V2 compatibility field.
	SourceKey string `json:"sourceKey"`
}

type ResolveSeoRefValueRequest struct {
	RefType       uint32 `json:"refType" form:"refType"`
	RefID         uint64 `json:"refId" form:"refId"`
	ResolverKey   string `json:"resolverKey" form:"resolverKey"`
	SourceService string `json:"sourceService" form:"sourceService"`
}

// ResolveSeoSourceValueRequest is the V2 source-key based resolve request.
type ResolveSeoSourceValueRequest struct {
	SourceKey string `json:"sourceKey" form:"sourceKey"`
	RefID     uint64 `json:"refId" form:"refId"`
}

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

	// V2 compatibility fields.
	SourceKey      string `json:"sourceKey"`
	Visibility     string `json:"visibility"`
	Indexable      bool   `json:"indexable"`
	DisabledReason string `json:"disabledReason"`
}
