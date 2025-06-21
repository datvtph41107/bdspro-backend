package entity

import "time"

type Wallet struct {
	Id uint32

	CreatedAt time.Time
	UpdatedAt time.Time

	UserId   uint32
	Balance  float64
	Currency string
}
