package models

import (
	_db "common/db"
	_enum "common/domain/enum"
	_model "common/models"
	"fmt"
	"time"
	"user/enums"
)

type UserProfileSearch struct {
	ProfileID uint64 `gorm:"primaryKey;column:profile_id" copier:"-" json:"profileId"`
	// Email     string `gorm:"size:50"`
	FullName string `gorm:"size:255" json:"fullName"`
	Avatar   string `gorm:"size:255" json:"avatar"`
}

// UserProfileEntity chứa thông tin hồ sơ người dùng
type UserProfileEntity struct {
	_model.BaseEntity
	ProfileID         uint64                `gorm:"primaryKey;autoIncrement;column:profile_id" copier:"-" json:"profileId"`
	Phone             string                `gorm:"size:20" json:"phone"`
	Phone2            string                `gorm:"size:20" json:"phone2"`
	Email             string                `gorm:"size:50" json:"email"`
	TaxCode           string                `gorm:"size:13" json:"taxCode"`
	FullName          string                `gorm:"size:255" json:"fullName"`
	Address           string                `gorm:"size:255" json:"address"`
	Gender            uint8                 `json:"gender"`
	Avatar            string                `gorm:"size:255" json:"avatar"`
	TickVerified      bool                  `gorm:"column:tick_verified;default:false" json:"tickVerified"`
	BackgroundImage   string                `gorm:"size:255" json:"backgroundImage"`
	ReferralCode      string                `gorm:"size:255" copier:"-" json:"referralCode"`
	ReferralBy        string                `gorm:"size:255" copier:"-" json:"referralBy"`
	RoleType          enums.ERole           `gorm:"default:10" copier:"-" json:"roleType"`
	RoleKey           uint32                `gorm:"default:10" copier:"-" json:"roleKey"`
	RoleRealEstate    enums.ERoleRealEstate `gorm:"default:10" copier:"-" json:"roleRealEstate"`
	Position          string                `gorm:"size:255" json:"position"`
	DepartmentID      uint64                `gorm:"column:department_id" json:"departmentId"`
	Birth             *time.Time            `gorm:"column:birth" json:"birth"`
	StartDate         *time.Time            `gorm:"column:start_date" json:"startDate"`
	DeletedAt         *time.Time            `gorm:"column:deleted_at" json:"-"`
	FrontIdentify     string                `gorm:"column:front_identify" json:"frontIdentify"`
	BackIdentify      string                `gorm:"column:back_identify" json:"backIdentify"`
	PlanID            *uint64               `gorm:"column:plan_id" json:"planId"`
	PlanAt            *time.Time            `gorm:"column:plan_at" json:"planAt"`
	Slogan            string                `gorm:"column:slogan" json:"slogan"`
	Website           string                `gorm:"column:website" json:"website"`
	Facebook          string                `gorm:"column:facebook" json:"facebook"`
	Instagram         string                `gorm:"column:instagram" json:"instagram"`
	Twitter           string                `gorm:"column:twitter" json:"twitter"`
	Linkedin          string                `gorm:"column:linkedin" json:"linkedin"`
	Youtube           string                `gorm:"column:youtube" json:"youtube"`
	Introduction      string                `gorm:"column:introduction" json:"introduction"`
	Visibility        uint8                 `gorm:"column:visibility" json:"visibility"`
	StatusOnline      enums.EOnlineStatus   `gorm:"column:status_online" json:"statusOnline"`
	ProfileVisibility uint8                 `gorm:"column:profile_visibility" json:"profileVisibility"`
	TabDefault        enums.ESettingTab     `gorm:"column:tab_default" json:"tabDefault"`
	// LastLoginAt       *time.Time            `gorm:"column:last_login_at" json:"lastLoginAt"`
	// TotalLogin        uint32                `gorm:"column:total_login;default:0" json:"totalLogin"`

	// transient
	Certifications []CertificationEntity `gorm:"-" json:"certifications"`
	Professions    []ProfessionEntity    `gorm:"-" json:"professions"`
	RateStats      *RateStats            `gorm:"-" json:"rateStats"`

	FriendStatus *FriendItem       `gorm:"-" json:"friendStatus"`
	Following    bool              `gorm:"-" json:"following"`
	Status       _enum.EUserStatus `gorm:"column:status;default:10" json:"status"`

	Bookmark bool `gorm:"-" json:"bookmark"`

	// WarningLevel level
	WarningLevel _enum.EWarningLevel `gorm:"column:warning;default:0" json:"warning"`

	// Role information
	RoleID *uint64 `gorm:"column:role_id" json:"roleId"`

	// transient - Role information from auth-service
	RoleName    string `gorm:"-" json:"roleName"`
	RoleKeyData uint32 `gorm:"-" json:"roleKeyData"` // Role key from role table
	RoleColor   string `gorm:"-" json:"roleColor"`
	RoleBgColor string `gorm:"-" json:"roleBgColor"`
}

// TableName đặt tên bảng trong DB
func (UserProfileEntity) TableName() string {
	return "user_profile"
}

type UserItem struct {
	ProfileID    uint64 `gorm:"primaryKey;column:profile_id" copier:"-" json:"profileId"`
	FullName     string `gorm:"size:255" json:"fullName"`
	Avatar       string `gorm:"size:255" json:"avatar"`
	TickVerified bool   `gorm:"column:tick_verified" json:"tickVerified"`
}

func (UserItem) TableName() string {
	return "user_profile"
}

// MakeUserProfileEntity tạo user với authId
func MakeUserProfileEntity(authId uint64, planID *uint64, fullname, email, phone, avatar string) *UserProfileEntity {
	return &UserProfileEntity{
		ReferralCode: generateReferral(authId),
		RoleType:     enums.ERoleUser,
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
		RoleType:     enums.ERoleUser,
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
func ProfileGetByPhone(phone string) (*UserProfileEntity, error) {
	var userProfile UserProfileEntity
	if err := _db.DB.Model(&UserProfileEntity{}).
		Where("(phone = ?) and deleted_at is null", phone).
		First(&userProfile).Error; err != nil {
		return nil, err
	}
	return &userProfile, nil
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
