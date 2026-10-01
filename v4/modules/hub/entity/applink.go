package entity

import "time"

// Applink maps a deep-link code (base64-encoded uint64) to a refID + action.
type Applink struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	Code      uint64    `gorm:"column:code;type:bigint;uniqueIndex:idx_applink_code"`
	RefID     uint64    `gorm:"column:ref_id;index:idx_applink_ref_action"`
	Action    int32     `gorm:"column:action;type:int;index:idx_applink_ref_action"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Applink) TableName() string { return "applink" }
