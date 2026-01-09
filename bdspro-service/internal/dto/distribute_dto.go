package dto

import (
	"bdspro/internal/enums"
	_dto "common/domain/dto"
	"errors"
	"time"
)

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
	_dto.Pagable            // page, limit, offset nếu có
	FromDate     *time.Time `form:"fromDate"`     // lọc FromDate >=
	ToDate       *time.Time `form:"toDate"`       // lọc ToDate <=
	CanDeal      *bool      `form:"canDeal"`      // lọc policy.CanDeal
	ChannelPrice *bool      `form:"channelPrice"` // lọc policy.ChannelPrice
}

type DistributeResultDTO struct {
	CreatedCount int        `json:"createdCount"`
	FromDate     *time.Time `json:"fromDate"`
	ToDate       *time.Time `json:"toDate"`
}

type RevokeDTO struct {
	Reason      uint32
	ReasonOther string
}

func (d *DistributeDTO) Validate() error {
	now := time.Now()
	if d.ProductID == 0 {
		return errors.New("ProductID must be greater than 0")
	}
	if d.Price < 0 {
		return errors.New("Policy.Price must be >= 0")
	}
	if len(d.PartnerIDs) == 0 {
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
