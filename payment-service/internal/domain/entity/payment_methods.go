package entity

import (
	paymentpb "pb/types/payment"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type PaymentMethod struct {
	Id             uint32
	OrganizationId uint32
	Name           string
	Code           string
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (m *PaymentMethod) ToProto() *paymentpb.PaymentMethod {
	return &paymentpb.PaymentMethod{
		Id:             m.Id,
		OrganizationId: m.OrganizationId,
		Name:           m.Name,
		Code:           m.Code,
		IsActive:       m.IsActive,
		CreatedAt:      timestamppb.New(m.CreatedAt),
		UpdatedAt:      timestamppb.New(m.UpdatedAt),
	}
}