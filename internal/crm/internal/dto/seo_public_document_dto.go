package dto

import "encoding/json"

const SeoPublicDocumentSchemaVersion = "qhpro-public-seo-document/v1"

// SeoPublicDocument là hợp đồng nội dung công khai dùng chung cho web và ứng dụng.
type SeoPublicDocument struct {
	SchemaVersion   string             `json:"schemaVersion"`
	DocumentVersion string             `json:"documentVersion"`
	ModuleID        string             `json:"moduleId"`
	ResourceType    string             `json:"resourceType"`
	Identity        SeoPublicIdentity  `json:"identity"`
	Heading         SeoPublicHeading   `json:"heading"`
	Breadcrumbs     []SeoBreadcrumb    `json:"breadcrumbs"`
	Facts           []SeoFact          `json:"facts"`
	Sections        []SeoSection       `json:"sections"`
	RelatedLinks    []SeoRelatedLink   `json:"relatedLinks"`
	RelatedEntities []SeoRelatedEntity `json:"relatedEntities,omitempty"`
	Media           *SeoMediaBundle    `json:"media,omitempty"`
	FAQ             *SeoFAQ            `json:"faq,omitempty"`
	MapTarget       *SeoMapTarget      `json:"mapTarget,omitempty"`
	SEO             SeoPublicSEO       `json:"seo"`
	Source          SeoPublicSource    `json:"source"`
}

// SeoPublicIdentity chứa định danh và đường dẫn chính thức của trang.
type SeoPublicIdentity struct {
	ID            string `json:"id"`
	Slug          string `json:"slug"`
	CanonicalPath string `json:"canonicalPath"`
}

// SeoPublicHeading chứa nội dung mở đầu của trang.
type SeoPublicHeading struct {
	Eyebrow string `json:"eyebrow,omitempty"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

// SeoBreadcrumb là một cấp điều hướng trong đường dẫn phân cấp.
type SeoBreadcrumb struct {
	Label string `json:"label"`
	Path  string `json:"path"`
}

// SeoFact là một thông tin ngắn được trình bày theo nhãn và giá trị.
type SeoFact struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}

// SeoSection là một nhóm nội dung có tiêu đề riêng.
type SeoSection struct {
	ID     string            `json:"id"`
	Title  string            `json:"title"`
	Blocks []SeoContentBlock `json:"blocks"`
}

// SeoContentBlock là đơn vị nội dung nhỏ nhất được renderer hỗ trợ.
type SeoContentBlock struct {
	Type          string            `json:"type"`
	Text          string            `json:"text,omitempty"`
	Items         []string          `json:"-"`
	Tone          string            `json:"tone,omitempty"`
	Title         string            `json:"title,omitempty"`
	TimelineItems []SeoTimelineItem `json:"-"`
}

func (b *SeoContentBlock) UnmarshalJSON(data []byte) error {
	type blockWire struct {
		Type  string          `json:"type"`
		Text  string          `json:"text,omitempty"`
		Items json.RawMessage `json:"items,omitempty"`
		Tone  string          `json:"tone,omitempty"`
		Title string          `json:"title,omitempty"`
	}
	var wire blockWire
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	b.Type = wire.Type
	b.Text = wire.Text
	b.Tone = wire.Tone
	b.Title = wire.Title
	b.Items = nil
	b.TimelineItems = nil
	if len(wire.Items) == 0 || string(wire.Items) == "null" {
		return nil
	}
	if wire.Type == "timeline" {
		return json.Unmarshal(wire.Items, &b.TimelineItems)
	}
	return json.Unmarshal(wire.Items, &b.Items)
}

func (b SeoContentBlock) MarshalJSON() ([]byte, error) {
	type blockWire struct {
		Type  string `json:"type"`
		Text  string `json:"text,omitempty"`
		Items any    `json:"items,omitempty"`
		Tone  string `json:"tone,omitempty"`
		Title string `json:"title,omitempty"`
	}
	var items any
	if b.Type == "timeline" {
		if len(b.TimelineItems) > 0 {
			items = b.TimelineItems
		}
	} else if len(b.Items) > 0 {
		items = b.Items
	}
	return json.Marshal(blockWire{Type: b.Type, Text: b.Text, Items: items, Tone: b.Tone, Title: b.Title})
}

// SeoTimelineItem là một mốc thời gian trong nội dung.
type SeoTimelineItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Date        string `json:"date,omitempty"`
	Description string `json:"description,omitempty"`
}

// SeoRelatedLink là liên kết điều hướng bổ sung của trang.
type SeoRelatedLink struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Path        string `json:"path"`
	Description string `json:"description,omitempty"`
}

// SeoRelatedEntity là dữ liệu liên quan có thể mở trang chi tiết hoặc bản đồ.
type SeoRelatedEntity struct {
	ID           string `json:"id"`
	RelationType string `json:"relationType"`
	EntityType   string `json:"entityType"`
	EntityID     string `json:"entityId"`
	Title        string `json:"title"`
	Description  string `json:"description,omitempty"`
	PublicPath   string `json:"publicPath,omitempty"`
	MapPath      string `json:"mapPath,omitempty"`
}

// SeoMedia mô tả một ảnh hoặc tài nguyên trực quan của trang.
type SeoMedia struct {
	ID     string     `json:"id"`
	Role   string     `json:"role"`
	URL    string     `json:"url,omitempty"`
	Alt    string     `json:"alt"`
	Width  int        `json:"width,omitempty"`
	Height int        `json:"height,omitempty"`
	Bounds [4]float64 `json:"bounds,omitempty"`
	Order  int        `json:"order"`
}

// SeoMediaBundle gom ảnh chính và các ảnh bổ sung.
type SeoMediaBundle struct {
	Primary      *SeoMedia  `json:"primary,omitempty"`
	Items        []SeoMedia `json:"items"`
	Completeness string     `json:"completeness"`
	FitMode      string     `json:"fitMode,omitempty"`
	Note         string     `json:"note,omitempty"`
}

// SeoFAQ là nhóm câu hỏi thường gặp của trang.
type SeoFAQ struct {
	GroupKey string       `json:"groupKey"`
	Items    []SeoFAQItem `json:"items"`
}

// SeoFAQItem là một cặp câu hỏi và câu trả lời.
type SeoFAQItem struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// SeoMapTarget là đường dẫn mở vị trí tương ứng trên bản đồ.
type SeoMapTarget struct {
	Path  string `json:"path"`
	Label string `json:"label"`
}

// SeoPublicSEO chứa các thẻ kỹ thuật gửi cho máy tìm kiếm.
type SeoPublicSEO struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	CanonicalPath string `json:"canonicalPath"`
	Indexable     bool   `json:"indexable"`
	Follow        bool   `json:"follow"`
	ImageURL      string `json:"imageUrl,omitempty"`
}

// SeoPublicSource mô tả nguồn và thời điểm cập nhật nội dung.
type SeoPublicSource struct {
	SourceName  string `json:"sourceName,omitempty"`
	SourceURL   string `json:"sourceUrl,omitempty"`
	UpdatedAt   string `json:"updatedAt"`
	PublishedAt string `json:"publishedAt,omitempty"`
}
