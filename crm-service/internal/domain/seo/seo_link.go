package seo_domain

import (
	"time"

	"gorm.io/gorm"
)

const (
	// SeoLinkTypeNavigation là liên kết điều hướng thông thường.
	SeoLinkTypeNavigation = "navigation"
	// SeoLinkTypeBreadcrumb là liên kết trong chuỗi phân cấp.
	SeoLinkTypeBreadcrumb = "breadcrumb"
	// SeoLinkTypeRelated là liên kết tới nội dung liên quan.
	SeoLinkTypeRelated = "related"
	// SeoLinkTypeContextual là liên kết nằm trong nội dung.
	SeoLinkTypeContextual = "contextual"
	// SeoLinkTypeFooter là liên kết ở chân trang.
	SeoLinkTypeFooter = "footer"
)

// SeoInternalLink là liên kết HTML từ một trang SEO đến trang SEO khác.
type SeoInternalLink struct {
	// ID là khóa chính của liên kết.
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// ParentSeoID là mã trang chứa liên kết.
	ParentSeoID uint64 `gorm:"column:parent_seo_id;not null;index"`
	// ChildSeoID là mã trang được liên kết tới.
	ChildSeoID uint64 `gorm:"column:child_seo_id;not null;index"`

	// Link là URL được ghi vào HTML.
	Link string `gorm:"type:text;not null"`

	// Title là nội dung hiển thị của liên kết.
	Title string `gorm:"type:text;not null"`

	// LinkType xác định vị trí hoặc mục đích kỹ thuật của liên kết.
	LinkType string `gorm:"column:link_type;type:varchar(50);not null;default:'navigation';index"`

	// Priority xác định thứ tự render liên kết, từ 1 đến 10.
	Priority uint32 `gorm:"not null;default:5;index"`

	// ParentSeo là dữ liệu trang chứa liên kết.
	ParentSeo *SeoDomain `gorm:"foreignKey:ParentSeoID;references:ID"`
	// ChildSeo là dữ liệu trang được liên kết tới.
	ChildSeo *SeoDomain `gorm:"foreignKey:ChildSeoID;references:ID"`

	// CreatedAt là thời điểm tạo liên kết.
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	// UpdatedAt là thời điểm cập nhật liên kết gần nhất.
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	// DeletedAt hỗ trợ xóa mềm liên kết.
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (SeoInternalLink) TableName() string {
	return "seo_links"
}
