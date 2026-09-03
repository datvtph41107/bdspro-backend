package models

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	_model "common/models"
	"fmt"
	"time"
	"user/enums"

	"github.com/lib/pq"
)

type UserProfileSearch struct {
	ProfileID uint64 `gorm:"primaryKey;column:profile_id" copier:"-" json:"profileId"`
	// Email     string `gorm:"size:50"`
	FullName string `gorm:"size:255" json:"fullName"`
	Avatar   string `gorm:"size:255" json:"avatar"`
}

// UserProfileEntity chứa thông tin hồ sơ người dùng
type UserProfileEntity struct {
	_model.BaseEntityNotId
	ProfileID         uint64                `gorm:"primaryKey;autoIncrement;column:profile_id" copier:"-" json:"profileId"`
	Phone             string                `gorm:"size:20" json:"phone"`
	Phone2            string                `gorm:"size:20" json:"phone2"`
	Email             string                `gorm:"size:50" json:"email"`
	TaxCode           string                `gorm:"size:13" json:"taxCode"`
	FullName          string                `gorm:"size:255" json:"fullName"`
	Address           string                `gorm:"size:255" json:"address"`
	Gender            uint32                `json:"gender"`
	Avatar            string                `gorm:"size:255" json:"avatar"`
	TickVerified      bool                  `gorm:"column:tick_verified;default:false" json:"tickVerified"`
	BackgroundImage   string                `gorm:"size:255" json:"backgroundImage"`
	ReferralCode      string                `gorm:"size:255" copier:"-" json:"referralCode"`
	ReferralBy        string                `gorm:"size:255" copier:"-" json:"referralBy"`
	RoleType          enums.ERole           `gorm:"default:10" copier:"-" json:"roleType"`
	RoleKey           uint32                `gorm:"default:10" copier:"-" json:"roleKey"`
	RoleRealEstate    enums.ERoleRealEstate `gorm:"default:10" copier:"-" json:"roleRealEstate"`
	Position          string                `gorm:"size:255" json:"position"`
	Workplace         string                `gorm:"size:255;column:workplace" json:"workplace"`
	RoleTitle         string                `gorm:"size:255;column:role_title" json:"roleTitle"`
	DepartmentID      uint64                `gorm:"column:department_id" json:"departmentId"`
	Birth             *time.Time            `gorm:"column:birth" json:"birth"`
	StartDate         *time.Time            `gorm:"column:start_date" json:"startDate"`
	DeletedAt         *time.Time            `gorm:"column:deleted_at" json:"-"`
	FrontIdentify     string                `gorm:"column:front_identify" json:"frontIdentify"`
	BackIdentify      string                `gorm:"column:back_identify" json:"backIdentify"`
	ProvinceID        *uint64               `gorm:"column:province_id" json:"provinceId"`
	WardID            *uint64               `gorm:"column:ward_id" json:"wardId"`
	PlanID            *uint64               `gorm:"column:plan_id" json:"planId"`
	PlanAt            *time.Time            `gorm:"column:plan_at" json:"planAt"`
	Slogan            string                `gorm:"column:slogan" json:"slogan"`
	WebsiteURL        string                `gorm:"column:website_url" json:"websiteUrl"`
	ZaloURL           string                `gorm:"column:zalo_url" json:"zaloUrl"`
	FacebookURL       string                `gorm:"column:facebook_url" json:"facebookUrl"`
	InstagramURL      string                `gorm:"column:instagram_url" json:"instagramUrl"`
	TwitterURL        string                `gorm:"column:twitter_url" json:"twitterUrl"`
	LinkedinURL       string                `gorm:"column:linkedin_url" json:"linkedinUrl"`
	YoutubeURL        string                `gorm:"column:youtube_url" json:"youtubeUrl"`
	Introduction      string                `gorm:"column:introduction" json:"introduction"`
	Visibility        _enum.EVisibility     `gorm:"column:visibility" json:"visibility"`                // ai xem được trang cá nhân
	StatusOnline      enums.EOnlineStatus   `gorm:"column:status_online" json:"statusOnline"`           // online, gần đây, không hiển thị
	ProfileVisibility _enum.EVisibility     `gorm:"column:profile_visibility" json:"profileVisibility"` // ai xem được liên hệ
	TabDefault        enums.ESettingTab     `gorm:"column:tab_default" json:"tabDefault"`
	SignatureVisible  bool                  `gorm:"column:signature_visible;default:false" json:"signatureVisible"`

	VisibilityIntroduce  _enum.EVisibility `gorm:"column:visibility_introduce" json:"visibilityIntroduce"`
	VisibilityProfession _enum.EVisibility `gorm:"column:visibility_profession" json:"visibilityProfession"`
	VisibilityMainArea   _enum.EVisibility `gorm:"column:visibility_main_area" json:"visibilityMainArea"`
	VisibilityFriends    _enum.EVisibility `gorm:"column:visibility_friends" json:"visibilityFriends"`
	VisibilitySignature  _enum.EVisibility `gorm:"column:visibility_signature" json:"visibilitySignature"`
	ViewRoles            pq.Int32Array     `gorm:"type:integer[];column:view_roles" json:"viewRoles"` // Array of EViewRole (10, 20, 30, 40)

	// LastLoginAt       *time.Time            `gorm:"column:last_login_at" json:"lastLoginAt"`
	// TotalLogin        uint32                `gorm:"column:total_login;default:0" json:"totalLogin"`

	// transient
	Certifications  []CertificationEntity `gorm:"-" json:"certifications"`
	Professions     []ProfessionEntity    `gorm:"-" json:"professions"`
	PurposeUses     []PurposeUseEntity    `gorm:"-" json:"purposeUses"`
	TagJobs         []TagEntity           `gorm:"-" json:"tagJobs"`
	TagSpecialities []TagEntity           `gorm:"-" json:"tagSpecialities"`
	TagMainAreas    []TagEntity           `gorm:"-" json:"tagMainAreas"`
	TagProjects     []TagEntity           `gorm:"-" json:"tagProjects"`
	Tags            []TagEntity           `gorm:"-" json:"tags"`
	ProfileMedias   []ProfileMediaEntity  `gorm:"-" json:"profileMedias"`
	RateStats       *_dto.RateStats       `gorm:"-" json:"rateStats"`

	FriendStatus *_dto.FriendItemDTO `gorm:"-" json:"friendStatus"`
	Following    bool                `gorm:"-" json:"following"`
	Status       _enum.EUserStatus   `gorm:"column:status;default:10" json:"status"`

	Bookmark bool `gorm:"-" json:"bookmark"`

	LastLoginAt *time.Time `gorm:"column:last_login_at" json:"lastLoginAt"`
	TotalLogin  uint32     `gorm:"column:total_login" json:"totalLogin"`

	// WarningLevel level
	WarningLevel _enum.EWarningLevel `gorm:"column:warning;default:0" json:"warning"`

	// Role information
	RoleID *uint64 `gorm:"column:role_id" json:"roleId"`

	// transient - Role information from auth-service
	RoleName    string `gorm:"-" json:"roleName"`
	RoleKeyData uint32 `gorm:"-" json:"roleKeyData"` // Role key from role table
	RoleColor   string `gorm:"-" json:"roleColor"`
	RoleBgColor string `gorm:"-" json:"roleBgColor"`

	KYC *KYCEntity `gorm:"-" json:"kyc"`
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
	profile := &UserProfileEntity{
		ReferralCode: generateReferral(authId),
		RoleType:     enums.ERoleUser,
		FullName:     fullname,
		Email:        email,
		Phone:        phone,
		Avatar:       avatar,
		PlanID:       planID,
	}
	// Set tất cả visibility thành public mặc định
	SetAllVisibilityPublic(profile)
	return profile
}

// SetAllVisibilityPublic set tất cả visibility fields thành public (10)
func SetAllVisibilityPublic(profile *UserProfileEntity) {
	if profile == nil {
		return
	}
	// Set visibility chung
	profile.Visibility = _enum.EVisibilityDTOPublic
	profile.ProfileVisibility = _enum.EVisibilityDTOPublic

	// Set các visibility riêng lẻ
	profile.VisibilityIntroduce = _enum.EVisibilityDTOPublic
	profile.VisibilityProfession = _enum.EVisibilityDTOPublic
	profile.VisibilityMainArea = _enum.EVisibilityDTOPublic
	profile.VisibilityFriends = _enum.EVisibilityDTOPublic
	profile.VisibilitySignature = _enum.EVisibilityDTOPublic
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
