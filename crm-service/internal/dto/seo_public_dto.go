package dto

// SeoPublicPageRequest tìm trang công khai theo slug, canonical hoặc path.
type SeoPublicPageRequest struct {
	Slug         string `json:"slug" form:"slug"`
	CanonicalURL string `json:"canonicalUrl" form:"canonicalUrl"`
	Path         string `json:"path" form:"path"`
}

// SeoPublicPageResponse chứa HTML và header kỹ thuật của trang công khai.
type SeoPublicPageResponse struct {
	ID             uint64   `json:"id"`
	Slug           string   `json:"slug"`
	CanonicalURL   string   `json:"canonicalUrl"`
	HTML           string   `json:"html"`
	ETag           string   `json:"etag"`
	CacheControl   string   `json:"cacheControl"`
	StatusCode     int      `json:"statusCode"`
	ContentType    string   `json:"contentType"`
	NoIndex        bool     `json:"noIndex"`
	RenderStatus   string   `json:"renderStatus"`
	PageStatus     string   `json:"pageStatus"`
	Warnings       []string `json:"warnings"`
	RedirectTo     string   `json:"redirectTo"`
	StaticHTMLPath string   `json:"staticHtmlPath"`
}

// SeoSitemapXMLResponse chứa XML sitemap đã tạo.
type SeoSitemapXMLResponse struct {
	XML         string `json:"xml"`
	ContentType string `json:"contentType"`
	Total       int    `json:"total"`
}

// SeoPublicDocumentResponse chứa tài liệu SEO dùng chung cho web và ứng dụng.
type SeoPublicDocumentResponse struct {
	Document     *SeoPublicDocument `json:"document"`
	ETag         string             `json:"etag"`
	CacheControl string             `json:"cacheControl"`
	StatusCode   int                `json:"statusCode"`
	Warnings     []string           `json:"warnings,omitempty"`
}
