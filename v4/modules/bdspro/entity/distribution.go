package entity

import (
	"time"

	"mdn4/modules/bdspro/enums"
)

// Distribution maps to table distributions (product channel distribution).
// ORM::TABLE
type Distribution struct {
	ID            uint64                   `gorm:"primaryKey" json:"id"`
	CanDeal       bool                     `gorm:"column:can_deal;default:false" json:"canDeal"`
	Note          string                   `gorm:"column:note" json:"note,omitempty"`
	ProductID     uint64                   `gorm:"not null;index" json:"productId"`
	DistributedAt *time.Time               `gorm:"not null" json:"distributedAt"`
	RevokedAt     *time.Time               `json:"revokedAt,omitempty"`
	Reason        enums.DistributionReason `json:"reason"`
	ReasonOther   string                   `gorm:"type:varchar(100)" json:"reasonOther"`
	FromDate      *time.Time               `gorm:"column:from_date;not null" json:"fromDate"`
	ToDate        *time.Time               `gorm:"column:to_date;not null" json:"toDate"`
	CreatedAt     *time.Time               `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt     *time.Time               `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	PriceID       *uint64                  `json:"priceId,omitempty"`
	ProductPrice  *ProductPrice            `gorm:"foreignKey:PriceID;references:ID" json:"productPrice,omitempty"`
}

func (Distribution) TableName() string {
	return "distributions"
}
