package qh_domain

import (
	_models "common/domain/entity"
)

type QHLabel struct {
	_models.BaseEntity
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	LayerID     uint64 `gorm:"type:bigint;not null;index" json:"layerId"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Color       string `gorm:"type:varchar(20);default:'#1890ff'" json:"color"`
	Description string `gorm:"type:varchar(255)" json:"description,omitempty"`
	SortOrder   int    `gorm:"default:0;index" json:"sortOrder"`
	// CreatedAt   time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	// UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	// DeletedAt   *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (QHLabel) TableName() string {
	return "qh_labels"
}
