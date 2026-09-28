package dto

// PropertyListRequest is the query for GET /v2/bdspro/v2/property/me and detail.
type PropertyListRequest struct {
	ID        uint64 `form:"id" uri:"id"`
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	Timestamp int64  `form:"timestamp"`
}
