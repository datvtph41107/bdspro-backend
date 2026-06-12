package dto

import _dto "common/domain/dto"

type SeoDomainCreateRequest struct {
	Slug              string  `json:"slug"`
	OriginURL         string  `json:"originUrl"`
	CanonicalURL      *string `json:"canonicalUrl"`
	RefType           uint32  `json:"refType"`
	RefID             *uint64 `json:"refId"`
	Scope             string  `json:"scope"`
	Classify          uint32  `json:"classify"`
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	Content           string  `json:"content"`
	Summary           string  `json:"summary"`
	Published         bool    `json:"published"`
	PublishedAt       *string `json:"publishedAt"`
	IsSiteMap         bool    `json:"isSiteMap"`
	IsIndex           bool    `json:"isIndex"`
	IsRobot           *bool   `json:"isRobot"`
	SiteMapLastedAt   *string `json:"siteMapLastedAt"`
	SitemapPriority   float32 `json:"sitemapPriority"`
	SitemapChangeFreq string  `json:"sitemapChangeFreq"`
	NeedGenerate      *bool   `json:"needGenerate"`
	GeneratedAt       *string `json:"generatedAt"`
	SourceUpdatedAt   *string `json:"sourceUpdatedAt"`
	StaticHtmlPath    string  `json:"staticHtmlPath"`
	StaticHtmlHash    string  `json:"staticHtmlHash"`
	DeepLink          string  `json:"deepLink"`
	Metadata          *string `json:"metadata"`
	Note              string  `json:"note"`
}

type SeoDomainUpdateRequest struct {
	Slug              *string  `json:"slug"`
	OriginURL         *string  `json:"originUrl"`
	CanonicalURL      *string  `json:"canonicalUrl"`
	RefType           *uint32  `json:"refType"`
	RefID             *uint64  `json:"refId"`
	Scope             *string  `json:"scope"`
	Classify          *uint32  `json:"classify"`
	Title             *string  `json:"title"`
	Description       *string  `json:"description"`
	Content           *string  `json:"content"`
	Summary           *string  `json:"summary"`
	Published         *bool    `json:"published"`
	PublishedAt       *string  `json:"publishedAt"`
	IsSiteMap         *bool    `json:"isSiteMap"`
	IsIndex           *bool    `json:"isIndex"`
	IsRobot           *bool    `json:"isRobot"`
	SiteMapLastedAt   *string  `json:"siteMapLastedAt"`
	SitemapPriority   *float32 `json:"sitemapPriority"`
	SitemapChangeFreq *string  `json:"sitemapChangeFreq"`
	NeedGenerate      *bool    `json:"needGenerate"`
	GeneratedAt       *string  `json:"generatedAt"`
	SourceUpdatedAt   *string  `json:"sourceUpdatedAt"`
	StaticHtmlPath    *string  `json:"staticHtmlPath"`
	StaticHtmlHash    *string  `json:"staticHtmlHash"`
	DeepLink          *string  `json:"deepLink"`
	Metadata          *string  `json:"metadata"`
	Note              *string  `json:"note"`
}

type SeoDomainListRequest struct {
	_dto.Pagable
	Keyword      *string `json:"keyword" form:"keyword"`
	RefType      *uint32 `json:"refType" form:"refType"`
	RefID        *uint64 `json:"refId" form:"refId"`
	Scope        *string `json:"scope" form:"scope"`
	Classify     *uint32 `json:"classify" form:"classify"`
	Published    *bool   `json:"published" form:"published"`
	IsSiteMap    *bool   `json:"isSiteMap" form:"isSiteMap"`
	IsIndex      *bool   `json:"isIndex" form:"isIndex"`
	IsRobot      *bool   `json:"isRobot" form:"isRobot"`
	NeedGenerate *bool   `json:"needGenerate" form:"needGenerate"`
}

type SeoSitemapItem struct {
	CanonicalURL      string
	UpdatedAt         string
	SiteMapLastedAt   string
	SitemapPriority   float32
	SitemapChangeFreq string
}

type SeoSitemapRequest struct {
	Scope *string `json:"scope" form:"scope"`
}

type TouchSeoDomainByRefRequest struct {
	RefType uint32 `json:"refType"`
	RefID   uint64 `json:"refId"`
}

type InternalSeoRenderDataRequest struct {
	RefType uint32 `json:"refType"`
	RefID   uint64 `json:"refId"`
}

type SeoInternalLinkCreateRequest struct {
	Link        string `json:"link"`
	Title       string `json:"title"`
	LinkType    string `json:"linkType"`
	Priority    uint32 `json:"priority"`
	ParentSeoID uint64 `json:"parentSeoId"`
	ChildSeoID  uint64 `json:"childSeoId"`
}

type SeoInternalLinkUpdateRequest struct {
	Link        *string `json:"link"`
	Title       *string `json:"title"`
	LinkType    *string `json:"linkType"`
	Priority    *uint32 `json:"priority"`
	ParentSeoID *uint64 `json:"parentSeoId"`
	ChildSeoID  *uint64 `json:"childSeoId"`
}

type SeoInternalLinkListRequest struct {
	_dto.Pagable
	ParentSeoID *uint64 `json:"parentSeoId" form:"parentSeoId"`
	ChildSeoID  *uint64 `json:"childSeoId" form:"childSeoId"`
	LinkType    *string `json:"linkType" form:"linkType"`
	Keyword     *string `json:"keyword" form:"keyword"`
}

type SeoRelativeCreateRequest struct {
	ParentSeoID  uint64 `json:"parentSeoId"`
	ChildSeoID   uint64 `json:"childSeoId"`
	RelationType string `json:"relationType"`
	Priority     uint32 `json:"priority"`
}

type SeoRelativeListRequest struct {
	_dto.Pagable
	ParentSeoID  *uint64 `json:"parentSeoId" form:"parentSeoId"`
	ChildSeoID   *uint64 `json:"childSeoId" form:"childSeoId"`
	RelationType *string `json:"relationType" form:"relationType"`
}
