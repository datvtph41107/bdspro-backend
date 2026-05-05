package qh_domain

import "time"

type QHRegionLabel struct {
	RegionID uint64 `gorm:"primaryKey" json:"regionId"`
	LabelID  uint64 `gorm:"primaryKey" json:"labelId"`

	// Thông tin ai gán, khi nào gán
	AssignedBy uint64    `gorm:"default:0" json:"assignedBy"`
	AssignedAt time.Time `gorm:"autoCreateTime" json:"assignedAt"`

	// Nguồn gốc của việc gán label này
	SourceFile string `gorm:"type:varchar(255)" json:"sourceFile,omitempty"`
	Notes      string `gorm:"type:text" json:"notes,omitempty"`
}

func (QHRegionLabel) TableName() string { return "qh_region_labels" }
