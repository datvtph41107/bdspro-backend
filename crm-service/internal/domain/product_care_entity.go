package domain

import (
	_models "common/models"
)

type ProductCareEntity struct {
	_models.BaseEntityNotId
	LeadID    uint64 `gorm:"primaryKey" json:"lead_id"`
	ProductID uint64 `gorm:"primaryKey" json:"product_id"`
}

func (ProductCareEntity) TableName() string {
	return "tb_product_care"
}