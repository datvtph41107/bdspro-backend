package auth

import "time"

// AuthStatus is registration/credential lifecycle state keyed by auth_id.
// After a profile exists, account lock/status truth lives on user_profile.
type AuthStatus struct {
	IdDomain
	Active      bool      `gorm:"column:active" json:"active"`
	Verified    bool      `gorm:"column:verified" json:"verified"`
	LockedUntil time.Time `gorm:"column:locked_until" json:"lockedUntil"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (AuthStatus) TableName() string {
	return "auth_status"
}
