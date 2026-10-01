package dto

import "time"

// ── SEO Domain admin CRUD ─────────────────────────────────────────────────────

type SeoDomainListFilter struct {
	Page       int    `form:"page"`
	Size       int    `form:"size"`
	Keyword    string `form:"keyword"`
	PageStatus string `form:"pageStatus"`
	Scope      string `form:"scope"`
	Published  *bool  `form:"published"`
	CategoryID uint64 `form:"categoryId"`
	ProvinceID uint64 `form:"provinceId"`
	RefType    uint32 `form:"refType"`
}

type SeoDomainItem struct {
	ID              uint64         `json:"id"`
	Slug            string         `json:"slug"`
	CanonicalURL    string         `json:"canonicalUrl"`
	OriginURL       string         `json:"originUrl,omitempty"`
	Title           string         `json:"title"`
	Description     string         `json:"description,omitempty"`
	Summary         string         `json:"summary,omitempty"`
	Scope           string         `json:"scope"`
	PageStatus      string         `json:"pageStatus"`
	Published       bool           `json:"published"`
	PublishedAt     *time.Time     `json:"publishedAt,omitempty"`
	RefType         uint32         `json:"refType"`
	RefTypeName     string         `json:"refTypeName,omitempty"`
	RefID           *uint64        `json:"refId,omitempty"`
	RefSource       string         `json:"refSource,omitempty"`
	RefLabel        string         `json:"refLabel,omitempty"`
	SourceStatus    string         `json:"sourceStatus,omitempty"`
	CategoryID      *uint64        `json:"categoryId,omitempty"`
	CategoryName    string         `json:"categoryName,omitempty"`
	ProvinceID      *uint64        `json:"provinceId,omitempty"`
	ImageURL        string         `json:"imageUrl,omitempty"`
	TemplateKey     string         `json:"templateKey,omitempty"`
	RenderStatus    string         `json:"renderStatus,omitempty"`
	NeedGenerate    bool           `json:"needGenerate"`
	IsSiteMap       bool           `json:"isSiteMap"`
	IsIndex         bool           `json:"isIndex"`
	VisibleOn       uint32         `json:"visibleOn"`
	DeepLink        string         `json:"deepLink,omitempty"`
	ClickCount      int64          `json:"clickCount,omitempty"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}

type SeoDomainListResult struct {
	Data  []SeoDomainItem `json:"data"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
}

