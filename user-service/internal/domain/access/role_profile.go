package access

import (
	_models "common/models"
)

type RoleProfile struct {
	_models.BaseEntity
	ProfileID uint64 `gorm:"index;not null"`
	RoleID    uint64 `gorm:"index;not null"`
}

func (RoleProfile) TableName() string {
	return "role_profiles"
}
