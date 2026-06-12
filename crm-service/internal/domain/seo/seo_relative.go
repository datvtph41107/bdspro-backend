package seo_domain

import (
	"time"

	"gorm.io/gorm"
)

const (
	SeoRelationTypeParentChild = "parent_child"
	SeoRelationTypeRelated     = "related"
	SeoRelationTypeSameRegion  = "same_region"
	SeoRelationTypeSameTopic   = "same_topic"
)

// SeoRelative là quan hệ ngữ nghĩa giữa hai page SEO.
// Khác với seo_internal_link: relation không nhất thiết render thành link HTML.
type SeoRelative struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ParentSeoID uint64 `gorm:"column:parent_seo_id;not null;index"`
	ChildSeoID  uint64 `gorm:"column:child_seo_id;not null;index"`

	RelationType string `gorm:"column:relation_type;type:varchar(50);not null;default:'parent_child';index"`
	Priority     uint32 `gorm:"not null;default:5;index"`

	ParentSeo *SeoDomain `gorm:"foreignKey:ParentSeoID;references:ID"`
	ChildSeo  *SeoDomain `gorm:"foreignKey:ChildSeoID;references:ID"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (SeoRelative) TableName() string {
	return "seo_relative"
}
