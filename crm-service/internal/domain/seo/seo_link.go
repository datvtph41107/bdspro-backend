package seo_domain

import (
	"time"

	"gorm.io/gorm"
)

const (
	SeoLinkTypeNavigation = "navigation"
	SeoLinkTypeBreadcrumb = "breadcrumb"
	SeoLinkTypeRelated    = "related"
	SeoLinkTypeContextual = "contextual"
	SeoLinkTypeFooter     = "footer"
)

// SeoInternalLink là đường link thật render ra HTML giữa hai SEO page.
// Nó dùng để phân bổ authority, tạo breadcrumb, related pages và navigation.
type SeoInternalLink struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ParentSeoID uint64 `gorm:"column:parent_seo_id;not null;index"`
	ChildSeoID  uint64 `gorm:"column:child_seo_id;not null;index"`

	// Link là URL thật; nếu request không truyền thì usecase lấy child.CanonicalURL.
	Link string `gorm:"type:text;not null"`

	// Title là anchor text SEO.
	Title string `gorm:"type:text;not null"`

	// LinkType giúp phân biệt navigation/breadcrumb/related/contextual/footer.
	LinkType string `gorm:"column:link_type;type:varchar(50);not null;default:'navigation';index"`

	// Priority dùng khi render hoặc chọn link quan trọng.
	Priority uint32 `gorm:"not null;default:5;index"`

	ParentSeo *SeoDomain `gorm:"foreignKey:ParentSeoID;references:ID"`
	ChildSeo  *SeoDomain `gorm:"foreignKey:ChildSeoID;references:ID"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (SeoInternalLink) TableName() string {
	return "seo_links"
}
