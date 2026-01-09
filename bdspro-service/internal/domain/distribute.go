package domain

import (
	"bdspro/internal/enums"
	"time"
)

type DistributionPolicy struct {
	Price           int64                `json:"price"`
	CanDeal         bool                 `json:"canDeal"`
	ChannelPrice    bool                 `json:"channelPrice"`
	CommissionType  enums.CommissionType `json:"commissionType"`
	CommissionValue float64              `json:"commissionValue"`
	Note            string               `json:"note,omitempty"`
}

type DistributionEntity struct {
	ID                uint64                   `gorm:"primaryKey" json:"id"`
	ProductID         uint64                   `gorm:"not null;index" json:"productId"`
	Policy            DistributionPolicy       `gorm:"type:jsonb;not null" json:"policy"`
	DistributedAt     *time.Time               `gorm:"not null" json:"distributedAt"`
	RevokedAt         *time.Time               `json:"revokedAt,omitempty"`
	Status            enums.DistributionStatus `gorm:"not null" json:"status"`
	Reason            enums.DistributionReason `json:"reason"`                               // enum
	ReasonOther       string                   `gorm:"type:varchar(100)" json:"reasonOther"` // text bổ sung nếu Reason = OTHER
	FromDate          *time.Time               `gorm:"column:from_date;not null" json:"fromDate"`
	ToDate            *time.Time               `gorm:"column:to_date;not null" json:"toDate"`
	PartnerCount      int32                    `gorm:"-" json:"partnerCount,omitempty"`
	PartnerIds        []uint64                 `gorm:"-" json:"partnerIds,omitempty"`
	PartnerCommission float64                  `gorm:"column:partner_commission;not null;default:0" json:"partnerCommission"` // Mới
}

func (DistributionEntity) TableName() string {
	return "distributions"
}
