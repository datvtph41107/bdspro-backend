package dto

import (
	_dto "common/domain/dto"
	"crm/internal/enums"
)

// SeoDomainCreateRequest chứa dữ liệu tạo mới một trang SEO.
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
	Scope             string  `json:"scope"`
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

	// Các trường trạng thái vòng đời và render của trang.
	PageStatus      string  `json:"pageStatus"`
	RenderStatus    string  `json:"renderStatus"`
	RenderedHTML    string  `json:"renderedHtml"`
	LastRenderError string  `json:"lastRenderError"`
	TemplateKey     string  `json:"templateKey"`
	TemplateVersion string  `json:"templateVersion"`
	CategoryID      *uint64 `json:"categoryId"`
}

// SeoDomainUpdateRequest chứa các trường được phép cập nhật một phần.
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
	Scope             *string  `json:"scope"`
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

	// Các trường trạng thái vòng đời và render của trang.
	PageStatus      *string `json:"pageStatus"`
	RenderStatus    *string `json:"renderStatus"`
	RenderedHTML    *string `json:"renderedHtml"`
	LastRenderError *string `json:"lastRenderError"`
	TemplateKey     *string `json:"templateKey"`
	TemplateVersion *string `json:"templateVersion"`
	CategoryID      *uint64 `json:"categoryId"`
}

type BatchUpdateSeoDomainRequest struct {
	IDs       []uint64 `json:"ids"`
	Fields    []string `json:"fields"`
	IsSiteMap *bool    `json:"isSiteMap"`
}

