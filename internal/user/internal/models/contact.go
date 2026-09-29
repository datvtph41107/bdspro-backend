package models

import (
	_models "common/models"
)

type ContactEntity struct {
	_models.BaseEntity
	// OwnerID   uint64  `gorm:"column:owner_id" json:"ownerId"` => createdBy
	ProfileID *uint64 `gorm:"column:profile_id" json:"profileId,omitempty"` // Có thể NULL, dùng "omitempty" để bỏ qua nếu giá trị là nil
	Phone     string  `gorm:"column:phone" json:"phone"`
	FullName  string  `gorm:"column:full_name" json:"fullName"`
}

func (ContactEntity) TableName() string {
	return "tb_contact"
}
