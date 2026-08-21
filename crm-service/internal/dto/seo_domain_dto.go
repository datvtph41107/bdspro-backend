package dto

import (
	_dto "common/domain/dto"
	"crm/internal/enums"
)

type SeoDomainCreateRequest struct {
	Slug              string  `json:"slug"`
	OriginURL         string  `json:"originUrl"`
	CanonicalURL      *string `json:"canonicalUrl"`
	RefType           uint32  `json:"refType"`
	RefID             *uint64 `json:"refId"`
	RefSource         string  `json:"refSource"`
	RefLabel          string  `json:"refLabel"`
	RefURL            string  `json:"refUrl"`
	SourceStatus      string  `json:"sourceStatus"`
	RefMissing        *bool   `json:"refMissing"`
	RefSnapshotJSON   *string `json:"refSnapshotJson"`
	RefHash           string  `json:"refHash"`
	RefLastSyncedAt   *string `json:"refLastSyncedAt"`
	QualityScore      *uint32 `json:"qualityScore"`
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
	ImageURL          string  `json:"imageUrl"`
	Metadata          *string `json:"metadata"`
	Note              string  `json:"note"`

	// V2 lifecycle/render fields.
	PageStatus      string `json:"pageStatus"`
	RenderStatus    string `json:"renderStatus"`
	RenderedHTML    string `json:"renderedHtml"`
	LastRenderError string `json:"lastRenderError"`
	TemplateKey     string `json:"templateKey"`
	TemplateVersion string `json:"templateVersion"`
	CategoryID      *uint64 `json:"categoryId"`
}

type SeoDomainUpdateRequest struct {
	Slug              *string  `json:"slug"`
	OriginURL         *string  `json:"originUrl"`
	CanonicalURL      *string  `json:"canonicalUrl"`
	RefType           *uint32  `json:"refType"`
	RefID             *uint64  `json:"refId"`
	RefSource         *string  `json:"refSource"`
	RefLabel          *string  `json:"refLabel"`
	RefURL            *string  `json:"refUrl"`
	SourceStatus      *string  `json:"sourceStatus"`
	RefMissing        *bool    `json:"refMissing"`
	RefSnapshotJSON   *string  `json:"refSnapshotJson"`
	RefHash           *string  `json:"refHash"`
	RefLastSyncedAt   *string  `json:"refLastSyncedAt"`
	QualityScore      *uint32  `json:"qualityScore"`
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
	ImageURL          *string  `json:"imageUrl"`
	Metadata          *string  `json:"metadata"`
	Note              *string  `json:"note"`

	// V2 lifecycle/render fields.
	PageStatus      *string `json:"pageStatus"`
	RenderStatus    *string `json:"renderStatus"`
	RenderedHTML    *string `json:"renderedHtml"`
	LastRenderError *string `json:"lastRenderError"`
	TemplateKey     *string `json:"templateKey"`
	TemplateVersion *string `json:"templateVersion"`
	CategoryID      *uint64 `json:"categoryId"`
}

type SeoDomainListRequest struct {
	_dto.Pagable
	Keyword      *string `json:"keyword" form:"keyword"`
	RefType      *uint32 `json:"refType" form:"refType"`
	RefID        *uint64 `json:"refId" form:"refId"`
	RefSource    *string `json:"refSource" form:"refSource"`
	SourceStatus *string `json:"sourceStatus" form:"sourceStatus"`
	RefMissing   *bool   `json:"refMissing" form:"refMissing"`
	Scope        *string `json:"scope" form:"scope"`
	Classify     *uint32 `json:"classify" form:"classify"`
	Published    *bool   `json:"published" form:"published"`
	IsSiteMap    *bool   `json:"isSiteMap" form:"isSiteMap"`
	IsIndex      *bool   `json:"isIndex" form:"isIndex"`
	IsRobot      *bool   `json:"isRobot" form:"isRobot"`
	NeedGenerate *bool   `json:"needGenerate" form:"needGenerate"`

	// V2 lifecycle/render filters.
	PageStatus   *string `json:"pageStatus" form:"pageStatus"`
	RenderStatus *string `json:"renderStatus" form:"renderStatus"`
	TemplateKey  *string `json:"templateKey" form:"templateKey"`
}

// CountRecentSeoByRefTypeRequest đếm SEO theo từng refType với updatedAt trong N ngày gần nhất.
type CountRecentSeoByRefTypeRequest struct {
	Days uint32 `json:"days" form:"days"`
}

func (r *CountRecentSeoByRefTypeRequest) Normalize() {
	if r == nil {
		return
	}
	if r.Days == 0 {
		r.Days = 5
	}
}

type SeoRefTypeCountItem struct {
	RefType     enums.ESEORefType `json:"refType"`
	RefTypeName string            `json:"refTypeName"`
	Total       int64             `json:"total"`
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
	RefType   uint32 `json:"refType"`
	RefID     uint64 `json:"refId"`
	RefSource string `json:"refSource"`
}

