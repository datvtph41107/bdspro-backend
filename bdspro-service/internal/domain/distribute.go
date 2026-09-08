package domain

import (
	"bdspro/internal/enums"
	"time"
)

type DistributionEntity struct {
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
	PriceID       *uint64                  `json:"priceId,omitempty"` // id của product_price
	ProductPrice  *ProductPrice            `gorm:"foreignKey:PriceID;references:ID" json:"productPrice,omitempty"`

	// PartnerCommission float64                  `gorm:"column:partner_commission;not null;default:0" json:"partnerCommission"`
	// Price             int64                    `gorm:"column:price" json:"price"`
	// CommissionType    enums.CommissionType     `gorm:"column:commission_type;default:10" json:"commissionType"`
	// CommissionValue   float64                  `gorm:"column:commission_value;default:0" json:"commissionValue"`
}

func (DistributionEntity) TableName() string {
	return "distributions"
}
