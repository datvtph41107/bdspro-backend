package dto

import "time"

// OrderType represents the type of order
type OrderType int8

const (
	OrderBill       OrderType = 1
	OrderMembership OrderType = 2
)

// OrderRequest represents the request to create a payment order
type OrderRequest struct {
	Amount      int       `json:"amount"`
	Description string    `json:"description"`
	TargetID    *uint64   `json:"targetId,omitempty"`
	Type        OrderType `json:"type"`
}

// OrderResponse represents the response from payment service
type OrderResponse struct {
	PaymentURL string `json:"paymentUrl"`
}

// VerifyOrderDTO represents the webhook notification from payment service
type VerifyOrderDTO struct {
	TargetID  uint64     `json:"targetId" binding:"required"`
	Type      OrderType  `json:"type" binding:"required"`
	Signature string     `json:"signature" binding:"required"`
	PaidAt    *time.Time `json:"paidAt" binding:"required"`
}
