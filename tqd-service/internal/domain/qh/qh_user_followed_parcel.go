package qh_domain

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type QHUserFollowedParcel struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	UserID uint64 `gorm:"column:user_id;not null;index:idx_user_followed_parcels_user_id;uniqueIndex:uidx_user_followed_parcels_active_user_parcel,priority:1,where:deleted_at IS NULL" json:"userId"`

	ParcelID uint64 `gorm:"column:parcel_id;not null;index:idx_user_followed_parcels_parcel_id;uniqueIndex:uidx_user_followed_parcels_active_user_parcel,priority:2,where:deleted_at IS NULL" json:"parcelId"`

	Note string `gorm:"column:note;type:text;default:''" json:"note"`

	Metadata datatypes.JSON `gorm:"column:metadata;type:jsonb;default:'{}'::jsonb" json:"metadata"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
}

func (QHUserFollowedParcel) TableName() string {
	return "user_followed_parcels"
}
