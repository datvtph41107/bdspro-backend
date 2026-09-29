package domain

import (
	_models "common/domain/entity"
	"crm/internal/enums"
)

type OriginProfile struct {
	_models.BaseEntityNotId
	OriginID    uint64                   `gorm:"column:origin_id;primaryKey" json:"originId"`
	DisplayName string                   `gorm:"column:display_name" json:"displayName"`
	Phone       string                   `gorm:"column:phone" json:"phone"`
	Avatar      string                   `gorm:"column:avatar" json:"avatar"`
	OwnerID     uint64                   `gorm:"column:owner_id"` // references cua UserId
	OwnerOf     enums.EOriginProfileType `gorm:"column:owner_of;not null;default:10" json:"ownerOf"`
}

func (OriginProfile) TableName() string {
	return "origin_profile"
}