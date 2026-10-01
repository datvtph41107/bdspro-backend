package dto

// GetEnumsRequest mirrors sharepb.SyncRequest query for GET /v2/tqd/public/enums.
type GetEnumsRequest struct {
	Timestamp int64 `form:"timestamp"`
}
