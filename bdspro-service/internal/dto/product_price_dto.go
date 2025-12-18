package dto

import (
	_dto "common/domain/dto"
)

type PriceHistorySearch struct {
	_dto.Pagable
	ProductID uint64 `json:"productId"`
}

type PriceHistoryResponse struct {
	ID                 uint64   `json:"id"`
	Currency           string   `json:"currency"`
	SalePrice          *float64 `json:"salePrice"`
	SaleCommission     *float64 `json:"saleCommission"`
	SaleCommissionType *uint32  `json:"saleCommissionType"`
	Deposite           *float64 `json:"deposite"`
	RentPrice          *float64 `json:"rentPrice"`
	RentCommission     *float64 `json:"rentCommission"`
	RentCommissionType *uint32  `json:"rentCommissionType"`
	RentPaymentCycle   *uint32  `json:"rentPaymentCycle"`
	CreatedAt          string   `json:"createdAt"`
	UpdatedAt          string   `json:"updatedAt"`
}
