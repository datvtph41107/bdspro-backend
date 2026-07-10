package seo_domain

import "crm/internal/enums"

// ngoài ref_type và ref_id chính từ seo thì đây là bảng chứa id trỏ đến các domain tương ứng
// nếu các domain đó cần link tới bài viết. có thể coi đây là 1 bảng ds các
// adminunit/parcel/project/document/report/... liên quan
// cụ thể là dùng ở chô chi tiết thư viện tab "tin"
type SeoRef struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	SeoID   uint64            `gorm:"column:seo_domain_id;not null;index"`
	RefID   uint64            `gorm:"column:ref_id;not null;index"`
	RefType enums.ESEORefType `gorm:"column:ref_type;type:int4;not null;index"`
}

func (SeoRef) TableName() string {
	return "seo_refs"
}
