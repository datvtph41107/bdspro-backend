package dto

import "time"

// DealSearchRequest mirrors bdspropb.DealSearchRequest query params.
type DealSearchRequest struct {
	Page           uint32  `form:"page"`
	Size           uint32  `form:"size"`
	Text           string  `form:"text"`
	OrganizationID *uint64 `form:"organizationId"`
	GroupID        *uint64 `form:"groupId"`
	OwnerID        *uint64 `form:"ownerId"`
	OwnerOf        uint32  `form:"ownerOf"`
	Timestamp      int64   `form:"timestamp"`
}

// DealDetailRequest is GET /v4/deal/{id} query (sync timestamp).
type DealDetailRequest struct {
	ID        uint64 `path:"id"`
	Timestamp int64  `form:"timestamp"`
}

// DealListByIDRequest is shared page/size/timestamp for deal-scoped lists.
type DealListByIDRequest struct {
	DealID    uint64  `path:"id|dealId"`
	Page      uint32  `form:"page"`
	Size      uint32  `form:"size"`
	Timestamp int64   `form:"timestamp"`
	MemberID  *uint64 `form:"memberId"`
}

// DealMembersRequest is GET /v4/deal/{id}/members (was /v2.1/org/deal/{id}/members).
type DealMembersRequest struct {
	DealID         uint64 `path:"id"`
	Keyword        string `form:"keyword"`
	Page           uint32 `form:"page"`
	Size           uint32 `form:"size"`
	DoneInvestment *bool  `form:"doneInvestment"`
	Timestamp      int64  `form:"timestamp"`
}

// DealSummaryRequest is GET /v4/deal/{id}/summary query flags (was /v2.1/org/deal/{dealId}/summary).
type DealSummaryRequest struct {
	DealID          uint64 `path:"id"`
	AmountPending   bool   `form:"amountPending"`
	AmountApproved  bool   `form:"amountApproved"`
	AmountRejected  bool   `form:"amountRejected"`
	NumPending      bool   `form:"numPending"`
	NumApproved     bool   `form:"numApproved"`
	NumRejected     bool   `form:"numRejected"`
	NumInvestment   bool   `form:"numInvestment"`
	NumMemberSubmit bool   `form:"numMemberSubmit"`
	AmountTarget    bool   `form:"amountTarget"`
	PercentDone     bool   `form:"percentDone"`
	AmountTotal     bool   `form:"amountTotal"`
	AmountCost      bool   `form:"amountCost"`
	AmountProfit    bool   `form:"amountProfit"`
	NumTransaction  bool   `form:"numTransaction"`
	All             bool   `form:"all"`
	Timestamp       int64  `form:"timestamp"`
}

// DealSummaryRow is the aggregated investment/cost dashboard scan result.
type DealSummaryRow struct {
	AmountPending   float64 `gorm:"column:amount_pending"`
	AmountApproved  float64 `gorm:"column:amount_approved"`
	AmountRejected  float64 `gorm:"column:amount_rejected"`
	NumPending      uint32  `gorm:"column:num_pending"`
	NumApproved     uint32  `gorm:"column:num_approved"`
	NumRejected     uint32  `gorm:"column:num_rejected"`
	NumInvestment   uint32  `gorm:"column:num_investment"`
	NumMemberSubmit uint32  `gorm:"column:num_member_submit"`
	AmountCost      float64 `gorm:"column:amount_cost"`
}

// DealStatisticsRequest is GET /v4/transaction/contract/deal/{id}/statistic.
type DealStatisticsRequest struct {
	DealID    uint64 `path:"id"`
	Payment   bool   `form:"payment"`
	Cost      bool   `form:"cost"`
	Timestamp int64  `form:"timestamp"`
}

// TxStats is approved/reject/pending/total buckets used by deal statistics.
type TxStats struct {
	Pending  float64
	Approved float64
	Rejected float64
	Total    float64
}

// DealCostStatsRow is cost/revenue aggregates from tx_deal_costs.
type DealCostStatsRow struct {
	AmountCost    TxStats
	AmountRevenue TxStats
	NumCost       TxStats
	NumRevenue    TxStats
}

// DealPaymentStatsRow is payment aggregates from sale/rent contracts.
type DealPaymentStatsRow struct {
	Payment    TxStats
	NumPayment TxStats
}

// TxProcessRow is one first-of-each-action row for deal contract process.
type TxProcessRow struct {
	Action        uint32     `gorm:"column:action"`
	Timestamp     *time.Time `gorm:"column:created_at"`
	TransactionID uint64     `gorm:"column:tx_id"`
	Note          string     `gorm:"column:note"`
	Message       string     `gorm:"column:message"`
}

// TxContractDealListRow is a joined contract list row for a deal.
type TxContractDealListRow struct {
	ContractID      uint64     `gorm:"column:contract_id"`
	DealID          uint64     `gorm:"column:deal_id"`
	TransactionName string     `gorm:"column:transaction_name"`
	Method          uint32     `gorm:"column:method"`
	CustomerID      uint64     `gorm:"column:customer_id"`
	Status          uint32     `gorm:"column:status"`
	Amount          float64    `gorm:"column:amount"`
	Timestamp       *time.Time `gorm:"column:timestamp"`
}
