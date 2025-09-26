package providers

import "context"

// AuthProvider interface cho việc tương tác với auth-service
type AuthProvider interface {
	// GetUserStatus lấy thông tin status của user từ auth-service
	GetUserStatus(ctx context.Context, profileID uint64) (*UserStatusInfo, error)

	// GetUserStatusBatch lấy thông tin status của nhiều user cùng lúc
	GetUserStatusBatch(ctx context.Context, profileIDs []uint64) (map[uint64]*UserStatusInfo, error)
}

// UserStatusInfo thông tin status của user từ auth-service
type UserStatusInfo struct {
	ProfileID   uint64  `json:"profileId"`
	Verified    bool    `json:"verified"`
	VerifiedAt  *string `json:"verifiedAt,omitempty"`
	LockedUntil *string `json:"lockedUntil,omitempty"`
	Active      bool    `json:"active"`
}
