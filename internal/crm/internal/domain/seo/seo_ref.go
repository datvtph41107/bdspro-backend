package seo_domain

import "crm/internal/enums"

// SeoRef lưu một dữ liệu nguồn liên quan bổ sung của trang SEO.
// Bảng dùng khi một trang cần tham chiếu nhiều địa bàn, thửa đất, đồ án hoặc văn bản.
type SeoRef struct {
	// ID là khóa chính của quan hệ tham chiếu.
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// SeoID là mã trang SEO chứa tham chiếu.
	SeoID uint64 `gorm:"column:seo_domain_id;not null;index"`
	// RefID là mã bản ghi tại hệ thống nguồn.
	RefID uint64 `gorm:"column:ref_id;not null;index"`
	// RefType là loại dữ liệu nguồn được tham chiếu.
	RefType enums.ESEORefType `gorm:"column:ref_type;type:int4;not null;index"`
}

func (SeoRef) TableName() string {
	return "seo_refs"
}
