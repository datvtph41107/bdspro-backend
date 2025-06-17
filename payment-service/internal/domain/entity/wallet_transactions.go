package entity

import (
	paymentpb "pb/types/payment"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type WalletTransaction struct {
	Id                uint32
	WalletId         uint32
	Type             string
	Amount           float64
	Status           string
	RelatedService   string
	RelatedId        string
	ExternalPaymentId string
	PaymentMethod    string
	Description      string
	CreatedBy        uint32
	ApprovedBy       uint32
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (t *WalletTransaction) ToProto() *paymentpb.WalletTransaction {
	return &paymentpb.WalletTransaction{
		Id:               t.Id,
		WalletId:         t.WalletId,
		Type:             t.Type,
		Amount:           t.Amount,
		Status:           t.Status,
		RelatedService:   t.RelatedService,
		RelatedId:        t.RelatedId,
		ExternalPaymentId: t.ExternalPaymentId,
		CreatedAt:        timestamppb.New(t.CreatedAt),
		UpdatedAt:        timestamppb.New(t.UpdatedAt),
	}
}