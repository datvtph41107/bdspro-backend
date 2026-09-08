package access

import (
	"time"
)

// AuthAdminWithRole chứa thông tin admin kèm với role
type AuthAdminWithRole struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"userId"`
	Username  string    `json:"username"`
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Avatar    string    `json:"avatar"`
	RoleKey   uint32    `json:"roleKey"`
	Role      *Role     `json:"role,omitempty"`
	Status    uint32    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
