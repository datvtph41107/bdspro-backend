package domain

import (
	_models "common/models"
	"time"
)

type PropertyUser struct {
	_models.BaseEntity
	OwnerOriginID uint64 `gorm:"column:owner_origin_id;not null;index:idx_property_user_user_profile,priority:1"`
	PropertyID    uint64 `gorm:"column:property_id;not null;index:idx_property_user_user_profile,priority:2"`

	OwnerAt time.Time `gorm:"column:is_owner;"`
	RoleID  *uint64   `gorm:"column:role_id;"`

	OriginProfileID *uint64 `gorm:"column:origin_id;type:bigint;index:idx_property_origin"`
}

func (PropertyUser) TableName() string {
	return "property_user"
}
