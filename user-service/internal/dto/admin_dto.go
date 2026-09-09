package dto

import (
	_dto "common/domain/dto"
	"fmt"
	"strings"
	"time"

	"user/internal/enums"
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

// AdminListRequest DTO cho danh sách admin
type AdminListRequest struct {
	_dto.Pagable
	// Page   int32  `json:"page" form:"page" binding:"min=1" validate:"min=1"`
	// Size   int32  `json:"size" form:"size" binding:"min=1,max=100" validate:"min=1,max=100"`
	RoleID uint64 `json:"roleId" form:"roleId"`
	Name   string `json:"name" form:"name"`
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

// ChangePasswordRequest DTO cho thay đổi mật khẩu admin
type ChangePasswordRequest struct {
	AuthID      uint64 `json:"authId" binding:"required" validate:"required"`
	OldPassword string `json:"oldPassword" binding:"required" validate:"required"`
	NewPassword string `json:"newPassword" binding:"required" validate:"required"`
}

// ResetPasswordRequest DTO cho reset mật khẩu admin
type ResetPasswordRequest struct {
	AuthID      uint64 `json:"authId" binding:"required" validate:"required"`
	NewPassword string `json:"newPassword" binding:"required" validate:"required"`
}

// Validate thực hiện validation custom cho CreateAdminRequest
func (req *CreateAdminRequest) Validate() error {
	// Validation cho email nội bộ
	if req.Email != "" && !strings.Contains(req.Email, "@bdspro.vn") {
		// Có thể bỏ comment dòng này nếu muốn bắt buộc email nội bộ
		// return fmt.Errorf("email phải là email nội bộ (@bdspro.vn)")
	}

	// Validation cho phone nội bộ (có thể thêm logic kiểm tra số nội bộ)
	if req.Phone != "" && len(req.Phone) != 10 {
		return fmt.Errorf("số điện thoại phải có 10 chữ số")
	}

	// Validation cho password strength
	if req.Password != "" && len(req.Password) < 6 {
		return fmt.Errorf("mật khẩu phải có ít nhất 6 ký tự")
	}

	// Validation cho username
	if req.Username != "" && len(req.Username) < 3 {
		return fmt.Errorf("username phải có ít nhất 3 ký tự")
	}

	return nil
}
