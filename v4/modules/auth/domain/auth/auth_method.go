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

// GORM hook để validate
func (a *AuthMethod) BeforeCreate(tx *gorm.DB) error {
	if !a.Provider.IsValid() {
		return fmt.Errorf("invalid provider type: %s", a.Provider)
	}
	return nil
}

// AuthMethod is a login credential. Account lock/status lives on user_profile.
type AuthMethod struct {
	ID          uint64       `gorm:"primaryKey;autoIncrement"`
	Provider    ProviderType `gorm:"column:provider;type:varchar(20);not null"`
	AuthName    string       `gorm:"column:auth_name;size:255;not null"`
	Password    string       `gorm:"column:password"`
	IsSensitive bool         `gorm:"column:is_sensitive"`
	Avatar      string       `gorm:"size:255"`
	FullName    string       `gorm:"size:255"`
	Email       string       `gorm:"size:255"`
	Phone       string       `gorm:"size:255"`
	RoleKey     uint32       `gorm:"column:role_key"`
	PrivateKey  string       `gorm:"column:private_key;size:512"`
	PublicKey   string       `gorm:"column:public_key;size:512"`
	AuthKey     string       `gorm:"column:auth_key;size:512"`
	CreatedAt   time.Time    `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time    `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   *time.Time   `gorm:"column:deleted_at"`
	UserID      uint64       `gorm:"column:user_id"`
}

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
