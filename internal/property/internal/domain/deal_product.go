package domain

import _models "common/models"

type DealProduct struct {
	_models.BaseEntity
	DealID    uint64 `gorm:"not null;uniqueIndex:uq_deal_products,priority:1"`
	ProductID uint64 `gorm:"not null;uniqueIndex:uq_deal_products,priority:2"`
}

func (DealProduct) TableName() string {
	return "deal_products"
}
