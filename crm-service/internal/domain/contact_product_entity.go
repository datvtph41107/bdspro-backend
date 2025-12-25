package domain

import (
	_models "common/models"
)

type ContactProductEntity struct {
	_models.BaseEntityNotId
	ContactID uint64 `gorm:"primaryKey" json:"contact_id"`
	ProductID uint64 `gorm:"primaryKey" json:"product_id"`
}

func (ContactProductEntity) TableName() string {
	return "tb_contact_product"
}
