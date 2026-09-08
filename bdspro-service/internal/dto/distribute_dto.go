package dto

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	_dto "common/domain/dto"
	"errors"
	"time"

	"github.com/lib/pq"
)

// type PartnerInfo struct {
// 	ID      uint64         `json:"id"`
// 	OwnerOf enums.EOwnerOf `json:"ownerOf"`
// }
//  Partners        []PartnerInfo `json:"partners"`

type DistributeDTO struct {
	ProductID       uint64               `json:"productId"`
	PartnerIDs      []uint64             `json:"partnerIds"`
	FromDate        *time.Time           `json:"fromDate,omitempty"`
	ToDate          *time.Time           `json:"toDate,omitempty"`
	Price           int64                `json:"price"`
	CanDeal         bool                 `json:"canDeal"`
	ChannelPrice    bool                 `json:"channelPrice"`
	CommissionType  enums.CommissionType `json:"commissionType"`
	CommissionValue float64              `json:"commissionValue"`
	Note            string               `json:"note,omitempty"`
}

type FilterDistributeDTO struct {
	_dto.Pagable
	FromDate     *time.Time                 `form:"fromDate"`
	ToDate       *time.Time                 `form:"toDate"`
	CanDeal      *bool                      `form:"canDeal"`
	ChannelPrice *bool                      `form:"channelPrice"`
	Text         string                     `form:"query"`
	Status       []enums.DistributionStatus `form:"status"`
	Sort         string                     `form:"sort"`
}

type DistributionWithPartners struct {
	domain.DistributionEntity
	PartnerIds   pq.Int64Array            `gorm:"column:partner_ids;type:bigint[]" json:"partnerIds"`
	PartnerCount int32                    `gorm:"column:partner_count" json:"partnerCount"`
	Status       enums.DistributionStatus `json:"status"`

	Price            int64                `gorm:"column:price" json:"price"`
	CommissionType   enums.CommissionType `gorm:"column:commission_type" json:"commissionType"`
	CommissionValue  float64              `gorm:"column:commission_value" json:"commissionValue"`
	ProductUpdatedAt *time.Time           `gorm:"column:product_updated_at;->;-:migration" json:"productUpdatedAt"`
}

type DistributeResultDTO struct {
	Price             int64                    `json:"price"`
	CanDeal           bool                     `json:"canDeal"`
	ChannelPrice      bool                     `json:"channelPrice"`
	Status            enums.DistributionStatus `json:"status"`
	CommissionType    enums.CommissionType     `json:"commissionType"`
	CommissionValue   float64                  `json:"commissionValue"`
	Note              string                   `json:"note"`
	PartnerCommission float64                  `json:"partnerCommission"`
	CreatedCount      int                      `json:"createdCount"`
	FromDate          *time.Time               `json:"fromDate"`
	ToDate            *time.Time               `json:"toDate"`
	ProductPrice      *ProductPriceDTO         `json:"productPrice"`
}

type RevokeDTO struct {
	Reason      uint32
	ReasonOther string
}

func (d *DistributeDTO) Validate(isUpdate bool) error {
	now := time.Now()
	if d.ProductID == 0 {
		return errors.New("ProductID must be greater than 0")
	}
	if d.Price < 0 {
		return errors.New("Price must be >= 0")
	}
	if len(d.PartnerIDs) == 0 && !isUpdate {
		return errors.New("PartnerIDs cannot be empty")
	}
	if d.FromDate == nil {
		d.FromDate = &now
	}
	if d.ToDate == nil {
		t := now.AddDate(0, 1, 0)
		d.ToDate = &t
	}
	if d.FromDate.After(*d.ToDate) {
		return errors.New("FromDate cannot be after ToDate")
	}

	return nil
}
