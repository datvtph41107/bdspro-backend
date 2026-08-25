package seo_domain

import "time"

// SeoCount is the aggregate click counter per SEO page.
type SeoCount struct {
	SeoID     uint64    `gorm:"column:seo_id;primaryKey"`
	Count     int64     `gorm:"column:count;type:bigint;not null;default:0"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
}

func (SeoCount) TableName() string {
	return "seo_count"
}

// SeoCountLog stores one flushed click-count batch for an SEO page.
type SeoCountLog struct {
	SeoID     uint64    `gorm:"column:seo_id;type:bigint;not null;index"`
	Clicked   int64     `gorm:"column:clicked;type:bigint;not null;default:0"`
	Timestamp time.Time `gorm:"column:timestamp;type:timestamptz;not null;index"`
}

func (SeoCountLog) TableName() string {
	return "seo_count_log"
}
