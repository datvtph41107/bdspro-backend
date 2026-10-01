// Package entity defines GORM models for v4/modules/seo.
// No imports from crm-service; ref-type constants are cloned locally.
package entity

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ── SEO ref type constants (cloned from crm/enums) ───────────────────────────

type SEORefType uint32

const (
	SEORefTypeAdmUnit  SEORefType = 100
	SEORefTypeParcel   SEORefType = 200
	SEORefTypeRegion   SEORefType = 300
	SEORefTypeProject  SEORefType = 400
	SEORefTypeDocument SEORefType = 500
	SEORefTypeMap      SEORefType = 600
	SEORefTypeReport   SEORefType = 700
	SEORefTypeLegacy   SEORefType = 800
)

func (e SEORefType) Name() string {
	switch e {
	case SEORefTypeAdmUnit:
		return "Đơn vị hành chính"
	case SEORefTypeParcel:
		return "Thửa đất"
	case SEORefTypeRegion:
		return "Vùng"
	case SEORefTypeProject:
		return "Đồ án"
	case SEORefTypeDocument:
		return "Tài liệu"
	case SEORefTypeMap:
		return "Bản đồ"
	case SEORefTypeReport:
		return "Báo cáo"
	case SEORefTypeLegacy:
		return "Nội dung thủ công"
	default:
		return ""
	}
}

// ── Page status constants ─────────────────────────────────────────────────────

const (
	PageStatusDraft     = "draft"
	PageStatusPublished = "published"
	PageStatusArchived  = "archived"

	ScopePublic   = "public"
	ScopeInternal = "internal"
	ScopeApp      = "app"
	ScopeAdmin    = "admin"

	SourceStatusManual          = "manual"
	SourceStatusLinked          = "linked"
	SourceStatusStale           = "stale"
	SourceStatusMissing         = "missing"
	SourceStatusPermissionDenied = "permission_denied"

	RenderStatusNone      = "none"
	RenderStatusPending   = "pending"
	RenderStatusRendering = "rendering"
	RenderStatusSuccess   = "success"
	RenderStatusFailed    = "failed"
)

// ── SeoDomain ─────────────────────────────────────────────────────────────────

