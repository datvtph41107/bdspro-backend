package dto

import (
	"bdspro/internal/enums"
	_dto "common/domain/dto"
)

type PriceHistorySearch struct {
	_dto.Pagable
	ProductID       uint64 `json:"productId"`
	TransactionType *enums.TransactionType
}

type PriceHistoryItemDTO struct {
	ID              uint64                `json:"id"`
	ProductID       uint64                `json:"productId"`
	TransactionType enums.TransactionType `json:"transactionType"`
	Currency        string                `json:"currency"`
	Old             *PriceDTO             `json:"old,omitempty"`
	New             *PriceDTO             `json:"new"`
	CreatedAt       string                `json:"createdAt"`
	ChangeNote      string                `json:"changeNote"`
}
type PriceDTO struct {
	Sale *SalePriceDTO `json:"sale,omitempty"`
	Rent *RentPriceDTO `json:"rent,omitempty"`
}

type SalePriceDTO struct {
	Price          *float64 `json:"price"`
	Commission     *float64 `json:"commission"`
	CommissionType *uint32  `json:"commissionType"`
	Deposite       *float64 `json:"deposite"`
}

type RentPriceDTO struct {
	Price          *float64 `json:"price"`
	Commission     *float64 `json:"commission"`
	CommissionType *uint32  `json:"commissionType"`
	PaymentCycle   *uint32  `json:"paymentCycle"`
}
