package seo_domain

import (
	"time"

	"crm/internal/enums"
)

// SeoRef lưu một dữ liệu nguồn liên quan bổ sung của trang SEO.
// Bảng dùng khi một trang cần tham chiếu nhiều địa bàn, thửa đất, đồ án hoặc văn bản.
type SeoRef struct {
	// ID là khóa chính của quan hệ tham chiếu.
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	SeoID   uint64            `gorm:"column:seo_domain_id;not null;uniqueIndex:idx_seo_refs_domain_ref;index"`
	RefID   uint64            `gorm:"column:ref_id;not null;uniqueIndex:idx_seo_refs_domain_ref;index"`
	RefType enums.ESEORefType `gorm:"column:ref_type;type:int4;not null;default:800;uniqueIndex:idx_seo_refs_domain_ref;index"`

	SeoDomain *SeoDomain `gorm:"foreignKey:SeoID;references:ID"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (SeoRef) TableName() string {
	return "seo_refs"
}