type InternalSeoRenderDataRequest struct {
	RefType   uint32 `json:"refType"`
	RefID     uint64 `json:"refId"`
	RefSource string `json:"refSource"`
}

type SitemapRequest struct {
	BaseURL *string `json:"baseUrl" form:"baseUrl"`
}

type CreateSeoFromParcelRequest struct {
	ParcelID uint64 `json:"parcelId"`
}

type GenerateSeoFromParcelsRequest struct {
	Limit     uint32   `json:"limit"`
	ParcelIDs []uint64 `json:"parcelIds"`
}

type GenerateSeoFromParcelsResult struct {
	Total   int64    `json:"total"`
	Created int64    `json:"created"`
	Skipped int64    `json:"skipped"`
	Failed  int64    `json:"failed"`
	Errors  []string `json:"errors"`
}

type ParcelSeoSource struct {
	ParcelID  uint64  `json:"parcelId" gorm:"column:parcel_id"`
	AdrSearch string  `json:"adrSearch" gorm:"column:adr_search"`
	SeoID     *uint64 `json:"seoId" gorm:"column:seo_id"`
}

type LayerSeoSource struct {
	LayerID      uint64  `json:"layerId" gorm:"column:layer_id"`
	DisplayName  string  `json:"displayName" gorm:"column:display_name"`
	SeoID        *uint64 `json:"seoId" gorm:"column:seo_id"`
	ThumbnailURL string  `json:"thumbnailUrl" gorm:"column:thumbnail_url"`
	ImageURL     string  `json:"imageUrl" gorm:"column:image_url"`
	Avatar       string  `json:"avatar" gorm:"column:avatar"`
}

type GenerateSeoFromLayersRequest struct {
	LayerIDs []uint64 `json:"layerIds"`
}

type GenerateSeoFromLayersResult struct {
	Total   int64    `json:"total"`
	Created int64    `json:"created"`
	Skipped int64    `json:"skipped"`
	Failed  int64    `json:"failed"`
	Errors  []string `json:"errors"`
}

type PlanningProjectSeoSource struct {
	PlanningProjectID uint64  `json:"planningProjectId" gorm:"column:planning_project_id"`
	Name              string  `json:"name" gorm:"column:name"`
	Slug              string  `json:"slug" gorm:"column:slug"`
	SeoID             *uint64 `json:"seoId" gorm:"column:seo_id"`
}

type GenerateSeoFromPlanningProjectsRequest struct {
	PlanningProjectIDs []uint64 `json:"planningProjectIds"`
}

type GenerateSeoFromPlanningProjectsResult struct {
	Total   int64    `json:"total"`
	Created int64    `json:"created"`
	Skipped int64    `json:"skipped"`
	Failed  int64    `json:"failed"`
	Errors  []string `json:"errors"`
}

type PlanningDocumentSeoSource struct {
	PlanningDocumentID uint64  `json:"planningDocumentId" gorm:"column:planning_document_id"`
	Title              string  `json:"title" gorm:"column:title"`
	Slug               string  `json:"slug" gorm:"column:slug"`
	SeoID              *uint64 `json:"seoId" gorm:"column:seo_id"`
	Filepath           string  `json:"filepath" gorm:"column:filepath"`
	Thumbnail          string  `json:"thumbnail" gorm:"column:thumbnail"`
	Ext                string  `json:"ext" gorm:"column:ext"`
}

type GenerateSeoFromPlanningDocumentsRequest struct {
	DocIDs []uint64 `json:"docIds"`
}

type GenerateSeoFromPlanningDocumentsResult struct {
	Total   int64    `json:"total"`
	Created int64    `json:"created"`
	Skipped int64    `json:"skipped"`
	Failed  int64    `json:"failed"`
	Errors  []string `json:"errors"`
}

type AdministrativeUnitSeoSource struct {
	Identity           uint64  `json:"identity"`
	UnitType           uint32  `json:"unitType"`
	SourceID           uint64  `json:"sourceId"`
	Name               string  `json:"name"`
	Slug               string  `json:"slug"`
	CanonicalPath      string  `json:"canonicalPath"`
	AdministrativeCode string  `json:"administrativeCode"`
	SeoID              *uint64 `json:"seoId"`
	CompositeID        uint64  `json:"compositeId"`
}

type GenerateSeoFromAdministrativeUnitsRequest struct {
	CompositeIDs []uint64 `json:"compositeIds"`
}

type GenerateSeoFromAdministrativeUnitsResult struct {
	Total   int64    `json:"total"`
	Created int64    `json:"created"`
	Skipped int64    `json:"skipped"`
	Failed  int64    `json:"failed"`
	Errors  []string `json:"errors"`
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

type SeoGenerationLogListRequest struct {
	_dto.Pagable
	SeoDomainID *uint64 `json:"seoDomainId" form:"seoDomainId"`
	Status      *string `json:"status" form:"status"`
	TriggerType *string `json:"triggerType" form:"triggerType"`
}
