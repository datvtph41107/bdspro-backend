package seo_domain

import (
	"time"

	"gorm.io/gorm"
)

const (
	SeoGenerationStatusPending = "pending"
	SeoGenerationStatusRunning = "running"
	SeoGenerationStatusSuccess = "success"
	SeoGenerationStatusFailed  = "failed"
)

// SeoGenerationLog lưu kết quả job generate HTML/sitemap để retry và debug vận hành.
type SeoGenerationLog struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	SeoDomainID uint64 `gorm:"column:seo_domain_id;not null;index"`
	Status      string `gorm:"type:varchar(30);not null;index"`

	ErrorMessage string     `gorm:"column:error_message;type:text"`
	StartedAt    *time.Time `gorm:"column:started_at;index"`
	FinishedAt   *time.Time `gorm:"column:finished_at;index"`

	SeoDomain *SeoDomain `gorm:"foreignKey:SeoDomainID;references:ID"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (SeoGenerationLog) TableName() string {
	return "seo_generation_log"
}
