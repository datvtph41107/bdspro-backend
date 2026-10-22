package dto

import _dto "common/domain/dto"

// ListingListRequest lists listings for the current user (parity PostSearchRequest subset).
type ListingListRequest struct {
	_dto.Pagable
	Text            string `form:"text"`
	TransactionType int32  `form:"transactionType"`
	ProductID       uint64 `form:"productId"`
	Timestamp       int64  `form:"timestamp"`
}
