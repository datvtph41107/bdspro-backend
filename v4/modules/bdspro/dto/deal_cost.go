package dto

import _dto "common/domain/dto"

// DealCostSearchRequest mirrors bdspropb.DealCostSearchDTO.
type DealCostSearchRequest struct {
	_dto.Pagable
	DealID    uint64 `path:"dealId"`
	Timestamp int64  `form:"timestamp"`
}
