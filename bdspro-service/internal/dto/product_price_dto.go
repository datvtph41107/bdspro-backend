package dto

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	_dto "common/domain/dto"
)

type PriceHistorySearch struct {
	_dto.Pagable
	ProductID       uint64 `json:"productId"`
	TransactionType *enums.TransactionType
}

type PriceUpdateContext struct {
	IsOwner        bool
	CanUpdatePrice bool
	ProductUser    *domain.ProductUser
	Distribution   *domain.DistributionEntity
	ProducutPrice  *domain.ProductPrice
	Product        *domain.Product
}

type PriceRow struct {
	ProductID     uint64  `gorm:"column:product_id"`
	OwnerID       uint64  `gorm:"column:owner_id"`
	LastPriceID   *uint64 `gorm:"column:last_price_id"`
	Area          float64 `gorm:"column:area"`
	ProductUserID *uint64 `gorm:"column:product_user_id"`
	DistributeID  *uint64 `gorm:"column:distribute_id"`
	PriceID       *uint64 `gorm:"column:price_id"`
	ChannelPrice  *bool   `gorm:"column:channel_price"`
}

type PriceBenefitItem struct {
	ActorType   string // OWNER | DISTRIBUTION | USER
	ActorID     uint64
	BuyPrice    float64
	SellPrice   float64
	Benefit     float64
	BenefitRate float64 // %
}

type PriceBenefitDetail struct {
	ProductID     uint64
	Chain         []PriceBenefitItem
	SystemBenefit float64
}

type PriceHistoryItemDTO struct {
	ID              uint64                `json:"id"`
	ProductID       uint64                `json:"productId"`
	TransactionType enums.TransactionType `json:"transactionType"`
	Currency        string                `json:"currency"`
	Old             *PriceDTO             `json:"old,omitempty"`
	New             *PriceDTO             `json:"new"`
	IsCurrent       bool                  `json:"isCurrent"`
	CreatedAt       string                `json:"createdAt"`
	ChangeNote      string                `json:"changeNote"`
	DistributeID    *uint64               `json:"distributeId,omitempty"`
}
type PriceDTO struct {
	Price          *float64     `json:"price"`
	Commission     *float64     `json:"commission"`
	CommissionType *uint32      `json:"commissionType"`
	CreatedBy      *UserInfoDTO `json:"createdBy,omitempty"`
}

// type SalePriceDTO struct {
// 	Price          *float64 `json:"price"`
// 	Commission     *float64 `json:"commission"`
// 	CommissionType *uint32  `json:"commissionType"`
// 	Deposite       *float64 `json:"deposite"`
// }

// type RentPriceDTO struct {
// 	Price          *float64 `json:"price"`
// 	Commission     *float64 `json:"commission"`
// 	CommissionType *uint32  `json:"commissionType"`
// 	PaymentCycle   *uint32  `json:"paymentCycle"`
// }

type UserInfoDTO struct {
	ID     uint64
	Name   string
	Avatar string
	Role   string
}

type ProductPriceDetailDTO struct {
	Price        *float64 `json:"price"`
	PricePerM2   int64    `json:"pricePerM2"`
	Currency     string   `json:"currency"`
	Cost         *float64 `json:"cost,omitempty"`
	Commission   *float64 `json:"commission,omitempty"`
	Channel      *bool    `json:"channel,omitempty"`
	ChannelValue *int64   `json:"channelValue,omitempty"`
	Profit       *float64 `json:"profit,omitempty"`
	Note         *string  `json:"note,omitempty"`
}
