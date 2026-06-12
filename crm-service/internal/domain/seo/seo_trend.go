package seo_domain

import (
	"time"

	"gorm.io/gorm"
)

// SeoTrend lưu metric SEO theo ngày để đo hiệu quả về impressions/clicks/CTR/position.
type SeoTrend struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	SeoDomainID uint64    `gorm:"column:seo_domain_id;not null;index"`
	Date        time.Time `gorm:"type:date;not null;index"`

	Impressions uint64  `gorm:"not null;default:0"`
	Clicks      uint64  `gorm:"not null;default:0"`
	CTR         float64 `gorm:"column:ctr;type:double precision;not null;default:0"`
	AvgPosition float64 `gorm:"column:avg_position;type:double precision;not null;default:0"`
	Source      string  `gorm:"type:varchar(50);not null;default:'manual';index"`

	SeoDomain *SeoDomain `gorm:"foreignKey:SeoDomainID;references:ID"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (SeoTrend) TableName() string {
	return "seo_trend"
}
