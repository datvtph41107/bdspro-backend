package entity

import (
	_models "common/models"
	"organization/internal/enums"
	"time"
)

type Investment struct {
	_models.BaseEntity
	DealID         uint64               `gorm:"not null"`
	InvestType     enums.DealMemberType `gorm:"not null;default:10"`
	Amount         float64              `gorm:"not null"`
	TransferTime   *time.Time           ``
	Note           string               `gorm:"not null"`
	AttachDocument []AttachDocument     `gorm:"foreignKey:OwnerId;references:ID"`

	Status    enums.ApprovedStatus `gorm:"not null;default:10"`
	ChangedAt *time.Time           ``
}

func (Investment) TableName() string {
	return "investments"
}
