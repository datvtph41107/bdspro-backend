package entity

import "time"

type DealStatus string

const (
	DealStatusPending DealStatus = "pending"
	DealStatusActive DealStatus = "active"
	DealStatusInactive DealStatus = "inactive"
)

type Deal struct {
	Id uint32
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uint32
	UpdatedBy uint32

	GroupId uint32
	Name string
	Amount float64
	Status DealStatus
}