type CreateSeoDomainRequest struct {
	Slug            string         `json:"slug"`
	CanonicalURL    string         `json:"canonicalUrl"`
	OriginURL       string         `json:"originUrl"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	Summary         string         `json:"summary"`
	Content         string         `json:"content"`
	ImageURL        string         `json:"imageUrl"`
	Scope           string         `json:"scope"`
	PageStatus      string         `json:"pageStatus"`
	RefType         uint32         `json:"refType"`
	RefID           *uint64        `json:"refId"`
	RefSource       string         `json:"refSource"`
	RefLabel        string         `json:"refLabel"`
	RefURL          string         `json:"refUrl"`
	CategoryID      *uint64        `json:"categoryId"`
	ProvinceID      *uint64        `json:"provinceId"`
	VisibleOn       uint32         `json:"visibleOn"`
	IsSiteMap       bool           `json:"isSiteMap"`
	IsIndex         bool           `json:"isIndex"`
	TemplateKey     string         `json:"templateKey"`
	TemplateVersion string         `json:"templateVersion"`
	DeepLink        string         `json:"deepLink"`
	Note            string         `json:"note"`
	SitemapPriority float32        `json:"sitemapPriority"`
	SitemapChangeFreq string       `json:"sitemapChangeFreq"`
}

type UpdateSeoDomainRequest struct {
	ID              uint64         `path:"id"`
	Slug            *string        `json:"slug"`
	CanonicalURL    *string        `json:"canonicalUrl"`
	OriginURL       *string        `json:"originUrl"`
	Title           *string        `json:"title"`
	Description     *string        `json:"description"`
	Summary         *string        `json:"summary"`
	Content         *string        `json:"content"`
	ImageURL        *string        `json:"imageUrl"`
	Scope           *string        `json:"scope"`
	PageStatus      *string        `json:"pageStatus"`
	RefType         *uint32        `json:"refType"`
	RefID           *uint64        `json:"refId"`
	RefSource       *string        `json:"refSource"`
	RefLabel        *string        `json:"refLabel"`
	RefURL          *string        `json:"refUrl"`
	CategoryID      *uint64        `json:"categoryId"`
	ProvinceID      *uint64        `json:"provinceId"`
	VisibleOn       *uint32        `json:"visibleOn"`
	IsSiteMap       *bool          `json:"isSiteMap"`
	IsIndex         *bool          `json:"isIndex"`
	TemplateKey     *string        `json:"templateKey"`
	TemplateVersion *string        `json:"templateVersion"`
	DeepLink        *string        `json:"deepLink"`
	Note            *string        `json:"note"`
	SitemapPriority *float32       `json:"sitemapPriority"`
	SitemapChangeFreq *string      `json:"sitemapChangeFreq"`
}

type SeoDomainIDPath struct {
	ID uint64 `path:"id"`
}

type SeoDomainByRefRequest struct {
	RefType uint32 `form:"refType"`
	RefID   uint64 `form:"refId"`
}

// SeoDomainRefsFilter is GET /v2/crm/admin/seo-domains/refs.
type SeoDomainRefsFilter struct {
	Page    int    `form:"page"`
	Size    int    `form:"size"`
	RefType uint32 `form:"refType"`
	RefID   uint64 `form:"refId"`
}

// ── Document / page ───────────────────────────────────────────────────────────

type GetDocumentRequest struct {
	Module string `path:"module"`
	Slug   string `path:"slug"`
}

type GetPageRequest struct {
	Path string `path:"path"`
}

type PublicSEODocument struct {
	ID           uint64         `json:"id"`
	Slug         string         `json:"slug"`
	Title        string         `json:"title"`
	Description  string         `json:"description,omitempty"`
	CanonicalURL string         `json:"canonicalUrl"`
	Content      string         `json:"content,omitempty"`
	RenderedHTML string         `json:"renderedHtml,omitempty"`
	ImageURL     string         `json:"imageUrl,omitempty"`
	PublishedAt  *time.Time     `json:"publishedAt,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

// ── Generation log ────────────────────────────────────────────────────────────

type GenerationLogFilter struct {
	SeoDomainID uint64 `form:"seoDomainId"`
	Page        int    `form:"page"`
	Size        int    `form:"size"`
}

type GenerationLogItem struct {
	ID          uint64     `json:"id"`
	SeoDomainID uint64     `json:"seoDomainId"`
	Status      string     `json:"status"`
	Message     string     `json:"message,omitempty"`
	StartedAt   *time.Time `json:"startedAt,omitempty"`
	FinishedAt  *time.Time `json:"finishedAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// ── Click stats ───────────────────────────────────────────────────────────────

type ClickStatsRequest struct {
	SeoDomainID uint64 `form:"seoDomainId"`
	Days        int    `form:"days"`
}

type ClickStatsItem struct {
	SeoDomainID uint64 `json:"seoDomainId"`
	ClickCount  int64  `json:"clickCount"`
	ViewCount   int64  `json:"viewCount"`
}

// ── Internal links ────────────────────────────────────────────────────────────

type CreateSeoLinkRequest struct {
	ParentSeoID  uint64  `path:"id" json:"parentSeoId"`
	ChildSeoID   *uint64 `json:"childSeoId"`
	Href         string  `json:"href"`
	Link         string  `json:"link"`
	Title        string  `json:"title"`
	AnchorText   string  `json:"anchorText"`
	LinkType     string  `json:"linkType"`
	Priority     int     `json:"priority"`
	RelationType string  `json:"relationType"`
}

type UpdateSeoLinkRequest struct {
	ID         uint64  `path:"id"`
	ChildSeoID *uint64 `json:"childSeoId"`
	Link       *string `json:"link"`
	Title      *string `json:"title"`
	LinkType   *string `json:"linkType"`
	Priority   *int    `json:"priority"`
}

type SeoRefTypeItem struct {
	Code uint32 `json:"code"`
	Name string `json:"name"`
}

// ── Admin SEO source registry / ref-values (CRM path parity) ──────────────────

type SeoSourceCapabilities struct {
	Searchable       bool `json:"searchable"`
	Resolvable       bool `json:"resolvable"`
	UniquePerRef     bool `json:"uniquePerRef"`
	SupportsTemplate bool `json:"supportsTemplate"`
	SupportsPreview  bool `json:"supportsPreview"`
	SupportsGenerate bool `json:"supportsGenerate"`
	SupportsHealth   bool `json:"supportsHealth"`
}

// SeoSourceItem is the admin source catalog entry (also used for legacy ref-types).
type SeoSourceItem struct {
	Value            uint32                `json:"value"`
	Key              string                `json:"key"`
	Label            string                `json:"label"`
	Description      string                `json:"description"`
	Family           string                `json:"family"`
	SourceService    string                `json:"sourceService"`
	ResolverKey      string                `json:"resolverKey"`
	Status           string                `json:"status"`
	DisabledReason   string                `json:"disabledReason,omitempty"`
	Capabilities     SeoSourceCapabilities `json:"capabilities"`
	DefaultPageType  string                `json:"defaultPageType"`
	AdminURLPattern  string                `json:"adminUrlPattern,omitempty"`
	PublicURLPattern string                `json:"publicUrlPattern,omitempty"`
	Searchable       bool                  `json:"searchable"`
	UniquePerRef     bool                  `json:"uniquePerRef"`
	SupportsTemplate bool                  `json:"supportsTemplate"`
	Resolvable       bool                  `json:"resolvable"`
}

type SeoSourceKeyPath struct {
	SourceKey string `path:"sourceKey"`
}

type SearchSeoRefValuesRequest struct {
	RefType       uint32 `form:"refType"`
	Text          string `form:"text"`
	Page          int    `form:"page"`
	Size          int    `form:"size"`
	ResolverKey   string `form:"resolverKey"`
	SourceService string `form:"sourceService"`
}

type ResolveSeoRefValuePath struct {
	RefType       uint32 `path:"refType"`
	RefID         uint64 `path:"refId"`
	ResolverKey   string `form:"resolverKey"`
	SourceService string `form:"sourceService"`
}

type SearchSeoSourceValuesRequest struct {
	SourceKey string `path:"sourceKey"`
	Text      string `form:"text"`
	Page      int    `form:"page"`
	Size      int    `form:"size"`
}

type ResolveSeoSourceValuePath struct {
	SourceKey string `path:"sourceKey"`
	RefID     uint64 `path:"refId"`
}

type SeoRefValueItem struct {
	RefType           uint32 `json:"refType"`
	RefTypeKey        string `json:"refTypeKey"`
	RefID             uint64 `json:"refId"`
	Label             string `json:"label"`
	Description       string `json:"description,omitempty"`
	ImageURL          string `json:"imageUrl,omitempty"`
	AdminURL          string `json:"adminUrl,omitempty"`
	PublicURL         string `json:"publicUrl,omitempty"`
	AlreadyHasSeoPage bool   `json:"alreadyHasSeoPage"`
	SeoPageID         uint64 `json:"seoPageId,omitempty"`
	SourceService     string `json:"sourceService,omitempty"`
	ResolverKey       string `json:"resolverKey,omitempty"`
	SourceStatus      string `json:"sourceStatus,omitempty"`
	SourceKey         string `json:"sourceKey"`
}

type SeoRefSnapshotItem struct {
	RefType                  uint32 `json:"refType"`
	RefTypeKey               string `json:"refTypeKey"`
	RefID                    uint64 `json:"refId"`
	Label                    string `json:"label"`
	Description              string `json:"description,omitempty"`
	ImageURL                 string `json:"imageUrl,omitempty"`
	AdminURL                 string `json:"adminUrl,omitempty"`
	PublicURL                string `json:"publicUrl,omitempty"`
	SuggestedSlug            string `json:"suggestedSlug,omitempty"`
	SuggestedTitle           string `json:"suggestedTitle,omitempty"`
	SuggestedMetaTitle       string `json:"suggestedMetaTitle,omitempty"`
	SuggestedMetaDescription string `json:"suggestedMetaDescription,omitempty"`
	SourceService            string `json:"sourceService,omitempty"`
	ResolverKey              string `json:"resolverKey,omitempty"`
	SourceStatus             string `json:"sourceStatus,omitempty"`
	Hash                     string `json:"hash,omitempty"`
	DataJSON                 string `json:"dataJson,omitempty"`
	SourceKey                string `json:"sourceKey"`
	Visibility               string `json:"visibility,omitempty"`
	Indexable                bool   `json:"indexable"`
	DisabledReason           string `json:"disabledReason,omitempty"`
}

type PublicSeoByRefPath struct {
	RefType uint32 `path:"refType"`
	RefID   uint64 `path:"refId"`
}

type TouchByRefRequest struct {
	RefType uint32 `json:"refType"`
	RefID   uint64 `json:"refId"`
}

type AssignCategoryRequest struct {
	ID         uint64 `path:"id"`
	CategoryID uint64 `json:"categoryId"`
}

type SeoLinkItem struct {
	ID          uint64  `json:"id"`
	ParentSeoID uint64  `json:"parentSeoId"`
	ChildSeoID  *uint64 `json:"childSeoId,omitempty"`
	Href        string  `json:"href"`
	AnchorText  string  `json:"anchorText,omitempty"`
	Priority    int     `json:"priority"`
}

// ── Generate from ref ─────────────────────────────────────────────────────────

type GenerateFromRefRequest struct {
	RefType   uint32 `json:"refType"`
	RefID     uint64 `json:"refId"`
	RefSource string `json:"refSource"`
	Force     bool   `json:"force"`
}

type GenerateIDsRequest struct {
	IDs   []uint64 `json:"ids"`
	Limit uint32   `json:"limit"`
	// Parcel / layer / project / document / adm-unit specific field names (CRM parity).
	ParcelIDs           []uint64 `json:"parcelIds"`
	LayerIDs            []uint64 `json:"layerIds"`
	PlanningProjectIDs  []uint64 `json:"planningProjectIds"`
	DocIDs              []uint64 `json:"docIds"`
	CompositeIDs        []uint64 `json:"compositeIds"`
	ParcelID            uint64   `json:"parcelId"`
}

type GenerateBatchResult struct {
	Items   []GenerateBatchItem `json:"items"`
	Total   int64               `json:"total"`
	Created int64               `json:"created"`
	Skipped int64               `json:"skipped"`
	Failed  int64               `json:"failed"`
	Errors  []string            `json:"errors"`
}

type GenerateBatchItem struct {
	RefID        uint64 `json:"refId"`
	SeoID        uint64 `json:"seoId"`
	CanonicalURL string `json:"canonicalUrl"`
	Slug         string `json:"slug"`
	Created      bool   `json:"created"`
}

// ── Batch update ──────────────────────────────────────────────────────────────

type BatchUpdateSeoDomainRequest struct {
	IDs      []uint64 `json:"ids"`
	Fields   []string `json:"fields"`
	IsSiteMap *bool   `json:"isSiteMap"`
}

// ── Public click / count stats ────────────────────────────────────────────────

type LastDatesRequest struct {
	LastDates uint32 `form:"lastDates"`
	Days      uint32 `form:"days"`
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

type GetSeoTopClickedResponse struct {
	TotalClicks uint64              `json:"totalClicks"`
	Data        []SeoTopClickedItem `json:"data"`
	LastDates   uint32              `json:"lastDates"`
}

type SeoClickDailyItem struct {
	Date  string `json:"date"`
	Count uint64 `json:"count"`
}

type GetSeoClickDailyStatsResponse struct {
	Data      []SeoClickDailyItem `json:"data"`
	LastDates uint32              `json:"lastDates"`
}

type SeoRefTypeCountItem struct {
	RefType     uint32 `json:"refType"`
	RefTypeName string `json:"refTypeName"`
	Total       int64  `json:"total"`
}

type CountRecentSeoByRefTypeResponse struct {
	Data []SeoRefTypeCountItem `json:"data"`
	Days uint32                `json:"days"`
}
