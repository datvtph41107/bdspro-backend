package dto

import _dto "common/domain/dto"

type SeoDomainCreateRequest struct {
	OriginURL       string  `json:"originUrl"`
	CanonicalURL    *string `json:"canonicalUrl"`
	Classify        uint32  `json:"classify"`
	SiteMapLastedAt *string `json:"siteMapLastedAt"`
	IsSiteMap       bool    `json:"isSiteMap"`
	IsIndex         bool    `json:"isIndex"`
	IsRobot         *bool   `json:"isRobot"`
	Metadata        *string `json:"metadata"`
	Note            string  `json:"note"`
	RefType         uint32  `json:"refType"`
	RefID           *uint64 `json:"refId"`
}

type SeoDomainUpdateRequest struct {
	OriginURL       *string `json:"originUrl"`
	CanonicalURL    *string `json:"canonicalUrl"`
	Classify        *uint32 `json:"classify"`
	SiteMapLastedAt *string `json:"siteMapLastedAt"`
	IsSiteMap       *bool   `json:"isSiteMap"`
	IsIndex         *bool   `json:"isIndex"`
	IsRobot         *bool   `json:"isRobot"`
	Metadata        *string `json:"metadata"`
	Note            *string `json:"note"`
	RefType         *uint32 `json:"refType"`
	RefID           *uint64 `json:"refId"`
}

type SeoDomainListRequest struct {
	_dto.Pagable
	Keyword   *string `json:"keyword" form:"keyword"`
	Classify  *uint32 `json:"classify" form:"classify"`
	IsSiteMap *bool   `json:"isSiteMap" form:"isSiteMap"`
	IsIndex   *bool   `json:"isIndex" form:"isIndex"`
	IsRobot   *bool   `json:"isRobot" form:"isRobot"`
	RefType   *uint32 `json:"refType" form:"refType"`
	RefID     *uint64 `json:"refId" form:"refId"`
}

type SeoInternalLinkCreateRequest struct {
	Link        string `json:"link"`
	Title       string `json:"title"`
	ParentSeoID uint64 `json:"parentSeoId"`
	ChildSeoID  uint64 `json:"childSeoId"`
}

type SeoInternalLinkUpdateRequest struct {
	Link        *string `json:"link"`
	Title       *string `json:"title"`
	ParentSeoID *uint64 `json:"parentSeoId"`
	ChildSeoID  *uint64 `json:"childSeoId"`
}

type SeoInternalLinkListRequest struct {
	_dto.Pagable
	ParentSeoID *uint64 `json:"parentSeoId" form:"parentSeoId"`
	ChildSeoID  *uint64 `json:"childSeoId" form:"childSeoId"`
	Keyword     *string `json:"keyword" form:"keyword"`
}

type SeoRelativeCreateRequest struct {
	ParentSeoID uint64 `json:"parentSeoId"`
	ChildSeoID  uint64 `json:"childSeoId"`
}

type SeoRelativeListRequest struct {
	_dto.Pagable
	ParentSeoID *uint64 `json:"parentSeoId" form:"parentSeoId"`
	ChildSeoID  *uint64 `json:"childSeoId" form:"childSeoId"`
}
