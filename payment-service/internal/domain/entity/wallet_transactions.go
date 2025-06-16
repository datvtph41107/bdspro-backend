package entity

import "time"

type WalletTransaction struct {
	Id       uint32
	WalletId uint32
	Type     string
	Amount   float64
	Status   string
	RelatedService string
	RelatedId string
	ExternalPaymentId string
	CreatedAt time.Time
	UpdatedAt time.Time
}