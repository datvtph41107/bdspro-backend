package dto

import _dto "common/domain/dto"

// CostTypeSearchRequest mirrors bdspropb.CostTypeSearchDTO query params.
type CostTypeSearchRequest struct {
	_dto.Pagable
	TypeName  string `form:"typeName"`
	CostType  uint32 `form:"costType"`
	Timestamp int64  `form:"timestamp"`
}

// CostTypePathRequest is PUT/DELETE /v2.1/transaction/cost-type/{id}.
type CostTypePathRequest struct {
	ID uint64 `path:"id"`
}
