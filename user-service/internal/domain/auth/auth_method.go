package auth

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Tạo enum type
type ProviderType string

const (
	ProviderGoogle   ProviderType = "GOOGLE"
	ProviderZalo     ProviderType = "ZALO"
	ProviderFacebook ProviderType = "FACEBOOK"
	ProviderPhone    ProviderType = "PHONE"
	ProviderAdmin    ProviderType = "ADMIN"
)

// Validate method
func (p ProviderType) IsValid() bool {
	switch p {
	case ProviderGoogle, ProviderZalo, ProviderFacebook, ProviderPhone, ProviderAdmin:
		return true
	}
	return false
}

// Trong struct
// GORM hook để validate
func (a *AuthMethod) BeforeCreate(tx *gorm.DB) error {
	if !a.Provider.IsValid() {
		return fmt.Errorf("invalid provider type: %s", a.Provider)
	}
	return nil
}

// AuthMethod tương đương với Java Entity
type AuthMethod struct {
	ID          uint64       `gorm:"primaryKey;autoIncrement"`
	Provider    ProviderType `gorm:"column:provider;type:varchar(20);not null"` // enum: email, zalo, facebook, phone, admin
	AuthName    string       `gorm:"column:auth_name;size:255;not null"`
	Password    string       `gorm:"column:password"`
	IsSensitive bool         `gorm:"column:is_sensitive"`
	Avatar      string       `gorm:"size:255"`
	FullName    string       `gorm:"size:255"`
	Email       string       `gorm:"size:255"`
	Phone       string       `gorm:"size:255"`
	RoleKey     uint32       `gorm:"column:role_key"`             // Role được gán cho auth method
	Status      uint8        `gorm:"column:status;default:1"`     // 1: active, 0: inactive, 2: temporarily_locked, 3: permanently_locked
	LockedAt    *time.Time   `gorm:"column:locked_at"`            // Thời gian bị khóa
	LockedUntil *time.Time   `gorm:"column:locked_until"`         // Thời gian khóa đến khi nào (cho khóa tạm thời)
	LockReason  string       `gorm:"column:lock_reason;size:500"` // Lý do khóa
	LockedBy    *uint64      `gorm:"column:locked_by"`            // ID admin khóa tài khoản
	PrivateKey  string       `gorm:"column:private_key;size:512"` // Diffie-Hellman private key
	PublicKey   string       `gorm:"column:public_key;size:512"`  // Diffie-Hellman public key
	AuthKey     string       `gorm:"column:auth_key;size:512"`    // Auth key được mã hóa
	CreatedAt   time.Time    `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time    `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   *time.Time   `gorm:"column:deleted_at"`
	UserID      uint64       `gorm:"column:user_id"` // Không cập nhật & xóa, chỉ tạo mới
}

// TableName đặt tên bảng tương đương với Java
func (AuthMethod) TableName() string {
	return "auth_method"
}

func (a *AuthMethod) Validate() error {
	if a.AuthName == "" || a.Provider == "" {
		return fmt.Errorf("auth_name or provider is required")
	}
	if !a.Provider.IsValid() {
		return fmt.Errorf("provider is invalid")
	}
	return nil
}

// IsLocked kiểm tra tài khoản có bị khóa không
func (a *AuthMethod) IsLocked() bool {
	return a.Status == uint8(StatusTemporarilyLocked) || a.Status == uint8(StatusPermanentlyLocked)
}

// IsTemporarilyLocked kiểm tra tài khoản có bị khóa tạm thời không
func (a *AuthMethod) IsTemporarilyLocked() bool {
	return a.Status == uint8(StatusTemporarilyLocked)
}

// IsPermanentlyLocked kiểm tra tài khoản có bị khóa vĩnh viễn không
func (a *AuthMethod) IsPermanentlyLocked() bool {
	return a.Status == uint8(StatusPermanentlyLocked)
}

// IsLockExpired kiểm tra khóa tạm thời có hết hạn không
func (a *AuthMethod) IsLockExpired() bool {
	if !a.IsTemporarilyLocked() || a.LockedUntil == nil {
		return false
	}
	return time.Now().After(*a.LockedUntil)
}

// CanLogin kiểm tra tài khoản có thể đăng nhập không
func (a *AuthMethod) CanLogin() bool {
	if a.Status == uint8(StatusInactive) {
		return false
	}
	if a.IsPermanentlyLocked() {
		return false
	}
	if a.IsTemporarilyLocked() && !a.IsLockExpired() {
		return false
	}
	return true
}
