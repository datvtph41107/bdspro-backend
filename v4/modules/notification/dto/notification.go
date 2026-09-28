package dto

// CountRequest mirrors notificationpb.CountRequest query params.
type CountRequest struct {
	Timestamp int64 `form:"timestamp"`
}

// CountResponse mirrors notificationpb.CountResponse for GET /v2/notification/count.
type CountResponse struct {
	Count     int32 `json:"count"`
	Unread    int32 `json:"unread"`
	Timestamp int64 `json:"timestamp"`
	Unchanged bool  `json:"unchanged"`
}

// ListRequest mirrors notificationpb.ListRequest query params.
type ListRequest struct {
	Page      int32  `form:"page"`
	Size      int32  `form:"size"`
	Sort      string `form:"sort"`
	IsRead    bool   `form:"isRead"`
	Timestamp int64  `form:"timestamp"`
}

// ReadManyRequest mirrors notificationpb.ReadManyRequest body.
type ReadManyRequest struct {
	IDs []uint64 `json:"ids"`
}
