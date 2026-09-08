package domain

import (
	"bdspro/internal/enums"
	"time"
)

// PaymentMethod represents a payment method
type PaymentMethod struct {
	Id         uint32
	Name       string
	MethodCode string
	Status     enums.PaymentStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
