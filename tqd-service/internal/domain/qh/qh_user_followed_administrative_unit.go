package qh_domain

import "time"

// QHUserFollowedAdministrativeUnit stores a user follow of a province/ward.
// Identity is compositeId ({sourceId}{1=tỉnh|2=phường}). Unfollow hard-deletes the row.
type QHUserFollowedAdministrativeUnit struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint64    `gorm:"column:user_id;not null;uniqueIndex:uidx_user_followed_adm_unit,priority:1" json:"userId"`
	CompositeID uint64    `gorm:"column:composite_id;not null;uniqueIndex:uidx_user_followed_adm_unit,priority:2" json:"compositeId"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (QHUserFollowedAdministrativeUnit) TableName() string {
	return "user_followed_administrative_units"
}
