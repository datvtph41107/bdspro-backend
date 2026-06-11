package seo_domain

import (
	"time"

	"gorm.io/gorm"
)

type SeoInternalLink struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	Link string `gorm:"type:text"`

	Title string `gorm:"type:text"`

	ParentSeoID uint64 `gorm:"not null;index"`

	ChildSeoID uint64 `gorm:"not null;index"`

	ParentSeo *SeoDomain `gorm:"foreignKey:ParentSeoID;references:ID"`

	ChildSeo *SeoDomain `gorm:"foreignKey:ChildSeoID;references:ID"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (SeoInternalLink) TableName() string {
	return "seo_links"
}
