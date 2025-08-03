package dto

import (
	"time"
	"organization/internal/enums"
)

type DealSetting struct {
	FromDate                time.Time
	ToDate                  time.Time
	AllowSharing            bool
	MemberCanAddTransaction bool
	OnlyOwnerGetCommission  bool
	InternalNote            string
	OwnerId                 uint64
	Status                  enums.DealStatus
}