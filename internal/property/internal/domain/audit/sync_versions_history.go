package domain

import (
	"time"
)

type SyncVersionHistory struct {
	ID uint64 `gorm:"column:id;primaryKey" json:"id"`

	// sync context
	Resource string  `gorm:"column:resource;type:varchar(50);not null" json:"resource"`
	OwnerID  *int64  `gorm:"column:owner_id;type:bigint;index:idx_sync_history_owner,priority:1" json:"ownerId,omitempty"`
	SyncType *string `gorm:"column:sync_type;type:varchar(20)" json:"syncType,omitempty"` // check | delta | push

	// version info
	ClientVersion *int64 `gorm:"column:client_version;type:bigint" json:"clientVersion,omitempty"`
	ServerVersion *int64 `gorm:"column:server_version;type:bigint" json:"serverVersion,omitempty"`

	// stats
	ChangeCount    *int `gorm:"column:change_count;type:int" json:"changeCount,omitempty"`
	ResponseTimeMs *int `gorm:"column:response_time_ms;type:int" json:"responseTimeMs,omitempty"`
	DataSizeBytes  *int `gorm:"column:data_size_bytes;type:int" json:"dataSizeBytes,omitempty"`

	// client info
	ClientID *string `gorm:"column:client_id;type:varchar(128)" json:"clientId,omitempty"`
	DeviceID *string `gorm:"column:device_id;type:varchar(128)" json:"deviceId,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;autoCreateTime;index:idx_sync_history_owner,priority:2" json:"createdAt"`
}

func (SyncVersionHistory) TableName() string {
	return "sync_version_histories"
}
