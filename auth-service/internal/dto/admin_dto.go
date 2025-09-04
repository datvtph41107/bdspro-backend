package dto

import (
	"time"

	"auth/internal/enums"
)

// CreateAdminRequest DTO cho tạo admin mới
type CreateAdminRequest struct {
	Username string `json:"username" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required" validate:"required"`
	FullName string `json:"fullName" binding:"required" validate:"required"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Phone    string `json:"phone" binding:"required" validate:"required"`
	Avatar   string `json:"avatar"`
}

// UpdateAdminRequest DTO cho cập nhật admin
type UpdateAdminRequest struct {
	AuthID   uint64 `json:"authId" binding:"required" validate:"required"`
	FullName string `json:"fullName" binding:"required" validate:"required"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Phone    string `json:"phone" binding:"required" validate:"required"`
	Avatar   string `json:"avatar"`
}

// DeleteAdminRequest DTO cho xóa admin
type DeleteAdminRequest struct {
	AuthID uint64 `json:"authId" binding:"required" validate:"required"`
}

// AssignRoleRequest DTO cho gán quyền admin
type AssignRoleRequest struct {
	AuthID  uint64 `json:"authId" binding:"required" validate:"required"`
	RoleKey string `json:"roleKey" binding:"required" validate:"required"`
}

// AdminListRequest DTO cho danh sách admin
type AdminListRequest struct {
	Page int32 `json:"page" form:"page" binding:"min=1" validate:"min=1"`
	Size int32 `json:"size" form:"size" binding:"min=1,max=100" validate:"min=1,max=100"`
}

// AdminSearchRequest DTO cho tìm kiếm admin
type AdminSearchRequest struct {
	Page       int32  `json:"page" form:"page" binding:"min=1" validate:"min=1"`
	Size       int32  `json:"size" form:"size" binding:"min=1,max=100" validate:"min=1,max=100"`
	SearchName string `json:"searchName" form:"searchName"`
	Role       string `json:"role" form:"role"`
	Department string `json:"department" form:"department"`
	Status     string `json:"status" form:"status"`
}

// AdminListResponse DTO cho response danh sách admin
type AdminListResponse struct {
	Data  []AdminItem `json:"data"`
	Total int32       `json:"total"`
}

// AdminItem DTO cho item admin trong danh sách
type AdminItem struct {
	AuthID    uint64      `json:"authId"`
	ProfileID uint64      `json:"profileId"`
	Username  string      `json:"username"`
	FullName  string      `json:"fullName"`
	Email     string      `json:"email"`
	Phone     string      `json:"phone"`
	Avatar    string      `json:"avatar"`
	Role      string      `json:"role"`
	RoleType  enums.ERole `json:"roleType"`
	CreatedAt time.Time   `json:"createdAt"`
}

// AdminResponse DTO cho response admin
type AdminResponse struct {
	AuthID    uint64 `json:"authId"`
	ProfileID uint64 `json:"profileId"`
	Username  string `json:"username"`
	FullName  string `json:"fullName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Avatar    string `json:"avatar"`
	Role      string `json:"role"`
}

// AdminRoleDTO DTO cho role admin từ organization service
type AdminRoleDTO struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}
