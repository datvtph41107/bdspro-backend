package dto

import (
	"time"

	"mdn4/modules/bdspro/enums"
)

// DealSetting is PUT deal setting payload (parity with bdspro-service dto.DealSetting).
type DealSetting struct {
	FromDate                *time.Time
	ToDate                  *time.Time
	AllowSharing            bool
	MemberCanAddTransaction bool
	OnlyOwnerGetCommission  bool
	InternalNote            string
	OwnerID                 uint64
	Status                  enums.DealStatus
	Name                    string
	AccountNumber           string
	AccountName             string
	BankID                  uint64
	AllowManualInput        bool
	TargetProfit            float64
	ShareVisibility         enums.EVisibility
	Description             string
}

// DealSearchRequest mirrors bdspropb.DealSearchRequest query params.
type DealSearchRequest struct {
	Page             uint32  `form:"page"`
	Size             uint32  `form:"size"`
	Text             string  `form:"text"`
	OrganizationID   *uint64 `form:"organizationId"`
	GroupID          *uint64 `form:"groupId"`
	OwnerID          *uint64 `form:"ownerId"`
	OwnerOf          uint32  `form:"ownerOf"`
	InvitedProfileID *uint64 `form:"invitedProfileId"`
	Timestamp        int64   `form:"timestamp"`
}

// DealDetailRequest is GET /v4/deal/{id} query (sync timestamp).
type DealDetailRequest struct {
	ID        uint64 `path:"id"`
	Timestamp int64  `form:"timestamp"`
}

// DealListByIDRequest is shared page/size/timestamp for deal-scoped lists (products, invitations, histories, …).
type DealListByIDRequest struct {
	DealID    uint64  `path:"id|dealId"`
	Page      uint32  `form:"page"`
	Size      uint32  `form:"size"`
	Timestamp int64   `form:"timestamp"`
	MemberID  *uint64 `form:"memberId"`
}

// AddDealInternalNoteRequest is POST /v4/deal/{id}/internal-note.
type AddDealInternalNoteRequest struct {
	DealID     uint64 `path:"id|dealId"`
	Content    string `json:"content"`
	ActionType uint32 `json:"actionType"`
}

// DealMemberCommissionUpdate is one member commission row (batch update).
type DealMemberCommissionUpdate struct {
	MemberID        uint64  `json:"memberId"`
	CommissionValue float64 `json:"commissionValue"`
	CommissionType  uint32  `json:"commissionType"`
	Note            string  `json:"note"`
}

// UpdateDealMembersCommissionRequest is POST /v4/deal/{id}/commission.
type UpdateDealMembersCommissionRequest struct {
	DealID  uint64                       `path:"id|dealId"`
	Members []DealMemberCommissionUpdate `json:"members"`
}

// SaveDealProductsRequest is POST /v4/deal/{id}/products (add + remove by creator).
type SaveDealProductsRequest struct {
	DealID           uint64   `path:"id|dealId"`
	ProductIDs       []uint64 `json:"productIds"`
	RemoveProductIDs []uint64 `json:"removeProductIds"`
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

// DealSummaryRequest is GET /v4/deal/{id}/summary. With no query flags, all summary fields are returned.
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

// DealStatisticsRequest is GET .../contract/deal/{id}/statistic(s).
type DealStatisticsRequest struct {
	DealID    uint64 `path:"id|dealId"`
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

// DealProductItem is one deal_products row plus nested product detail
// (same shape as GET product detail) for GET /v4/deal/{id}/products.
type DealProductItem struct {
	ID        uint64  `json:"id"`
	CreatedAt string  `json:"createdAt,omitempty"`
	UpdatedAt string  `json:"updatedAt,omitempty"`
	CreatedBy *uint64 `json:"createdBy,omitempty"`
	UpdatedBy *uint64 `json:"updatedBy,omitempty"`
	DealID    uint64  `json:"dealId"`
	ProductID uint64  `json:"productId"`
	// Product is set by the service as *bdspropb.ProductResponse; controller
	// converts it to the HTTP envelope object.
	Product any `json:"product,omitempty"`
}

// ProductDealItem is one deal_products row plus nested deal detail
// (same shape as GET /v4/deal/{id}) for GET /v4/bdspro/product/{id}/deals.
type ProductDealItem struct {
	ID        uint64  `json:"id"`
	CreatedAt string  `json:"createdAt,omitempty"`
	UpdatedAt string  `json:"updatedAt,omitempty"`
	CreatedBy *uint64 `json:"createdBy,omitempty"`
	UpdatedBy *uint64 `json:"updatedBy,omitempty"`
	DealID    uint64  `json:"dealId"`
	ProductID uint64  `json:"productId"`
	// Deal is set by the service as *bdspropb.GroupDeal; controller
	// converts it to the HTTP envelope object.
	Deal any `json:"deal,omitempty"`
}
