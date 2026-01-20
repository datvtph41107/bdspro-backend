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
	IsOwner      bool
	ProductUser  *domain.ProductUser
	Distribution *domain.DistributionEntity
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
