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
	ProfileID     uint64     `gorm:"primaryKey;column:profile_id" copier:"-"`
	Phone         string     `gorm:"size:20"`
	Email         string     `gorm:"size:50"`
	TaxCode       string     `gorm:"size:13"`
	FullName      string     `gorm:"size:255"`
	Address       string     `gorm:"size:255"`
	Avatar        string     `gorm:"size:255"`
	ReferralCode  string     `gorm:"size:255" copier:"-"`
	ReferralBy    string     `gorm:"size:255" copier:"-"`
	RoleType      RoleEnum   `gorm:"type:varchar(10)" copier:"-"`
	Position      string     `gorm:"size:255"`
	DepartmentID  uint64     `gorm:"column:department_id"`
	Birth         *time.Time `gorm:"column:birth"`
	StartDate     *time.Time `gorm:"column:start_date"`
	DeletedAt     *time.Time `gorm:"column:deleted_at"`
	FrontIdentify string     `gorm:"column:front_identify"`
	BackIdentify  string     `gorm:"column:back_identify"`
	PlanID        *uint64    `gorm:"column:plan_id"`
	// CreatedAt    time.Time      `gorm:"autoCreateTime"`
	// UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	// DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// TableName đặt tên bảng trong DB
func (UserProfileEntity) TableName() string {
	return "user_profile"
}

// NewUserProfileEntity tạo user với authId
func NewUserProfileEntity(authId uint64, planID *uint64, fullname, email, phone, avatar string) *UserProfileEntity {
	return &UserProfileEntity{
		ReferralCode: generateReferral(authId),
		RoleType:     RoleUser,
		FullName:     fullname,
		Email:        email,
		Phone:        phone,
		Avatar:       avatar,
		PlanID:       planID,
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
	if err := _db.DB.Where("profile_id = ? and deleted_at is null", profileId).First(&userProfile).Error; err != nil {
		return nil, err
	}
	return &userProfile, nil
}

// GetByProfileId truy vấn UserProfileEntity theo profileId
func ProfileGetByPhone(phone string) (*UserProfileEntity, error) {
	var userProfile *UserProfileEntity
	if err := _db.DB.Model(&UserProfileEntity{}).
		Where("(phone like '%' || ? || '%') and deleted_at is null", phone).
		Scan(&userProfile).Error; err != nil {
		return nil, err
	}
	return userProfile, nil
}

// GetByProfileId truy vấn UserProfileEntity theo profileId
func ProfileSearchForChat(phone string) (*[]UserProfileEntity, error) {
	var userProfile *[]UserProfileEntity

	query := "%" + phone + "%"
	if err := _db.DB.Debug().
		Where("(phone ILIKE ? OR full_name ILIKE ?) AND deleted_at IS NULL", query, query).
		// Where("(phone like '%' || ? || '%' or name like '%' || ? || '%') and deleted_at is null", phone, phone).
		Find(&userProfile).Error; err != nil {
		return nil, err
	}
	return userProfile, nil
}

// GetByProfileId truy vấn UserProfileEntity theo profileId
func ProfileByIDs(ids []uint64) (*[]UserProfileEntity, error) {
	var userProfile *[]UserProfileEntity

	// query := "%" + phone + "%"
	if err := _db.DB.Debug().
		Where("profile_id in (?) AND deleted_at IS NULL", ids).
		Find(&userProfile).Error; err != nil {
		return nil, err
	}
	return userProfile, nil
}
