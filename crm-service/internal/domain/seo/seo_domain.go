package seo_domain

import (
	"crm/internal/enums"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	// SeoScopePublic cho phép trang xuất hiện trên website công khai.
	SeoScopePublic = "public"
	// SeoScopeInternal chỉ cho phép hệ thống nội bộ sử dụng trang.
	SeoScopeInternal = "internal"
	// SeoScopeApp dành cho nội dung hiển thị trong ứng dụng.
	SeoScopeApp = "app"
	// SeoScopeAdmin dành cho màn hình quản trị.
	SeoScopeAdmin = "admin"
)

// Các giá trị dưới đây quy định tần suất gợi ý cho sitemap.
const (
	SeoChangeFreqAlways  = "always"
	SeoChangeFreqHourly  = "hourly"
	SeoChangeFreqDaily   = "daily"
	SeoChangeFreqWeekly  = "weekly"
	SeoChangeFreqMonthly = "monthly"
	SeoChangeFreqYearly  = "yearly"
	SeoChangeFreqNever   = "never"
)

const (
	// SeoSourceStatusManual cho biết trang được nhập thủ công.
	SeoSourceStatusManual = "manual"
	// SeoSourceStatusLinked cho biết trang đang liên kết với dữ liệu nguồn.
	SeoSourceStatusLinked = "linked"
	// SeoSourceStatusStale cho biết dữ liệu nguồn đã mới hơn nội dung trang.
	SeoSourceStatusStale = "stale"
	// SeoSourceStatusMissing cho biết bản ghi nguồn không còn tồn tại.
	SeoSourceStatusMissing = "missing"
	// SeoSourceStatusPermissionDenied cho biết không thể đọc dữ liệu nguồn.
	SeoSourceStatusPermissionDenied = "permission_denied"
)

// Các trạng thái vòng đời của một trang SEO.
const (
	SeoPageStatusDraft     = "draft"
	SeoPageStatusPublished = "published"
	SeoPageStatusArchived  = "archived"
)

// Các trạng thái của quá trình tạo HTML.
const (
	SeoRenderStatusNone      = "none"
	SeoRenderStatusPending   = "pending"
	SeoRenderStatusRendering = "rendering"
	SeoRenderStatusSuccess   = "success"
	SeoRenderStatusFailed    = "failed"
)

