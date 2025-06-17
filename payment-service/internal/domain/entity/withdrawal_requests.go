package entity

import (
	paymentpb "pb/types/payment"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type WithdrawalRequest struct {
	Id              uint32
	WalletId        uint32
	Amount          float64
	Status          string
	PaymentMethod   string
	BankAccountInfo string
	RequestedBy     uint32
	ApprovedBy      uint32
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (w *WithdrawalRequest) ToProto() *paymentpb.WithdrawalRequest {
	return &paymentpb.WithdrawalRequest{
		Id:              w.Id,
		WalletId:        w.WalletId,
		Amount:          w.Amount,
		Status:          w.Status,
		PaymentMethod:   w.PaymentMethod,
		BankAccountInfo: w.BankAccountInfo,
		RequestedBy:     w.RequestedBy,
		ApprovedBy:     w.ApprovedBy,
		CreatedAt:       timestamppb.New(w.CreatedAt),
		UpdatedAt:       timestamppb.New(w.UpdatedAt),
	}
}
