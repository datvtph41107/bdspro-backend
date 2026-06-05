package qh_domain

import (
	"time"
	"tqd/internal/enums"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type QHUserViewHistory struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	UserID uint64 `gorm:"column:user_id;not null;index:idx_user_view_history_user_id" json:"userId"`

	EntityType enums.WorkspaceEntityType `gorm:"column:entity_type;not null;index:idx_user_view_history_entity" json:"entityType"`
	EntityID   uint64                    `gorm:"column:entity_id;not null;index:idx_user_view_history_entity" json:"entityId"`

	ParcelID *uint64 `gorm:"column:parcel_id;index:idx_user_view_history_parcel_id" json:"parcelId,omitempty"`

	Source enums.ViewHistorySource `gorm:"column:source;not null;default:99" json:"source"`

	ViewedAt time.Time `gorm:"column:viewed_at;not null;index:idx_user_view_history_viewed_at" json:"viewedAt"`

	Metadata datatypes.JSON `gorm:"column:metadata;type:jsonb;default:'{}'::jsonb" json:"metadata"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
}

func (QHUserViewHistory) TableName() string {
	return "user_view_history"
}
