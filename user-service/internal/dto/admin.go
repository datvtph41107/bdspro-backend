package dto

import (
	_dto "common/dto"
	"time"
)

// ListAllUsersRequest DTO cho request lấy danh sách tất cả user
type ListAllUsersRequest struct {
	_dto.Pagable
	Search    string     `json:"search" form:"search"`
	SortBy    string     `json:"sortBy" form:"sortBy"`
	SortOrder string     `json:"sortOrder" form:"sortOrder"`
	Status    *uint32    `json:"status" form:"status"`
	RoleType  string     `json:"roleType" form:"roleType"`
	StartDate *time.Time `json:"startDate" form:"startDate" time_format:"2006-01-02T15:04:05"`
	EndDate   *time.Time `json:"endDate" form:"endDate" time_format:"2006-01-02T15:04:05"`
}

// ListAllUsersResponse DTO cho response danh sách tất cả user
type ListAllUsersResponse struct {
	Data       []AdminUserItem `json:"data"`
	Total      int32           `json:"total"`
	Page       uint32          `json:"page"`
	Size       uint32          `json:"size"`
	TotalPages uint32          `json:"totalPages"`
}

// AdminUserItem DTO cho thông tin user trong admin view
type AdminUserItem struct {
	ProfileID      uint64     `json:"profileId"`
	FullName       string     `json:"fullName"`
	Email          string     `json:"email"`
	Phone          string     `json:"phone"`
	Avatar         string     `json:"avatar"`
	Address        string     `json:"address"`
	Gender         uint8      `json:"gender"`
	Birth          *time.Time `json:"birth"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	Status         uint8      `json:"status"`
	RoleType       string     `json:"roleType"`
	RoleRealEstate string     `json:"roleRealEstate"`
	Position       string     `json:"position"`
	DepartmentID   uint64     `json:"departmentId"`
	TickVerified   bool       `json:"tickVerified"`
	LastLoginAt    *time.Time `json:"lastLoginAt"`
	TotalLogin     uint32     `json:"totalLogin"`
}

// Admin API Request DTOs
type BlockUserRequest struct {
	ProfileID uint64 `json:"profileId"`
	Reason    string `json:"reason"`
}

type UnblockUserRequest struct {
	ProfileID uint64 `json:"profileId"`
	Reason    string `json:"reason"`
}

type ApproveUserRequest struct {
	ProfileID uint64 `json:"profileId"`
	Note      string `json:"note"`
}
