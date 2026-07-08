package seo_domain

import (
	"crm/internal/enums"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	SeoScopePublic   = "public"
	SeoScopeInternal = "internal"
	SeoScopeApp      = "app"
	SeoScopeAdmin    = "admin"
)

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
	SeoSourceStatusManual           = "manual"
	SeoSourceStatusLinked           = "linked"
	SeoSourceStatusStale            = "stale"
	SeoSourceStatusMissing          = "missing"
	SeoSourceStatusPermissionDenied = "permission_denied"
)

const (
	SeoPageStatusDraft     = "draft"
	SeoPageStatusPublished = "published"
	SeoPageStatusArchived  = "archived"
)

const (
	SeoRenderStatusNone      = "none"
	SeoRenderStatusPending   = "pending"
	SeoRenderStatusRendering = "rendering"
	SeoRenderStatusSuccess   = "success"
	SeoRenderStatusFailed    = "failed"
)

// SeoDomain là SEO content registry trung tâm.
//
// Architectural role:
//   - It is not the source data owner.
//   - It stores the SEO page lifecycle, source snapshot, render state,
//     sitemap/index policy and operational metadata.
//   - Source data remains owned by BDSPro/TQD/CRM upstream.
//   - RefSnapshotJSON is an immutable-ish snapshot at create/sync time.
type SeoDomain struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// Identity: định danh page SEO ổn định.
	Slug         string `gorm:"type:varchar(255);index"`
	OriginURL    string `gorm:"column:origin_url;type:text;index"`
	CanonicalURL string `gorm:"column:canonical_url;type:text;not null;index"`

	// Source binding: page này đại diện cho dữ liệu nào trong hệ thống.
	RefType enums.ESEORefType `gorm:"column:ref_type;type:smallint;default:0;index"`
	RefID   *uint64           `gorm:"column:ref_id;type:bigint;index"`

	// Source snapshot/cache: dữ liệu đã normalize từ resolver để editor/generator/audit dùng ổn định.
	RefSource       string         `gorm:"column:ref_source;type:varchar(80);not null;default:'manual';index"`
	RefLabel        string         `gorm:"column:ref_label;type:text"`
	RefURL          string         `gorm:"column:ref_url;type:text"`
	SourceStatus    string         `gorm:"column:source_status;type:varchar(40);not null;default:'manual';index"`
	RefMissing      bool           `gorm:"column:ref_missing;not null;default:false;index"`
	RefSnapshotJSON datatypes.JSON `gorm:"column:ref_snapshot_json;type:jsonb;not null;default:'{}'::jsonb"`
	RefHash         string         `gorm:"column:ref_hash;type:varchar(64);index"`
	RefLastSyncedAt *time.Time     `gorm:"column:ref_last_synced_at;index"`
	QualityScore    uint32         `gorm:"column:quality_score;not null;default:0;index"`

	// Scope phân vùng vận hành page: public/internal/app/admin.
	Scope string `gorm:"type:varchar(50);not null;default:'public';index"`

	// SEO classify: A/B/C/D/N.
	Classify enums.ESeoClassify `gorm:"type:smallint;not null;default:50;index"`

	// Page lifecycle.
	// PageStatus is the official lifecycle state for command APIs.
	// Legacy Published boolean is kept for compatibility and fast filters.
	PageStatus string `gorm:"column:page_status;type:varchar(30);not null;default:'draft';index"`

	// Content lõi để render SEO page.
	Title       string `gorm:"type:text"`
	Description string `gorm:"type:text"`
	Content     string `gorm:"type:text"`
	Summary     string `gorm:"type:text"`

	// Publish lifecycle: page đã đủ điều kiện ra ngoài hay chưa.
	Published   bool       `gorm:"default:false;index"`
	PublishedAt *time.Time `gorm:"column:published_at;index"`

	// Search engine policy.
	IsSiteMap bool  `gorm:"column:is_site_map;default:false;index"`
	IsIndex   bool  `gorm:"column:is_index;default:false;index"`
	IsRobot   *bool `gorm:"column:is_robot;type:bool;default:false"`

	// Sitemap data.
	SiteMapLastedAt   *time.Time `gorm:"column:site_map_lasted_at;index"`
	SitemapPriority   float32    `gorm:"column:sitemap_priority;type:real;default:0.5"`
	SitemapChangeFreq string     `gorm:"column:sitemap_change_freq;type:varchar(20);default:'daily'"`

	// Render/generation pipeline.
	NeedGenerate    bool       `gorm:"column:need_generate;default:true;index"`
	GeneratedAt     *time.Time `gorm:"column:generated_at;index"`
	SourceUpdatedAt *time.Time `gorm:"column:source_updated_at;index"`

	RenderStatus    string `gorm:"column:render_status;type:varchar(30);not null;default:'none';index"`
	RenderedHTML    string `gorm:"column:rendered_html;type:text"`
	LastRenderError string `gorm:"column:last_render_error;type:text"`
	StaticHtmlPath  string `gorm:"column:static_html_path;type:text"`
	StaticHtmlHash  string `gorm:"column:static_html_hash;type:varchar(64);index"`

	// Template metadata used by backend renderer.
	TemplateKey     string `gorm:"column:template_key;type:varchar(120);index"`
	TemplateVersion string `gorm:"column:template_version;type:varchar(80)"`

	// Mobile app deeplink.
	DeepLink string `gorm:"column:deep_link;type:text"`

	// Metadata chỉ giữ phần mở rộng: schema JSON-LD, breadcrumb, og, robots, custom.
	Metadata datatypes.JSON `gorm:"type:jsonb;default:'{}'::jsonb"`

	Note string `gorm:"type:text"`

	InternalLinks []SeoInternalLink `gorm:"foreignKey:ParentSeoID;references:ID"`
	Relatives     []SeoRelative     `gorm:"foreignKey:ParentSeoID;references:ID"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (SeoDomain) TableName() string {
	return "seo_domain"
}

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

func (s *SeoDomain) IsPublicIndexable() bool {
	if s == nil {
		return false
	}
	return s.PageStatus == SeoPageStatusPublished &&
		s.Published &&
		s.IsIndex &&
		s.Scope == SeoScopePublic &&
		s.Classify != enums.ESeoClassifyN &&
		s.SourceStatus != SeoSourceStatusMissing &&
		s.SourceStatus != SeoSourceStatusPermissionDenied
}

func (s *SeoDomain) IsSitemapEligible() bool {
	if s == nil {
		return false
	}
	return s.IsPublicIndexable() && s.IsSiteMap
}

func (s *SeoDomain) MarkRenderPending() {
	if s == nil {
		return
	}
	s.RenderStatus = SeoRenderStatusPending
	s.NeedGenerate = true
	s.LastRenderError = ""
}

func (s *SeoDomain) MarkRenderRunning() {
	if s == nil {
		return
	}
	s.RenderStatus = SeoRenderStatusRendering
}

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

func (s *SeoDomain) MarkRenderFailed(message string) {
	if s == nil {
		return
	}
	s.RenderStatus = SeoRenderStatusFailed
	s.NeedGenerate = true
	s.LastRenderError = message
}
