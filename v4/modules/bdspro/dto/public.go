package dto

// GetEnumsRequest mirrors sharepb.SyncRequest query for GET /v2/bdspro/public/enums.
type GetEnumsRequest struct {
	Timestamp int64 `form:"timestamp"`
}

// GetEnumByNameRequest is GET /v2/bdspro/public/enums/{text}.
type GetEnumByNameRequest struct {
	Text string `path:"text"`
}