// SeoDomainListRequest chứa bộ lọc danh sách trang trong Admin.
type SeoDomainListRequest struct {
	_dto.Pagable
	Keyword      *string `json:"keyword" form:"keyword"`
	RefType      *uint32 `json:"refType" form:"refType"`
	RefID        *uint64 `json:"refId" form:"refId"`
	RefSource    *string `json:"refSource" form:"refSource"`
	SourceStatus *string `json:"sourceStatus" form:"sourceStatus"`
	RefMissing   *bool   `json:"refMissing" form:"refMissing"`
	Scope        *string `json:"scope" form:"scope"`
	Published    *bool   `json:"published" form:"published"`
	IsSiteMap    *bool   `json:"isSiteMap" form:"isSiteMap"`
	IsIndex      *bool   `json:"isIndex" form:"isIndex"`
	IsRobot      *bool   `json:"isRobot" form:"isRobot"`
	NeedGenerate *bool   `json:"needGenerate" form:"needGenerate"`

	// Các bộ lọc theo trạng thái vòng đời và render.
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

type SeoClickSource struct {
	RefType uint32 `json:"refType"`
	RefID   uint64 `json:"refId"`
}

type SeoTopClickedItem struct {
	ID         uint64         `json:"id"`
	Name       string         `json:"name"`
	Source     SeoClickSource `json:"source"`
	ClickCount uint64         `json:"clickCount"`
}

type GetSeoTopClickedRequest struct {
	LastDates uint32 `json:"lastDates" form:"lastDates"`
}

func NormalizeLastDates(lastDates uint32) uint32 {
	if lastDates == 0 {
		return 3
	}
	if lastDates > 30 {
		return 30
	}
	return lastDates
}

func (r *GetSeoTopClickedRequest) Normalize() {
	if r == nil {
		return
	}
	r.LastDates = NormalizeLastDates(r.LastDates)
}

type GetSeoTopClickedResponse struct {
	TotalClicks uint64              `json:"totalClicks"`
	Data        []SeoTopClickedItem `json:"data"`
	LastDates   uint32              `json:"lastDates"`
}

type GetSeoClickDailyStatsRequest struct {
	LastDates uint32 `json:"lastDates" form:"lastDates"`
}

func (r *GetSeoClickDailyStatsRequest) Normalize() {
	if r == nil {
		return
	}
	r.LastDates = NormalizeLastDates(r.LastDates)
}

type SeoClickDailyItem struct {
	Date  string `json:"date"`
	Count uint64 `json:"count"`
}

type GetSeoClickDailyStatsResponse struct {
	Data      []SeoClickDailyItem `json:"data"`
	LastDates uint32              `json:"lastDates"`
}

type SeoSitemapItem struct {
	CanonicalURL      string
	UpdatedAt         string
	SiteMapLastedAt   string
	SitemapPriority   float32
	SitemapChangeFreq string
}

// SeoSitemapRequest lọc sitemap theo phạm vi sử dụng.
type SeoSitemapRequest struct {
	Scope *string `json:"scope" form:"scope"`
}

// TouchSeoDomainByRefRequest đánh dấu trang cần đồng bộ lại từ dữ liệu nguồn.
type TouchSeoDomainByRefRequest struct {
	RefType   uint32 `json:"refType"`
	RefID     uint64 `json:"refId"`
	RefSource string `json:"refSource"`
}

// InternalSeoRenderDataRequest định danh trang bằng bản ghi nguồn.
type InternalSeoRenderDataRequest struct {
	RefType   uint32 `json:"refType"`
	RefID     uint64 `json:"refId"`
	RefSource string `json:"refSource"`
}

// SitemapRequest cho phép truyền URL gốc khi tạo sitemap.
type SitemapRequest struct {
	BaseURL *string `json:"baseUrl" form:"baseUrl"`
}

// CreateSeoFromParcelRequest tạo một trang SEO từ thửa đất.
type CreateSeoFromParcelRequest struct {
	ParcelID uint64 `json:"parcelId"`
}

// GenerateSeoFromParcelsRequest giới hạn số thửa đất xử lý trong một lượt.
type GenerateSeoFromParcelsRequest struct {
	Limit     uint32   `json:"limit"`
	ParcelIDs []uint64 `json:"parcelIds"`
}

// GenerateSeoFromParcelsResult tổng hợp kết quả tạo trang theo lô.
type GenerateSeoFromParcelsResult struct {
	Total   int64    `json:"total"`
	Created int64    `json:"created"`
	Skipped int64    `json:"skipped"`
	Failed  int64    `json:"failed"`
	Errors  []string `json:"errors"`
}

// ParcelSeoSource là dữ liệu thửa đất tối thiểu dùng để tạo trang.
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
	ProvinceID   *uint64 `json:"provinceId" gorm:"column:province_id"`
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
	ProvinceID        *uint64 `json:"provinceId" gorm:"column:province_id"`
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
	ProvinceID         *uint64 `json:"provinceId" gorm:"column:province_id"`
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
	ProvinceID         *uint64 `json:"provinceId,omitempty"`
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

// SeoInternalLinkUpdateRequest chứa dữ liệu sửa liên kết nội bộ.
type SeoInternalLinkUpdateRequest struct {
	Link        *string `json:"link"`
	Title       *string `json:"title"`
	LinkType    *string `json:"linkType"`
	Priority    *uint32 `json:"priority"`
	ParentSeoID *uint64 `json:"parentSeoId"`
	ChildSeoID  *uint64 `json:"childSeoId"`
}

// SeoInternalLinkListRequest chứa bộ lọc danh sách liên kết nội bộ.
type SeoInternalLinkListRequest struct {
	_dto.Pagable
	ParentSeoID *uint64 `json:"parentSeoId" form:"parentSeoId"`
	ChildSeoID  *uint64 `json:"childSeoId" form:"childSeoId"`
	LinkType    *string `json:"linkType" form:"linkType"`
	Keyword     *string `json:"keyword" form:"keyword"`
}

// SeoRelativeCreateRequest chứa dữ liệu tạo quan hệ giữa hai trang.
type SeoRelativeCreateRequest struct {
	ParentSeoID  uint64 `json:"parentSeoId"`
	ChildSeoID   uint64 `json:"childSeoId"`
	RelationType string `json:"relationType"`
	Priority     uint32 `json:"priority"`
}

// SeoRelativeListRequest chứa bộ lọc danh sách quan hệ trang.
type SeoRelativeListRequest struct {
	_dto.Pagable
	ParentSeoID  *uint64 `json:"parentSeoId" form:"parentSeoId"`
	ChildSeoID   *uint64 `json:"childSeoId" form:"childSeoId"`
	RelationType *string `json:"relationType" form:"relationType"`
}

// SeoGenerationLogListRequest chứa bộ lọc lịch sử tạo HTML.
type SeoGenerationLogListRequest struct {
	_dto.Pagable
	SeoDomainID *uint64 `json:"seoDomainId" form:"seoDomainId"`
	Status      *string `json:"status" form:"status"`
	TriggerType *string `json:"triggerType" form:"triggerType"`
}
