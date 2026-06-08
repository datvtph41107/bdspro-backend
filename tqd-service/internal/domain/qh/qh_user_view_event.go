package qh_domain

import (
	"time"
	"tqd/internal/enums"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type QHUserViewEvent struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	UserID uint64 `gorm:"column:user_id;not null;index:idx_user_view_events_user_id;uniqueIndex:uidx_user_view_event_client_event,priority:1,where:client_event_id <> '' AND deleted_at IS NULL" json:"userId"`

	ClientEventID string `gorm:"column:client_event_id;type:varchar(100);default:'';uniqueIndex:uidx_user_view_event_client_event,priority:2" json:"clientEventId"`

	SessionID string `gorm:"column:session_id;type:varchar(100);default:'';index:idx_user_view_events_session_id" json:"sessionId"`
	DeviceID  string `gorm:"column:device_id;type:varchar(100);default:'';index:idx_user_view_events_device_id" json:"deviceId"`

	EntityType enums.WorkspaceEntityType `gorm:"column:entity_type;not null;index:idx_user_view_events_entity" json:"entityType"`
	EntityID   uint64                    `gorm:"column:entity_id;not null;index:idx_user_view_events_entity" json:"entityId"`

	ParcelID *uint64 `gorm:"column:parcel_id;index:idx_user_view_events_parcel_id" json:"parcelId,omitempty"`
	RegionID *uint64 `gorm:"column:region_id;index:idx_user_view_events_region_id" json:"regionId,omitempty"`

	Source    enums.ViewHistorySource `gorm:"column:source;not null;default:99" json:"source"`
	SourceRef string                  `gorm:"column:source_ref;type:varchar(255);default:''" json:"sourceRef"`
	RouteName string                  `gorm:"column:route_name;type:varchar(100);default:''" json:"routeName"`

	Zoom      float64 `gorm:"column:zoom;default:0" json:"zoom"`
	CenterLat float64 `gorm:"column:center_lat;default:0" json:"centerLat"`
	CenterLon float64 `gorm:"column:center_lon;default:0" json:"centerLon"`

	ViewportMinLon float64 `gorm:"column:viewport_min_lon;default:0" json:"viewportMinLon"`
	ViewportMinLat float64 `gorm:"column:viewport_min_lat;default:0" json:"viewportMinLat"`
	ViewportMaxLon float64 `gorm:"column:viewport_max_lon;default:0" json:"viewportMaxLon"`
	ViewportMaxLat float64 `gorm:"column:viewport_max_lat;default:0" json:"viewportMaxLat"`

	VisibleMs   uint64 `gorm:"column:visible_ms;not null;default:0" json:"visibleMs"`
	CountIntent bool   `gorm:"column:count_intent;not null;default:false" json:"countIntent"`
	Counted     bool   `gorm:"column:counted;not null;default:false;index:idx_user_view_events_counted" json:"counted"`

	DedupeKey string `gorm:"column:dedupe_key;type:varchar(255);default:'';index:idx_user_view_events_dedupe_key" json:"dedupeKey"`

	Metadata datatypes.JSON `gorm:"column:metadata;type:jsonb;default:'{}'::jsonb" json:"metadata"`

	ViewedAt time.Time `gorm:"column:viewed_at;not null;index:idx_user_view_events_viewed_at" json:"viewedAt"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
}

func (QHUserViewEvent) TableName() string {
	return "user_view_events"
}
