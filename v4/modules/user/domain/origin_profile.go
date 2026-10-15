package entity

import (
	_models "common/domain/entity"
	"userv4/modules/crm/enums"
)

// ORM::TABLE
// tạo người dùng/tổ chức/admin hệ thống thì tạo bản ghi này
// id của người dùng/tổ chức/admin hệ thống lấy theo originId trong đây
type OriginProfile struct {
	_models.BaseEntityNotId
	OriginID    uint64 `gorm:"column:origin_id;primaryKey" json:"originId"`
	DisplayName string `gorm:"column:display_name" json:"displayName"`
	Phone       string `gorm:"column:phone" json:"phone"`
	Avatar      string `gorm:"column:avatar" json:"avatar"`

	OwnerOf enums.EOwnerOf `gorm:"column:owner_of;not null;default:10" json:"ownerOf"`
	// là id trong UserProfile.ProfileID/Organization.ID
	// deprecated: sau có thể bỏ do id từ các bảng UserProfile/Organization
	// khi tạo lấy theo originId
	OwnerID uint64 `gorm:"column:owner_id"` // references cua UserId
}

func (OriginProfile) TableName() string {
	return "origin_profile"
}
