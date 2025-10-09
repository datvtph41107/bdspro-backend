package providers

import (
	"context"
	"time"
)

// UserStatusInfo chứa thông tin trạng thái của user
type UserStatusInfo struct {
	ProfileID   uint64  `json:"profileId"`
	Verified    bool    `json:"verified"`
	VerifiedAt  *string `json:"verifiedAt,omitempty"`
	LockedUntil *string `json:"lockedUntil,omitempty"`
	Active      bool    `json:"active"`
}

// AuthMethodDomain chứa thông tin auth method
type AuthMethodDomain struct {
	ID               uint64     `json:"id"`
	UserID           uint64     `json:"userId"`
	Provider         string     `json:"provider"`
	Username         string     `json:"username"`
	Password         string     `json:"password"`
	Name             string     `json:"name"`
	Email            string     `json:"email"`
	Phone            string     `json:"phone"`
	Avatar           string     `json:"avatar"`
	JobTitle         string     `json:"jobTitle"`
	RoleKey          uint32     `json:"roleKey"`
	Attachments      string     `json:"attachments"`
	InternalNotes    string     `json:"internalNotes"`
	SendNotification bool       `json:"sendNotification"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	DeletedAt        *time.Time `json:"deletedAt,omitempty"`
}

// ProfileDTO chứa thông tin profile
type ProfileDTO struct {
	ProfileID uint64 `json:"profileId"`
	FullName  string `json:"fullName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Avatar    string `json:"avatar"`
	RoleType  string `json:"roleType"`
}

// AdminItem chứa thông tin admin item
type AdminItem struct {
	AuthID    uint64    `json:"authId"`
	ProfileID uint64    `json:"profileId"`
	Username  string    `json:"username"`
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Avatar    string    `json:"avatar"`
	Role      string    `json:"role"`
	RoleType  string    `json:"roleType"`
	CreatedAt time.Time `json:"createdAt"`
}

// AdminListRequest DTO cho request lấy danh sách admin
type AdminListRequest struct {
	Page uint32 `json:"page"`
	Size uint32 `json:"size"`
}

// AdminListResponse DTO cho response danh sách admin
type AdminListResponse struct {
	Data  []AdminItem `json:"data"`
	Total int32       `json:"total"`
}

// AuthProvider định nghĩa interface cho việc tương tác với auth-service
type AuthProvider interface {
	// GetUserStatus lấy thông tin trạng thái của 1 user
	GetUserStatus(ctx context.Context, profileID uint64) (*UserStatusInfo, error)

	// GetUserStatusBatch lấy thông tin trạng thái của nhiều user cùng lúc
	GetUserStatusBatch(ctx context.Context, profileIDs []uint64) (map[uint64]*UserStatusInfo, error)

	// AuthMethod Management methods
	FindByUsernameAndProvider(ctx context.Context, username, provider string) (*AuthMethodDomain, error)
	FindByEmail(ctx context.Context, email string) (*AuthMethodDomain, error)
	FindByID(ctx context.Context, id uint64) (*AuthMethodDomain, error)
	Create(ctx context.Context, authMethod *AuthMethodDomain) (*AuthMethodDomain, error)
	Update(ctx context.Context, authMethod *AuthMethodDomain) (*AuthMethodDomain, error)
	SoftDelete(ctx context.Context, id uint64) error
	FindAdminsByProvider(ctx context.Context, provider string, page, size int) ([]*AuthMethodDomain, int64, error)
}
