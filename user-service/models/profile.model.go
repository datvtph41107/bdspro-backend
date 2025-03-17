package models

import (
	_db "common/db"
	"fmt"
	"time"
)

// RoleEnum định nghĩa các vai trò người dùng
type RoleEnum string

const (
	RoleUser  RoleEnum = "USER"
	RoleAdmin RoleEnum = "ADMIN"
)

// UserProfileEntity chứa thông tin hồ sơ người dùng
type UserProfileEntity struct {
	ProfileID    uint64     `gorm:"primaryKey;column:profile_id"`
	Phone        string     `gorm:"size:20"`
	Email        string     `gorm:"size:50"`
	TaxCode      string     `gorm:"size:13"`
	FullName     string     `gorm:"size:255"`
	Address      string     `gorm:"size:255"`
	Avatar       string     `gorm:"size:255"`
	ReferralCode string     `gorm:"size:255"`
	ReferralBy   string     `gorm:"size:255"`
	RoleType     RoleEnum   `gorm:"type:varchar(10)"`
	Position     string     `gorm:"size:255"`
	DepartmentID uint64     `gorm:"column:department_id"`
	Birth        *time.Time `gorm:"column:birth"`
	StartDate    *time.Time `gorm:"column:start_date"`
	// CreatedAt    time.Time      `gorm:"autoCreateTime"`
	// UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	// DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// TableName đặt tên bảng trong DB
func (UserProfileEntity) TableName() string {
	return "user_profile"
}

// NewUserProfileEntity tạo user với authId
func NewUserProfileEntity(authId uint64) *UserProfileEntity {
	return &UserProfileEntity{
		ReferralCode: generateReferral(authId),
		RoleType:     RoleUser,
	}
}

// NewUserProfileEntityWithRequest tạo user từ authId và OtpVerifyRequest
func NewUserProfileEntityWithRequest(authId uint64, request OtpVerifyRequest) *UserProfileEntity {
	return &UserProfileEntity{
		ReferralCode: generateReferral(authId),
		RoleType:     RoleUser,
		FullName:     request.FullName,
	}
}

// Hàm giả lập sinh mã giới thiệu
func generateReferral(authId uint64) string {
	return "REF" + fmt.Sprintf("%06d", authId)
}

// OtpVerifyRequest mô phỏng request object
type OtpVerifyRequest struct {
	FullName string
}

// GetByProfileId truy vấn UserProfileEntity theo profileId
func ProfileGetByProfileId(profileId uint64) (*UserProfileEntity, error) {
	var userProfile UserProfileEntity
	if err := _db.DB.Where("profile_id = ?", profileId).First(&userProfile).Error; err != nil {
		return nil, err
	}
	return &userProfile, nil
}