// SeoDomain is the central SEO page record for modules/seo.
type SeoDomain struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	Slug         string `gorm:"type:varchar(255);index"`
	OriginURL    string `gorm:"column:origin_url;type:text"`
	CanonicalURL string `gorm:"column:canonical_url;type:text;not null;default:''"`

	RefType         SEORefType     `gorm:"column:ref_type;type:smallint;default:0"`
	RefID           *uint64        `gorm:"column:ref_id;type:bigint"`
	RefSource       string         `gorm:"column:ref_source;type:varchar(80);default:'manual'"`
	RefLabel        string         `gorm:"column:ref_label;type:text"`
	RefURL          string         `gorm:"column:ref_url;type:text"`
	SourceStatus    string         `gorm:"column:source_status;type:varchar(40);default:'manual'"`
	RefMissing      bool           `gorm:"column:ref_missing;default:false"`
	RefSnapshotJSON datatypes.JSON `gorm:"column:ref_snapshot_json;type:jsonb;not null;default:'{}'::jsonb"`
	RefHash         string         `gorm:"column:ref_hash;type:varchar(64)"`
	RefLastSyncedAt *time.Time     `gorm:"column:ref_last_synced_at"`

	Scope      string `gorm:"type:varchar(50);default:'public'"`
	PageStatus string `gorm:"column:page_status;type:varchar(30);default:'draft'"`

	Title       string `gorm:"type:text"`
	Description string `gorm:"type:text"`
	Content     string `gorm:"type:text"`
	Summary     string `gorm:"type:text"`

	VisibleOn  uint32  `gorm:"column:visible_on;type:smallint;default:30"`
	CategoryID *uint64 `gorm:"column:category_id;type:bigint"`
	ProvinceID *uint64 `gorm:"column:province_id;type:bigint"`
	ImageURL   string  `gorm:"column:image_url;type:text"`

	Published   bool       `gorm:"default:false"`
	PublishedAt *time.Time `gorm:"column:published_at"`

	IsSiteMap         bool       `gorm:"column:is_site_map;default:false"`
	IsIndex           bool       `gorm:"column:is_index;default:false"`
	IsRobot           *bool      `gorm:"column:is_robot;type:bool;default:false"`
	SiteMapLastedAt   *time.Time `gorm:"column:site_map_lasted_at"`
	SitemapPriority   float32    `gorm:"column:sitemap_priority;type:real;default:0.5"`
	SitemapChangeFreq string     `gorm:"column:sitemap_change_freq;type:varchar(20);default:'daily'"`

	NeedGenerate    bool       `gorm:"column:need_generate;default:true"`
	GeneratedAt     *time.Time `gorm:"column:generated_at"`
	SourceUpdatedAt *time.Time `gorm:"column:source_updated_at"`

	RenderStatus    string `gorm:"column:render_status;type:varchar(30);default:'none'"`
	RenderedHTML    string `gorm:"column:rendered_html;type:text"`
	LastRenderError string `gorm:"column:last_render_error;type:text"`
	StaticHtmlPath  string `gorm:"column:static_html_path;type:text"`
	StaticHtmlHash  string `gorm:"column:static_html_hash;type:varchar(64)"`

	TemplateKey     string         `gorm:"column:template_key;type:varchar(120)"`
	TemplateVersion string         `gorm:"column:template_version;type:varchar(80)"`
	DeepLink        string         `gorm:"column:deep_link;type:text"`
	Metadata        datatypes.JSON `gorm:"type:jsonb;default:'{}'::jsonb"`
	Note            string         `gorm:"type:text"`

	InternalLinks []SeoLink     `gorm:"foreignKey:ParentSeoID;references:ID"`
	Relatives     []SeoRelative `gorm:"foreignKey:ParentSeoID;references:ID"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	// Read-only aggregates from joins (not DB columns on seo_domain):
	ClickCount   int64  `gorm:"column:click_count;->"`
	CategoryName string `gorm:"column:category_name;->"`
}

func (SeoDomain) TableName() string { return "seo_domain" }

func (s *SeoDomain) NormalizeLifecycle() {
	if s == nil {
		return
	}
	if s.PageStatus == "" {
		if s.DeletedAt.Valid {
			s.PageStatus = PageStatusArchived
		} else if s.Published {
			s.PageStatus = PageStatusPublished
		} else {
			s.PageStatus = PageStatusDraft
		}
	}
	if s.RenderStatus == "" {
		if s.NeedGenerate {
			s.RenderStatus = RenderStatusPending
		} else if s.GeneratedAt != nil {
			s.RenderStatus = RenderStatusSuccess
		} else {
			s.RenderStatus = RenderStatusNone
		}
	}
}

func (s *SeoDomain) IsSitemapEligible() bool {
	return s != nil &&
		s.PageStatus == PageStatusPublished &&
		s.Published && s.IsIndex && s.IsSiteMap &&
		s.Scope == ScopePublic
}

// ── SeoLink ───────────────────────────────────────────────────────────────────

type SeoLink struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement"`
	ParentSeoID uint64         `gorm:"column:parent_seo_id;not null"`
	ChildSeoID  *uint64        `gorm:"column:child_seo_id"`
	Href        string         `gorm:"type:text;default:''"`
	AnchorText  string         `gorm:"column:anchor_text;type:text"`
	Priority    int            `gorm:"default:0"`
	ChildSeo    *SeoDomain     `gorm:"foreignKey:ChildSeoID"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (SeoLink) TableName() string { return "seo_links" }

// ── SeoRelative ───────────────────────────────────────────────────────────────

type SeoRelative struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement"`
	ParentSeoID  uint64         `gorm:"column:parent_seo_id;not null"`
	ChildSeoID   *uint64        `gorm:"column:child_seo_id"`
	RelationType string         `gorm:"column:relation_type;type:varchar(60)"`
	Priority     int            `gorm:"default:0"`
	ChildSeo     *SeoDomain     `gorm:"foreignKey:ChildSeoID"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (SeoRelative) TableName() string { return "seo_relative" }

// ── SeoGenerationLog ──────────────────────────────────────────────────────────

type SeoGenerationLog struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement"`
	SeoDomainID uint64     `gorm:"column:seo_domain_id;not null"`
	Status      string     `gorm:"type:varchar(30);default:'pending'"`
	Message     string     `gorm:"type:text"`
	StartedAt   *time.Time `gorm:"column:started_at"`
	FinishedAt  *time.Time `gorm:"column:finished_at"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
}

func (SeoGenerationLog) TableName() string { return "seo_generation_log" }

// ── SeoCount ──────────────────────────────────────────────────────────────────

type SeoCount struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	SeoDomainID uint64    `gorm:"column:seo_domain_id;not null;uniqueIndex"`
	ClickCount  int64     `gorm:"column:click_count;default:0"`
	ViewCount   int64     `gorm:"column:view_count;default:0"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (SeoCount) TableName() string { return "seo_count" }

// ── SeoCountLog ───────────────────────────────────────────────────────────────

type SeoCountLog struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	SeoDomainID uint64    `gorm:"column:seo_domain_id;not null"`
	EventType   string    `gorm:"column:event_type;type:varchar(20);default:'click'"`
	UserID      *uint64   `gorm:"column:user_id"`
	IPAddr      string    `gorm:"column:ip_addr;type:varchar(45)"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (SeoCountLog) TableName() string { return "seo_count_log" }

// ── SeoCategory ───────────────────────────────────────────────────────────────

type SeoCategory struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement"`
	Code        string         `gorm:"type:varchar(80);uniqueIndex"`
	Slug        string         `gorm:"type:varchar(120)"`
	Name        string         `gorm:"type:varchar(255)"`
	Description string         `gorm:"type:text"`
	Status      string         `gorm:"type:varchar(30);default:'active'"`
	SortOrder   int            `gorm:"column:sort_order;default:0"`
	Metadata    datatypes.JSON `gorm:"column:metadata;type:jsonb;not null;default:'{}'::jsonb"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (SeoCategory) TableName() string { return "seo_category" }

// ── UserSavedPlanningNews ─────────────────────────────────────────────────────

type UserSavedPlanningNews struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement"`
	UserID      uint64         `gorm:"column:user_id;not null;uniqueIndex:uq_user_saved_planning_news"`
	SeoDomainID uint64         `gorm:"column:seo_domain_id;not null;uniqueIndex:uq_user_saved_planning_news"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (UserSavedPlanningNews) TableName() string { return "user_saved_planning_news" }
