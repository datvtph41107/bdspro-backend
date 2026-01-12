package domain

import (
	_models "common/domain/entity"
	"crm/internal/enums"
)

type OriginProfile struct {
	_models.BaseEntity
	OwnerOf         enums.EOriginProfileType `gorm:"column:type;not null" json:"ownerOf"`
	DisplayName     string                   `gorm:"column:display_name" json:"displayName"`
	Phone           string                   `gorm:"column:phone" json:"phone"`
	Avatar          string                   `gorm:"column:avatar" json:"avatar"`
	OriginProfileID *uint64                  `gorm:"column:origin_profile_id"`
}

func (OriginProfile) TableName() string {
	return "origin_profile"
}
