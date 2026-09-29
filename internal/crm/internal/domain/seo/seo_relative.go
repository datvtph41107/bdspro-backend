package seo_domain

import (
	"time"

	"gorm.io/gorm"
)

const (
	// SeoRelationTypeParentChild là quan hệ cha con.
	SeoRelationTypeParentChild = "parent_child"
	// SeoRelationTypeRelated là quan hệ nội dung liên quan.
	SeoRelationTypeRelated = "related"
	// SeoRelationTypeSameRegion là quan hệ cùng địa bàn.
	SeoRelationTypeSameRegion = "same_region"
	// SeoRelationTypeSameTopic là quan hệ cùng chủ đề.
	SeoRelationTypeSameTopic = "same_topic"
)

// SeoRelative mô tả quan hệ dữ liệu giữa hai trang, không bắt buộc render thành liên kết.
type SeoRelative struct {
	// ID là khóa chính của quan hệ.
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// ParentSeoID là mã trang nguồn.
	ParentSeoID uint64 `gorm:"column:parent_seo_id;not null;index"`
	// ChildSeoID là mã trang liên quan.
	ChildSeoID uint64 `gorm:"column:child_seo_id;not null;index"`

	// RelationType là loại quan hệ giữa hai trang.
	RelationType string `gorm:"column:relation_type;type:varchar(50);not null;default:'parent_child';index"`
	// Priority xác định thứ tự ưu tiên khi đọc quan hệ.
	Priority uint32 `gorm:"not null;default:5;index"`

	// ParentSeo là dữ liệu trang nguồn.
	ParentSeo *SeoDomain `gorm:"foreignKey:ParentSeoID;references:ID"`
	// ChildSeo là dữ liệu trang liên quan.
	ChildSeo *SeoDomain `gorm:"foreignKey:ChildSeoID;references:ID"`

	// CreatedAt là thời điểm tạo quan hệ.
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	// UpdatedAt là thời điểm cập nhật quan hệ gần nhất.
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	// DeletedAt hỗ trợ xóa mềm quan hệ.
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (SeoRelative) TableName() string {
	return "seo_relative"
}
