package dto

import _dto "common/domain/dto"

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
	Metadata          *string `json:"metadata"`
	Note              string  `json:"note"`

	// Các trường trạng thái vòng đời và render của trang.
	PageStatus      string `json:"pageStatus"`
	RenderStatus    string `json:"renderStatus"`
	RenderedHTML    string `json:"renderedHtml"`
	LastRenderError string `json:"lastRenderError"`
	TemplateKey     string `json:"templateKey"`
	TemplateVersion string `json:"templateVersion"`
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
	Metadata          *string  `json:"metadata"`
	Note              *string  `json:"note"`

	// Các trường trạng thái vòng đời và render của trang.
	PageStatus      *string `json:"pageStatus"`
	RenderStatus    *string `json:"renderStatus"`
	RenderedHTML    *string `json:"renderedHtml"`
	LastRenderError *string `json:"lastRenderError"`
	TemplateKey     *string `json:"templateKey"`
	TemplateVersion *string `json:"templateVersion"`
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

// SeoSitemapItem là dữ liệu tối thiểu để tạo một URL trong sitemap.
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
	Limit uint32 `json:"limit"`
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

// SeoInternalLinkCreateRequest chứa dữ liệu tạo liên kết nội bộ.
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
