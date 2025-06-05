package dto

import (
	_dto "common/dto"
	"crm/internal/enums"
	"time"
)

type CustomerSearchDTO struct {
	_dto.Page
	Phone      string                  `form:"phone"`
	FullName   string                  `form:"fullName"`
	Steps      []enums.EStep           `form:"steps"`
	Sources    []enums.ESourceCustomer `form:"customerSources"`
	ChargeId   uint64                  `form:"chargeId"`
	ChargeType enums.EOwnerType        `form:"chargeType"`
	FromDate   *time.Time              `form:"fromDate"`
	ToDate     *time.Time              `form:"toDate"`
	// OwnerID    uint64                  `form:"ownerId"`
	// OwnerType  enums.EOwnerType        `form:"ownerType"`

}
