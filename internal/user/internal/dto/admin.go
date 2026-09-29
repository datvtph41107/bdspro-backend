package dto

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	"time"
)

// ListUsersRequest DTO cho request lấy danh sách tất cả user
type ListUsersRequest struct {
	_dto.Pagable
	Search    string     `json:"search" form:"search"`
	SortBy    string     `json:"sortBy" form:"sortBy"`
	SortOrder string     `json:"sortOrder" form:"sortOrder"`
	Status    *uint32    `json:"status" form:"status"`
	RoleType  string     `json:"roleType" form:"roleType"`
	RoleID    *uint64    `json:"roleId" form:"roleId"`
	StartDate *time.Time `json:"startDate" form:"startDate" time_format:"2006-01-02T15:04:05"`
	EndDate   *time.Time `json:"endDate" form:"endDate" time_format:"2006-01-02T15:04:05"`
	Bookmark  *bool      `json:"bookmark" form:"bookmark"`
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
	ProfileID      uint64              `json:"profileId"`
	FullName       string              `json:"fullName"`
	Email          string              `json:"email"`
	Phone          string              `json:"phone"`
	Avatar         string              `json:"avatar"`
	Address        string              `json:"address"`
	Gender         uint8               `json:"gender"`
	Birth          *time.Time          `json:"birth"`
	CreatedAt      time.Time           `json:"createdAt"`
	UpdatedAt      time.Time           `json:"updatedAt"`
	Status         _enum.EUserStatus   `json:"status"`
	WarningLevel   _enum.EWarningLevel `json:"warningLevel"`
	RoleType       string              `json:"roleType"`
	RoleRealEstate string              `json:"roleRealEstate"`
	Position       string              `json:"position"`
	Workplace      string              `json:"workplace"`
	DepartmentID   uint64              `json:"departmentId"`
	Bookmark       bool                `json:"isBookmark"`
	TickVerified   bool                `json:"tickVerified"`
	LastLoginAt    *time.Time          `json:"lastLoginAt"`
	TotalLogin     uint32              `json:"totalLogin"`
	VerifiedAt     *time.Time          `json:"verifiedAt"`
	LockedAt       *time.Time          `json:"lockedAt"`
	RoleID         *uint64             `json:"roleId"`
	RoleName       string              `json:"roleName"`
	RoleKey        string              `json:"roleKey"`
	RoleColor      string              `json:"roleColor"`
	RoleBgColor    string              `json:"roleBgColor"`
}

// Admin API Request DTOs
type LockUserRequest struct {
	ProfileID uint64            `json:"profileId"`
	Reason    string            `json:"reason"`
	LockType  _enum.EUserStatus `json:"lockType"`
}

type UnLockUserRequest struct {
	ProfileID uint64 `json:"profileId"`
	Reason    string `json:"reason"`
}

type ApproveUserRequest struct {
	ProfileID uint64 `json:"profileId"`
	Note      string `json:"note"`
}

// CreateUserRequest DTO cho request tạo user mới
type CreateUserRequest struct {
	Username string     `json:"username"` // Optional - sẽ tự động generate từ email nếu không có
	Password string     `json:"password"` // Optional - sẽ tự động generate nếu không có
	FullName string     `json:"fullName" validate:"required"`
	Email    string     `json:"email" validate:"required,email"`
	Phone    string     `json:"phone" validate:"required"`
	Avatar   string     `json:"avatar"`
	RoleID   uint64     `json:"roleId"`
	Gender   *uint32    `json:"gender"`
	Address  string     `json:"address"`
	Birth    *time.Time `json:"birth"`
}

// UpdateUserRequest DTO cho request cập nhật user
type UpdateUserRequest struct {
	ProfileID uint64            `json:"profileId" validate:"required"`
	FullName  string            `json:"fullName"`
	Email     string            `json:"email" validate:"omitempty,email"`
	Phone     string            `json:"phone"`
	Avatar    string            `json:"avatar"`
	RoleID    *uint64           `json:"roleId"`
	Gender    *uint32           `json:"gender"`
	Address   string            `json:"address"`
	Birth     *time.Time        `json:"birth"`
	Status    _enum.EUserStatus `json:"status"`
}

// DeleteUserRequest DTO cho request xóa user
type DeleteUserRequest struct {
	ProfileID uint64 `json:"profileId" validate:"required"`
	Reason    string `json:"reason"`
}

// CreateUserResponse DTO cho response tạo user
type CreateUserResponse struct {
	ProfileID uint64    `json:"profileId"`
	AuthID    uint64    `json:"authId"`
	Username  string    `json:"username"`
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Avatar    string    `json:"avatar"`
	RoleID    uint64    `json:"roleId"`
	CreatedAt time.Time `json:"createdAt"`
}

// UpdateUserResponse DTO cho response cập nhật user
type UpdateUserResponse struct {
	ProfileID uint64            `json:"profileId"`
	FullName  string            `json:"fullName"`
	Email     string            `json:"email"`
	Phone     string            `json:"phone"`
	Avatar    string            `json:"avatar"`
	RoleID    uint64            `json:"roleId"`
	Status    _enum.EUserStatus `json:"status"`
	UpdatedAt *time.Time        `json:"updatedAt"`
}

// GetUserDetailResponse DTO cho response chi tiết user
type GetUserDetailResponse struct {
	ProfileID       uint64            `json:"profileId"`
	FullName        string            `json:"fullName"`
	Email           string            `json:"email"`
	Phone           string            `json:"phone"`
	Phone2          string            `json:"phone2"`
	Avatar          string            `json:"avatar"`
	Address         string            `json:"address"`
	Gender          uint32            `json:"gender"`
	Birth           *time.Time        `json:"birth"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
	Status          _enum.EUserStatus `json:"status"`
	Position        string            `json:"position"`
	Workplace       string            `json:"workplace"`
	DepartmentID    uint64            `json:"departmentId"`
	TickVerified    bool              `json:"tickVerified"`
	BackgroundImage string            `json:"backgroundImage"`
	TaxCode         string            `json:"taxCode"`
	FrontIdentify   string            `json:"frontIdentify"`
	BackIdentify    string            `json:"backIdentify"`
	Slogan          string            `json:"slogan"`
	Website         string            `json:"website"`
	Facebook        string            `json:"facebook"`
	Instagram       string            `json:"instagram"`
	Twitter         string            `json:"twitter"`
	Linkedin        string            `json:"linkedin"`
	Youtube         string            `json:"youtube"`
	Introduction    string            `json:"introduction"`
	RoleID          *uint64           `json:"roleId"`
	RoleName        string            `json:"roleName"`
	RoleKey         uint32            `json:"roleKey"`
	RoleColor       string            `json:"roleColor"`
	RoleBgColor     string            `json:"roleBgColor"`
}
