package dto

import (
	"organization/internal/enums"
	"time"
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
	Name                    string
	BankName                string
	AccountNumber           string
	AccountName             string
	BankId                  uint64
	AllowManualInput        bool
}

type SummaryRequest struct {
	DealId          uint64
	AmountPending   bool
	AmountApproved  bool
	AmountRejected  bool
	NumPending      bool
	NumApproved     bool
	NumRejected     bool
	NumInvestment   bool
	NumMemberSubmit bool
	AmountTarget    bool
	PercentDone     bool
	AmountCost      bool
	AmountProfit    bool
	AmountTotal     bool
	NumTransaction  bool
	All             bool
}

type SummaryResponse struct {
	AmountPending   float64 `json:"amount_pending"`
	AmountApproved  float64 `json:"amount_approved"`
	AmountRejected  float64 `json:"amount_rejected"`
	AmountRest      float64 `json:"amount_rest"`
	NumPending      uint32  `json:"num_pending"`
	NumApproved     uint32  `json:"num_approved"`
	NumRejected     uint32  `json:"num_rejected"`
	NumInvestment   uint32  `json:"num_investment"`
	NumMemberSubmit uint32  `json:"num_member_submit"`
	AmountTarget    float64 `json:"amount_target"`
	PercentDone     float64 `json:"percent_done"`
	AmountCost      float64 `json:"amount_cost"`
	AmountProfit    float64 `json:"amount_profit"`
	AmountTotal     float64 `json:"amount_total"`
	NumTransaction  uint32  `json:"num_transaction"`
}
