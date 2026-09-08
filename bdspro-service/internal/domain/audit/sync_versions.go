package domain

import "time"

type SyncVersion struct {
	Resource   string    `gorm:"column:resource;type:varchar(50);not null;primaryKey" json:"resource"`
	OwnerID    int64     `gorm:"column:owner_id;type:bigint;not null;primaryKey" json:"ownerId"`
	MaxVersion int64     `gorm:"column:max_version;type:bigint;not null;default:0" json:"maxVersion"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamptz;not null;autoUpdateTime" json:"updatedAt"`
}

func (SyncVersion) TableName() string {
	return "sync_versions"
}
