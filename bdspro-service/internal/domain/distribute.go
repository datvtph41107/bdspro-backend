package domain

import (
	"bdspro/internal/enums"
	"time"
)

// type DistributionPolicy struct {
// 	Price           int64                `json:"price"`
// 	CanDeal         bool                 `json:"canDeal"`
// 	ChannelPrice    bool                 `json:"channelPrice"`
// 	CommissionType  enums.CommissionType `json:"commissionType"`
// 	CommissionValue float64              `json:"commissionValue"`
// 	Note            string               `json:"note,omitempty"`
// }

type DistributionEntity struct {
	ID                uint64                   `gorm:"primaryKey" json:"id"`
	Price             int64                    `gorm:"column:price" json:"price"`
	CanDeal           bool                     `gorm:"column:can_deal;default:false" json:"canDeal"`
	ChannelPrice      bool                     `gorm:"column:channel_price;default:false" json:"channelPrice"`
	CommissionType    enums.CommissionType     `gorm:"column:commission_type;default:10" json:"commissionType"`
	CommissionValue   float64                  `gorm:"column:commission_value;default:0" json:"commissionValue"`
	Note              string                   `gorm:"column:note" json:"note,omitempty"`
	ProductID         uint64                   `gorm:"not null;index" json:"productId"`
	DistributedAt     *time.Time               `gorm:"not null" json:"distributedAt"`
	RevokedAt         *time.Time               `json:"revokedAt,omitempty"`
	Status            enums.DistributionStatus `gorm:"not null" json:"status"`
	Reason            enums.DistributionReason `json:"reason"`
	ReasonOther       string                   `gorm:"type:varchar(100)" json:"reasonOther"`
	FromDate          *time.Time               `gorm:"column:from_date;not null" json:"fromDate"`
	ToDate            *time.Time               `gorm:"column:to_date;not null" json:"toDate"`
	PartnerCount      int32                    `gorm:"-" json:"partnerCount,omitempty"`
	PartnerIds        []uint64                 `gorm:"-" json:"partnerIds,omitempty"`
	PartnerCommission float64                  `gorm:"column:partner_commission;not null;default:0" json:"partnerCommission"`
	CreatedAt         *time.Time               `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt         *time.Time               `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (DistributionEntity) TableName() string {
	return "distributions"
}
