package domain

import (
	_models "common/models"
	"time"
)

type ContactProductEntity struct {
	_models.BaseEntityNotId
	ContactID       uint64  `gorm:"primaryKey" json:"contact_id"`
	ProductID       uint64  `gorm:"primaryKey" json:"product_id"`
	OriginProfileID *uint64 `gorm:"index"`
	// reorderPinAt
	PriorityPinAt *time.Time `gorm:"column:priority_pin_at;index"`
	// Tác động cuối như gọi trường này sinh ra dể reorder lại list
	LastInteractionAt *time.Time `gorm:"column:last_interaction_at;index"`
}

func (ContactProductEntity) TableName() string {
	return "tb_contact_product"
}