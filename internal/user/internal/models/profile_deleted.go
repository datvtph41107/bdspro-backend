package models

import (
	_enum "common/domain/enum"
	"time"
	"user/enums"
)

// ProfileDeleted lưu thông tin profile đã bị xóa
type ProfileDeleted struct {
	ID                uint64                `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	ProfileID         uint64                `gorm:"column:profile_id;index" json:"profileId"`
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
	ReferralCode      string                `gorm:"size:255" json:"referralCode"`
	ReferralBy        string                `gorm:"size:255" json:"referralBy"`
	RoleType          enums.ERole           `gorm:"default:10" json:"roleType"`
	RoleKey           uint32                `gorm:"default:10" json:"roleKey"`
	RoleRealEstate    enums.ERoleRealEstate `gorm:"default:10" json:"roleRealEstate"`
	Position          string                `gorm:"size:255" json:"position"`
	DepartmentID      uint64                `gorm:"column:department_id" json:"departmentId"`
	Birth             *time.Time            `gorm:"column:birth" json:"birth"`
	StartDate         *time.Time            `gorm:"column:start_date" json:"startDate"`
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
	Visibility        _enum.EVisibility     `gorm:"column:visibility" json:"visibility"`
	StatusOnline      enums.EOnlineStatus   `gorm:"column:status_online" json:"statusOnline"`
	ProfileVisibility _enum.EVisibility     `gorm:"column:profile_visibility" json:"profileVisibility"`
	TabDefault        enums.ESettingTab     `gorm:"column:tab_default" json:"tabDefault"`
	Status            _enum.EUserStatus     `gorm:"column:status;default:10" json:"status"`
	WarningLevel      _enum.EWarningLevel   `gorm:"column:warning;default:0" json:"warning"`
	RoleID            *uint64               `gorm:"column:role_id" json:"roleId"`

	// Metadata về việc xóa
	DeletedBy         *uint64    `gorm:"column:deleted_by" json:"deletedBy"`                   // Người thực hiện xóa
	DeletedReason     string     `gorm:"column:deleted_reason;type:text" json:"deletedReason"` // Lý do xóa
	DeletedAt         time.Time  `gorm:"column:deleted_at" json:"deletedAt"`                   // Thời điểm xóa
	OriginalCreatedAt *time.Time `gorm:"column:original_created_at" json:"originalCreatedAt"`  // Thời điểm tạo ban đầu
	OriginalUpdatedAt *time.Time `gorm:"column:original_updated_at" json:"originalUpdatedAt"`  // Thời điểm cập nhật cuối cùng
}

// TableName đặt tên bảng trong DB
func (ProfileDeleted) TableName() string {
	return "profile_deleted"
}

// NewProfileDeleted tạo ProfileDeleted từ UserProfileEntity
func NewProfileDeleted(profile *UserProfileEntity, deletedBy *uint64, reason string) *ProfileDeleted {
	now := time.Now()
	return &ProfileDeleted{
		ProfileID:         profile.ProfileID,
		Phone:             profile.Phone,
		Phone2:            profile.Phone2,
		Email:             profile.Email,
		TaxCode:           profile.TaxCode,
		FullName:          profile.FullName,
		Address:           profile.Address,
		Gender:            profile.Gender,
		Avatar:            profile.Avatar,
		TickVerified:      profile.TickVerified,
		BackgroundImage:   profile.BackgroundImage,
		ReferralCode:      profile.ReferralCode,
		ReferralBy:        profile.ReferralBy,
		RoleType:          profile.RoleType,
		RoleKey:           profile.RoleKey,
		RoleRealEstate:    profile.RoleRealEstate,
		Position:          profile.Position,
		DepartmentID:      profile.DepartmentID,
		Birth:             profile.Birth,
		StartDate:         profile.StartDate,
		FrontIdentify:     profile.FrontIdentify,
		BackIdentify:      profile.BackIdentify,
		PlanID:            profile.PlanID,
		PlanAt:            profile.PlanAt,
		Slogan:            profile.Slogan,
		Website:           profile.WebsiteURL,
		Facebook:          profile.FacebookURL,
		Instagram:         profile.InstagramURL,
		Twitter:           profile.TwitterURL,
		Linkedin:          profile.LinkedinURL,
		Youtube:           profile.YoutubeURL,
		Introduction:      profile.Introduction,
		Visibility:        profile.Visibility,
		StatusOnline:      profile.StatusOnline,
		ProfileVisibility: profile.ProfileVisibility,
		TabDefault:        profile.TabDefault,
		Status:            profile.Status,
		WarningLevel:      profile.WarningLevel,
		RoleID:            profile.RoleID,
		DeletedBy:         deletedBy,
		DeletedReason:     reason,
		DeletedAt:         now,
		OriginalCreatedAt: profile.CreatedAt,
		OriginalUpdatedAt: profile.UpdatedAt,
	}
}