// SeoDomain là bản ghi trung tâm của một trang SEO.
// Bảng chỉ giữ nội dung, trạng thái xuất bản, cấu hình máy tìm kiếm và kết quả render.
type SeoDomain struct {
	// ID là khóa chính của trang SEO.
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// Slug là phần định danh cuối đường dẫn.
	Slug string `gorm:"type:varchar(255);index"`
	// OriginURL là đường dẫn gốc trước khi chuẩn hóa.
	OriginURL string `gorm:"column:origin_url;type:text;index"`
	// CanonicalURL là URL chính thức dành cho máy tìm kiếm.
	CanonicalURL string `gorm:"column:canonical_url;type:text;not null;index"`

	// RefType xác định loại dữ liệu nguồn được liên kết.
	RefType enums.ESEORefType `gorm:"column:ref_type;type:smallint;default:0;index"`
	// RefID là mã bản ghi ở hệ thống nguồn; nil khi nhập thủ công.
	RefID *uint64 `gorm:"column:ref_id;type:bigint;index"`

	// RefSource là khóa bộ kết nối dùng để đọc dữ liệu nguồn.
	RefSource string `gorm:"column:ref_source;type:varchar(80);not null;default:'manual';index"`
	// RefLabel là tên hiển thị của bản ghi nguồn.
	RefLabel string `gorm:"column:ref_label;type:text"`
	// RefURL là đường dẫn tham chiếu đến dữ liệu nguồn.
	RefURL string `gorm:"column:ref_url;type:text"`
	// SourceStatus là trạng thái đồng bộ hiện tại của dữ liệu nguồn.
	SourceStatus string `gorm:"column:source_status;type:varchar(40);not null;default:'manual';index"`
	// RefMissing đánh dấu bản ghi nguồn không còn tồn tại.
	RefMissing bool `gorm:"column:ref_missing;not null;default:false;index"`
	// RefSnapshotJSON là bản chụp dữ liệu dùng để render ổn định.
	RefSnapshotJSON datatypes.JSON `gorm:"column:ref_snapshot_json;type:jsonb;not null;default:'{}'::jsonb"`
	// RefHash dùng để phát hiện dữ liệu nguồn đã thay đổi.
	RefHash string `gorm:"column:ref_hash;type:varchar(64);index"`
	// RefLastSyncedAt là thời điểm đồng bộ nguồn gần nhất.
	RefLastSyncedAt *time.Time `gorm:"column:ref_last_synced_at;index"`

	// Scope xác định nơi trang được phép sử dụng.
	Scope string `gorm:"type:varchar(50);not null;default:'public';index"`

	// PageStatus là trạng thái draft, published hoặc archived.
	PageStatus string `gorm:"column:page_status;type:varchar(30);not null;default:'draft';index"`

	// Title là tiêu đề chính và tiêu đề SEO mặc định.
	Title string `gorm:"type:text"`
	// Description là mô tả dùng cho thẻ meta description.
	Description string `gorm:"type:text"`
	// Content là nội dung đầy đủ của trang.
	Content string `gorm:"type:text"`
	// Summary là phần tóm tắt hiển thị đầu trang.
	Summary string `gorm:"type:text"`

	// Published cho biết trang đã được phát hành hay chưa.
	Published bool `gorm:"default:false;index"`
	// PublishedAt là thời điểm phát hành gần nhất.
	PublishedAt *time.Time `gorm:"column:published_at;index"`

	// IsSiteMap quyết định trang có được đưa vào sitemap hay không.
	IsSiteMap bool `gorm:"column:is_site_map;default:false;index"`
	// IsIndex quyết định máy tìm kiếm có được lập chỉ mục hay không.
	IsIndex bool `gorm:"column:is_index;default:false;index"`
	// IsRobot quyết định có cho bot thu thập trang hay không.
	IsRobot *bool `gorm:"column:is_robot;type:bool;default:false"`

	// SiteMapLastedAt là thời điểm dữ liệu sitemap được cập nhật.
	SiteMapLastedAt *time.Time `gorm:"column:site_map_lasted_at;index"`
	// SitemapPriority là độ ưu tiên kỹ thuật ghi trong sitemap.
	SitemapPriority float32 `gorm:"column:sitemap_priority;type:real;default:0.5"`
	// SitemapChangeFreq là tần suất thay đổi gợi ý cho máy tìm kiếm.
	SitemapChangeFreq string `gorm:"column:sitemap_change_freq;type:varchar(20);default:'daily'"`

	// NeedGenerate đánh dấu trang đang cần tạo lại HTML.
	NeedGenerate bool `gorm:"column:need_generate;default:true;index"`
	// GeneratedAt là thời điểm tạo HTML thành công gần nhất.
	GeneratedAt *time.Time `gorm:"column:generated_at;index"`
	// SourceUpdatedAt là thời điểm thay đổi gần nhất từ dữ liệu nguồn.
	SourceUpdatedAt *time.Time `gorm:"column:source_updated_at;index"`

	// RenderStatus là trạng thái hiện tại của quá trình tạo HTML.
	RenderStatus string `gorm:"column:render_status;type:varchar(30);not null;default:'none';index"`
	// RenderedHTML là HTML hoàn chỉnh được phục vụ công khai.
	RenderedHTML string `gorm:"column:rendered_html;type:text"`
	// LastRenderError lưu lỗi render gần nhất để chẩn đoán.
	LastRenderError string `gorm:"column:last_render_error;type:text"`
	// StaticHtmlPath là vị trí tệp HTML đã tạo.
	StaticHtmlPath string `gorm:"column:static_html_path;type:text"`
	// StaticHtmlHash dùng để nhận biết nội dung HTML thay đổi.
	StaticHtmlHash string `gorm:"column:static_html_hash;type:varchar(64);index"`

	// TemplateKey chọn mẫu HTML dùng để render.
	TemplateKey string `gorm:"column:template_key;type:varchar(120);index"`
	// TemplateVersion ghi phiên bản mẫu đã sử dụng.
	TemplateVersion string `gorm:"column:template_version;type:varchar(80)"`

	// DeepLink là đường dẫn mở đúng màn hình trong ứng dụng.
	DeepLink string `gorm:"column:deep_link;type:text"`

	// Metadata chứa JSON-LD, breadcrumb, Open Graph và cấu hình mở rộng.
	Metadata datatypes.JSON `gorm:"type:jsonb;default:'{}'::jsonb"`

	// Note là ghi chú vận hành dành cho quản trị viên.
	Note string `gorm:"type:text"`

	// InternalLinks là các liên kết HTML đi từ trang này.
	InternalLinks []SeoInternalLink `gorm:"foreignKey:ParentSeoID;references:ID"`
	// Relatives là các trang có quan hệ dữ liệu với trang này.
	Relatives []SeoRelative `gorm:"foreignKey:ParentSeoID;references:ID"`

	// CreatedAt là thời điểm tạo bản ghi.
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	// UpdatedAt là thời điểm cập nhật bản ghi gần nhất.
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	// DeletedAt hỗ trợ xóa mềm bản ghi.
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (SeoDomain) TableName() string {
	return "seo_domain"
}

// GetSchema trả về loại schema nội dung dùng khi tạo dữ liệu có cấu trúc.
func (s *SeoDomain) GetSchema() string {
	switch s.RefType {
	case enums.ESEORefTypeAdmUnit:
		return "administrative_unit"
	case enums.ESEORefTypeRegion:
		return "planning_region"
	case enums.ESEORefTypeParcel:
		return "planning_parcel"
	case enums.ESEORefTypeProject:
		return "planning_project"
	case enums.ESEORefTypeDocument:
		return "legal_document"
	case enums.ESEORefTypeMap:
		return "gis_layer"
	case enums.ESEORefTypeReport:
		return "report"
	default:
		return "info"
	}
}

// NormalizeLifecycle bổ sung trạng thái mặc định cho dữ liệu cũ.
func (s *SeoDomain) NormalizeLifecycle() {
	if s == nil {
		return
	}

	if s.PageStatus == "" {
		if s.DeletedAt.Valid {
			s.PageStatus = SeoPageStatusArchived
		} else if s.Published {
			s.PageStatus = SeoPageStatusPublished
		} else {
			s.PageStatus = SeoPageStatusDraft
		}
	}

	if s.RenderStatus == "" {
		if s.NeedGenerate {
			s.RenderStatus = SeoRenderStatusPending
		} else if s.StaticHtmlHash != "" || s.GeneratedAt != nil {
			s.RenderStatus = SeoRenderStatusSuccess
		} else {
			s.RenderStatus = SeoRenderStatusNone
		}
	}
}

// IsPublicIndexable kiểm tra trực tiếp điều kiện để trang được lập chỉ mục.
func (s *SeoDomain) IsPublicIndexable() bool {
	if s == nil {
		return false
	}
	return s.PageStatus == SeoPageStatusPublished &&
		s.Published &&
		s.IsIndex &&
		s.Scope == SeoScopePublic
}

// IsSitemapEligible kiểm tra trang có đủ điều kiện xuất hiện trong sitemap.
func (s *SeoDomain) IsSitemapEligible() bool {
	if s == nil {
		return false
	}
	return s.IsPublicIndexable() && s.IsSiteMap
}

// MarkRenderPending đưa trang vào trạng thái chờ tạo HTML.
func (s *SeoDomain) MarkRenderPending() {
	if s == nil {
		return
	}
	s.RenderStatus = SeoRenderStatusPending
	s.NeedGenerate = true
	s.LastRenderError = ""
}

// MarkRenderRunning đánh dấu tiến trình tạo HTML đang chạy.
func (s *SeoDomain) MarkRenderRunning() {
	if s == nil {
		return
	}
	s.RenderStatus = SeoRenderStatusRendering
}

// MarkRenderSuccess lưu kết quả tạo HTML thành công.
func (s *SeoDomain) MarkRenderSuccess(staticPath string, staticHash string, generatedAt time.Time) {
	if s == nil {
		return
	}
	s.RenderStatus = SeoRenderStatusSuccess
	s.NeedGenerate = false
	s.GeneratedAt = &generatedAt
	s.SiteMapLastedAt = &generatedAt
	s.StaticHtmlPath = staticPath
	s.StaticHtmlHash = staticHash
	s.LastRenderError = ""
}

// MarkRenderFailed lưu lỗi và cho phép thử tạo HTML lại.
func (s *SeoDomain) MarkRenderFailed(message string) {
	if s == nil {
		return
	}
	s.RenderStatus = SeoRenderStatusFailed
	s.NeedGenerate = true
	s.LastRenderError = message
}
