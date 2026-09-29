package entity

import (
	_models "common/models"
	"organization/internal/enums"
	"time"
)

type DealMilestone struct {
	_models.BaseEntity

	DealID        uint64              `gorm:"not null;index:idx_deal_milestones_deal_id"`
	Title         string              `gorm:"not null"`
	Description   string
	ExpectedDate  *time.Time
	CompletedDate *time.Time
	Status        enums.MilestoneStatus `gorm:"not null;default:10"`
	OrderIndex    uint32                `gorm:"not null;default:0"`
	Deal          *Deal                 `gorm:"foreignKey:DealID;references:ID"`
}

func (DealMilestone) TableName() string {
	return "deal_milestones"
} 