package dto

// ListProvincesRequest is the query shape for GET /v2/tqd/client/provinces.
type ListProvincesRequest struct {
	Scope     int32 `form:"scope"`
	Timestamp int64 `form:"timestamp"`
}

// ListWardsRequest is the query shape for GET /v2/tqd/client/wards.
type ListWardsRequest struct {
	Page      *int32 `form:"page"`
	PageSize  *int32 `form:"pageSize"`
	Timestamp int64  `form:"timestamp"`
}
