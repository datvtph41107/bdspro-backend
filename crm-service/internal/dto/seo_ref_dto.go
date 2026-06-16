package dto

type SearchSeoRefValuesRequest struct {
	RefType       uint32 `json:"refType" form:"refType"`
	Text          string `json:"text" form:"text"`
	Page          uint32 `json:"page" form:"page"`
	Size          uint32 `json:"size" form:"size"`
	ResolverKey   string `json:"resolverKey" form:"resolverKey"`
	SourceService string `json:"sourceService" form:"sourceService"`
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
}

type ResolveSeoRefValueRequest struct {
	RefType       uint32 `json:"refType" form:"refType"`
	RefID         uint64 `json:"refId" form:"refId"`
	ResolverKey   string `json:"resolverKey" form:"resolverKey"`
	SourceService string `json:"sourceService" form:"sourceService"`
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
}
