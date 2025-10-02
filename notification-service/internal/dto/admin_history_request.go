package dto

import (
	_dto "common/domain/dto"
	shared_enum "pb/enums"
	"time"
)

// AdminHistorySearchDTO đại diện cho request tìm kiếm lịch sử admin
type AdminHistorySearchDTO struct {
	_dto.Pagable
	AdminID    *uint64                     `form:"adminId"`
	OwnerID    *uint64                     `form:"ownerId"`
	OwnerType  *shared_enum.EOwnerType     `form:"ownerType"`
	TargetID   *uint64                     `form:"targetId"`
	TargetType *shared_enum.ETargetHistory `form:"targetType"`
	ActionType *shared_enum.EHistory       `form:"actionType"`
	AdminRole  *string                     `form:"adminRole"`
	FromDate   *time.Time                  `form:"fromDate"`
	ToDate     *time.Time                  `form:"toDate"`
	IsInternal *bool                       `form:"isInternal"`
	IPAddress  *string                     `form:"ipAddress"`
	Status     *string                     `form:"status"`
}

// AdminHistoryCreateDTO đại diện cho request tạo lịch sử admin
type AdminHistoryCreateDTO struct {
	TargetID   uint64                     `json:"targetId" binding:"required"`
	TargetType shared_enum.ETargetHistory `json:"targetType" binding:"required"`
	ActionType shared_enum.EHistory       `json:"actionType" binding:"required"`
	Title      string                     `json:"title,omitempty"`
	Note       []string                   `json:"note,omitempty"`
	PreStage   string                     `json:"preStage,omitempty"`
	AfterStage string                     `json:"afterStage,omitempty"`
	AdminID    uint64                     `json:"adminId" binding:"required"`
	OwnerID    *uint64                    `json:"ownerId"`
	OwnerType  shared_enum.EOwnerType     `json:"ownerType"`
	IsInternal bool                       `json:"isInternal"`
	AdminRole  string                     `json:"adminRole,omitempty"`
	IPAddress  string                     `json:"ipAddress,omitempty"`
	UserAgent  string                     `json:"userAgent,omitempty"`
}

// AdminActionStatsDTO đại diện cho thống kê hành động admin
type AdminActionStatsDTO struct {
	AdminID      uint64           `json:"adminId"`
	AdminRole    string           `json:"adminRole"`
	TotalActions int64            `json:"totalActions"`
	ActionCounts map[string]int64 `json:"actionCounts"`
	FromDate     time.Time        `json:"fromDate"`
	ToDate       time.Time        `json:"toDate"`
}
