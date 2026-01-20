package domain

import (
	"bdspro/internal/enums"
	"time"

	"gorm.io/gorm"
)

type ProductPrice struct {
	ID                 uint64          `gorm:"primaryKey" json:"id"`
	CreatedAt          *time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt          *time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt          *gorm.DeletedAt `gorm:"index" json:"-"` // Hỗ trợ soft delete
	ProductID          *uint64         `json:"productId,omitempty"`
	Currency           string          `json:"currency,omitempty"`
	SalePrice          *float64        `gorm:"type:DECIMAL(15,0)" json:"salePrice,omitempty"`
	SaleCommission     *float64        `gorm:"type:DECIMAL(15,0)" json:"saleCommission,omitempty"`
	SaleCommissionType *uint32         `json:"saleCommisionType,omitempty"`
	Deposite           *float64        `gorm:"type:DECIMAL(15,0)" json:"deposite,omitempty"`
	RentPrice          *float64        `gorm:"type:DECIMAL(15,0)" json:"rentPrice,omitempty"`
	RentCommission     *float64        `gorm:"type:DECIMAL(15,0)" json:"rentCommission,omitempty"`
	RentCommissionType *uint32         `json:"rentCommisionType,omitempty"`
	RentPaymentCycle   *uint32         `json:"rentPaymentCycle,omitempty"`
	DistributeID       *uint64         `json:"distributeId,omitempty"`                                 // id phân phối ban đầu
	ChannelPrice       bool            `gorm:"column:channel_price;default:false" json:"channelPrice"` // Giá kênh

	PriceStatus enums.ProductPriceStatus `json:"status,omitempty"`
	ChangeNote  string                   `gorm:"type:text" json:"note,omitempty"`
	ChangedBy   uint64                   `json:"changedBy"`
	// PriceOwner       *float64 `json:"priceOwner,omitempty"` // deprecated
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
