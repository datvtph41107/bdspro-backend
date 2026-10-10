package dto

// InvestmentListQuery mirrors bdspropb.InvestmentListRequest.
type InvestmentListQuery struct {
	DealID     uint64  `path:"dealId"`
	MemberID   *uint64 `form:"memberId"`
	InvestType uint32  `form:"investType"`
	Status     uint32  `form:"status"`
	Timestamp  int64   `form:"timestamp"`
}
