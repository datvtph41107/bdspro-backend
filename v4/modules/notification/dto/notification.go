package dto

import shared_enum "pb/enums"

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

// AdminHistorySearchRequest mirrors notipb.AdminHistorySearchRequest query params.
type AdminHistorySearchRequest struct {
	Page       int32   `form:"page"`
	Size       int32   `form:"size"`
	AdminID    *uint64 `form:"adminId"`
	OwnerID    *uint64 `form:"ownerId"`
	OwnerType  *int32  `form:"ownerType"`
	TargetID   *uint64 `form:"targetId"`
	TargetType *int32  `form:"targetType"`
	ActionType *int32  `form:"actionType"`
	AdminRole  *string `form:"adminRole"`
	FromDate   *string `form:"fromDate"`
	ToDate     *string `form:"toDate"`
	IsInternal *bool   `form:"isInternal"`
}

// AdminHistoryFilter is the repository filter for admin history search.
type AdminHistoryFilter struct {
	Page       int
	Size       int
	AdminID    *uint64
	OwnerID    *uint64
	OwnerType  *shared_enum.EOwnerType
	TargetID   *uint64
	TargetType *shared_enum.ETargetHistory
	ActionType *shared_enum.EHistory
	AdminRole  *string
	FromDate   *string // RFC3339
	ToDate     *string // RFC3339
	IsInternal *bool
}
