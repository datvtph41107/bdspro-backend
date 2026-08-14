package dto

// RenderSeoPreviewRequest chứa nội dung cần xem trước dưới dạng HTML.
type RenderSeoPreviewRequest struct {
	SourceKey    string  `json:"sourceKey"`
	RefID        *uint64 `json:"refId"`
	Slug         string  `json:"slug"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	CanonicalURL string  `json:"canonicalUrl"`
	Content      string  `json:"content"`
	Summary      string  `json:"summary"`
	Metadata     string  `json:"metadata"`
	TemplateKey  string  `json:"templateKey"`
}

// RenderExistingSeoPreviewRequest xem trước một trang đã tồn tại với nội dung tùy chọn.
type RenderExistingSeoPreviewRequest struct {
	ID          uint64  `json:"id"`
	Content     *string `json:"content"`
	TemplateKey *string `json:"templateKey"`
}

// RenderSeoPreviewResponse trả về HTML và thông tin tệp dự kiến.
type RenderSeoPreviewResponse struct {
	HTML       string   `json:"html"`
	StaticPath string   `json:"staticPath"`
	Hash       string   `json:"hash"`
	Warnings   []string `json:"warnings"`
}

// GenerateSeoPageRequest yêu cầu tạo HTML cho một trang SEO.
type GenerateSeoPageRequest struct {
	ID          uint64 `json:"id"`
	TriggerType string `json:"triggerType"`
	Reason      string `json:"reason"`
}

// GenerateSeoPageResponse trả về trạng thái và kết quả tạo HTML.
type GenerateSeoPageResponse struct {
	ID             uint64 `json:"id"`
	RenderStatus   string `json:"renderStatus"`
	StaticHTMLPath string `json:"staticHtmlPath"`
	StaticHTMLHash string `json:"staticHtmlHash"`
	GeneratedAt    string `json:"generatedAt"`
	ErrorMessage   string `json:"errorMessage"`
}
