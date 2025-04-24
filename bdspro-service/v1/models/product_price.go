package models

import _models "common/models"

type ProductPrice struct {
	_models.BaseEntity
	ProductID        *uint64  `json:"productId,omitempty"`
	Currency         string   `json:"currency,omitempty"`
	PriceOwner       uint64   `json:"priceOwner,omitempty"`
	SalePrice        *float64 `gorm:"type:DECIMAL(15,0)" json:"salePrice,omitempty"`
	SaleCommission   uint64   `json:"saleCommission,omitempty"`
	Deposite         *float64 `gorm:"type:DECIMAL(15,0)" json:"deposite,omitempty"`
	RentPrice        *float64 `gorm:"type:DECIMAL(15,0)" json:"rentPrice,omitempty"`
	RentCommission   float64  `json:"rentCommission,omitempty"`
	RentPaymentCycle *uint    `json:"rentPaymentCycle,omitempty"`
	// SaleCommissionType uint     `json:"saleCommissionType,omitempty"`
	// RentCommissionType uint     `json:"rentCommissionType,omitempty"`
	// Chu kỳ thanh toán ("tháng", "quý"...)

	// isCurrent          bool
	// PriceType
}

func (ProductPrice) TableName() string {
	return "product_price"
}

// type ProductPrice struct {
// 	_models.BaseEntity
// 	ProductId          uint64
// 	Value              uint64
// 	Currency           string
// 	PriceOwner         uint64
// 	SaleCommission     uint64
// 	SaleCommissionType uint
// 	RentCommission     uint64
// 	RentCommissionType uint
// 	// isCurrent          bool
// 	// PriceType
// }

// func (ProductPrice) TableName() string {
// 	return "product_price"
// }
