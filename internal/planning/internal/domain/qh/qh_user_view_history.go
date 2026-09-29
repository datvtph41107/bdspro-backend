package qh_domain

import (
	"time"
	"tqd/internal/enums"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type QHUserViewHistory struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	UserID uint64 `gorm:"column:user_id;not null;index:idx_user_view_history_user_id;uniqueIndex:uidx_user_view_history_user_entity,priority:1,where:deleted_at IS NULL" json:"userId"`

	EntityType enums.WorkspaceEntityType `gorm:"column:entity_type;not null;index:idx_user_view_history_entity;uniqueIndex:uidx_user_view_history_user_entity,priority:2" json:"entityType"`
	EntityID   uint64                    `gorm:"column:entity_id;not null;index:idx_user_view_history_entity;uniqueIndex:uidx_user_view_history_user_entity,priority:3" json:"entityId"`

	ParcelID *uint64 `gorm:"column:parcel_id;index:idx_user_view_history_parcel_id" json:"parcelId,omitempty"`
	RegionID *uint64 `gorm:"column:region_id;index:idx_user_view_history_region_id" json:"regionId,omitempty"`

	Source enums.ViewHistorySource `gorm:"column:source;not null;default:99" json:"source"`

	ViewedAt      time.Time  `gorm:"column:viewed_at;not null;index:idx_user_view_history_viewed_at" json:"viewedAt"`
	FirstViewedAt *time.Time `gorm:"column:first_viewed_at" json:"firstViewedAt,omitempty"`

	ViewCount        uint64 `gorm:"column:view_count;not null;default:0" json:"viewCount"`
	CountedViewCount uint64 `gorm:"column:counted_view_count;not null;default:0" json:"countedViewCount"`

	LastZoom      float64 `gorm:"column:last_zoom;default:0" json:"lastZoom"`
	LastCenterLat float64 `gorm:"column:last_center_lat;default:0" json:"lastCenterLat"`
	LastCenterLon float64 `gorm:"column:last_center_lon;default:0" json:"lastCenterLon"`

	ViewportMinLon float64 `gorm:"column:viewport_min_lon;default:0" json:"viewportMinLon"`
	ViewportMinLat float64 `gorm:"column:viewport_min_lat;default:0" json:"viewportMinLat"`
	ViewportMaxLon float64 `gorm:"column:viewport_max_lon;default:0" json:"viewportMaxLon"`
	ViewportMaxLat float64 `gorm:"column:viewport_max_lat;default:0" json:"viewportMaxLat"`

	LastClientEventID string `gorm:"column:last_client_event_id;type:varchar(100);default:''" json:"lastClientEventId"`
	LastVisibleMs     uint64 `gorm:"column:last_visible_ms;not null;default:0" json:"lastVisibleMs"`

	Metadata datatypes.JSON `gorm:"column:metadata;type:jsonb;default:'{}'::jsonb" json:"metadata"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
}

func (QHUserViewHistory) TableName() string {
	return "user_view_history"
}
