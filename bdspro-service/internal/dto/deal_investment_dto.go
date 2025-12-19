package dto

import (
	"bdspro/internal/enums"
	"time"
)

type InvestmentDTO struct {
	ID           uint64               `json:"id"`
	DealID       uint64               `json:"deal_id"`
	MemberID     uint64               `json:"member_id"`
	InvestType   enums.DealMemberType `json:"invest_type"`
	Amount       float64              `json:"amount"`
	TransferTime time.Time            `json:"transfer_time"`
	Note         string               `json:"note"`
	ProxyUserID  uint64               `json:"proxy_user_id"`
}

type InfoInvestmentDTO struct {
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
}
