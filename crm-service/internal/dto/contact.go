package dto

import (
	_dto "common/dto"
	"crm/internal/enums"
	"time"
)

type ContactSearchDTO struct {
	_dto.Page
	Phone              string                  `form:"phone"`
	FullName           string                  `form:"fullName"`
	Steps              []enums.EStep           `form:"steps"`
	Sources            []enums.ESourceCustomer `form:"customerSources"`
	ChargeId           uint64                  `form:"chargeId"`
	ChargeType         enums.EOwnerType        `form:"chargeType"`
	CategoryProductIDs []uint64                `form:"categoryProductIds"`
	FromDate           *time.Time              `form:"fromDate"`
	ToDate             *time.Time              `form:"toDate"`
}
