package domain

import (
	_models "common/models"
	"crm/internal/enums"
)

type ContactEntity struct {
	_models.BaseEntity
	// OwnerID   uint64  `gorm:"column:owner_id" json:"ownerId"` => createdBy
	Phone     string           `gorm:"column:phone" json:"phone"`
	FullName  string           `gorm:"column:full_name" json:"fullName" binding:"required"`
	OwnerID   uint64           `gorm:"column:owner_id" json:"ownerId"`
	OwnerType enums.EOwnerType `gorm:"column:owner_type;default:10" json:"ownerType" binding:"required"`
}

func (ContactEntity) TableName() string {
	return "db_contact"
}
