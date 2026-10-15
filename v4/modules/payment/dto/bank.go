package dto

// BankListRequest is GET /v2/payment/bank/all (client sync timestamp).
type BankListRequest struct {
	Timestamp int64 `form:"timestamp"`
}
