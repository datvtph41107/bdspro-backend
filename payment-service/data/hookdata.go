package data

import (
	"payment/models"
	"time"
)

type NotiComplete struct {
	ID        uint64
	ProfileId uint64
	OrderId   uint64
	Type      models.OrderType
	TargetId  uint64
	PaidAt    *time.Time
	Signature string
}
