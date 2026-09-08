package wallet

import (
	"time"
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
