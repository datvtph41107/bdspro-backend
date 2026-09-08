package dto

import (
	"encoding/json"
	"time"

	_dto "common/domain/dto"
	_enum "common/domain/enum"
	"user/enums"
	models "user/internal/models"
)

type FollowInfoRes struct {
	NumFollower  int64 `json:"numFollower"`
	NumFollowing int64 `json:"numFollowing"`
	NumFriend    int64 `json:"numFriend"`
}

// MarshalJSON ẩn RefreshToken khi serialize JSON
func (a AuthLoginResponse) MarshalJSON() ([]byte, error) {
	type Alias AuthLoginResponse
	return json.Marshal(&struct {
		Alias
		RefreshToken string `json:"refreshToken,omitempty"`
	}{
		Alias:        (Alias)(a),
		RefreshToken: "",
	})
}

// RefreshResponse đại diện cho phản hồi làm mới token
type RefreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// NewRefreshResponse tạo mới RefreshResponse
func NewRefreshResponse(accessToken string) RefreshResponse {
	return RefreshResponse{AccessToken: accessToken}
}

// RequestLoginResponse đại diện cho phản hồi yêu cầu đăng nhập
type RequestLoginResponse struct {
	AuthID     uint64 `json:"authId,omitempty"`
	Provider   string `json:"provider,omitempty"`
	OAuthID    string `json:"oauthId,omitempty"`
	Existed    bool   `json:"existed"`    // true: user đã tồn tại, false: user mới đăng ký
	SmsChannel string `json:"smsChannel"` // "zalo" nếu gửi ZNS thành công, "sms" nếu ZNS lỗi
}

// UserInfoResponse đại diện cho phản hồi thông tin người dùng
//
//	type UserInfoResponse struct {
//		Phone        string     `json:"phone,omitempty"`
//		Email        string     `json:"email,omitempty"`
//		TaxCode      string     `json:"taxCode,omitempty" validate:"omitempty,min=2,max=50"`
//		Fullname     string     `json:"fullname,omitempty" validate:"omitempty,min=2,max=50"`
//		Address      string     `json:"address,omitempty" validate:"omitempty,min=2,max=100"`
//		Avatar       string     `json:"avatar,omitempty" validate:"omitempty,url,max=255"`
//		Position     string     `json:"position,omitempty" validate:"omitempty,min=2,max=50"`
//		DepartmentID uint64     `json:"department_id,omitempty" validate:"omitempty,gte=1"`
//		Birth        *time.Time `json:"birth,omitempty" validate:"omitempty"`
//	}
type UserInfoResponse struct {
	ProfileID       uint64     `json:"profileId"`
	Phone           string     `json:"phone,omitempty"`
	Email           string     `json:"email,omitempty"`
	TaxCode         string     `json:"taxCode,omitempty"`
	FullName        string     `json:"fullname,omitempty"`
	Address         string     `json:"address,omitempty"`
	Avatar          string     `json:"avatar,omitempty"`
	ReferralCode    string     `json:"referralCode,omitempty"`
	ReferralBy      string     `json:"referralBy,omitempty"`
	RoleType        string     `json:"roleType,omitempty"`
	Position        string     `json:"position,omitempty"`
	DepartmentID    uint64     `json:"departmentId,omitempty"`
	Birth           *time.Time `json:"birth,omitempty"`
	StartDate       *time.Time `json:"startDate,omitempty"`
	FollowingNumber int        `json:"followingNumber,omitempty"`
	FollowerNumber  int        `json:"followerNumber,omitempty"`
	FriendNumber    int        `json:"friendNumber,omitempty"`
}

type PublicUserInfo struct {
	ProfileID         uint64                `json:"profileId,omitempty"`
	FullName          string                `json:"fullname,omitempty"`
	Avatar            string                `json:"avatar,omitempty"`
	Phone             string                `json:"phone,omitempty"`
	Phone2            string                `json:"phone2,omitempty"`
	Email             string                `json:"email,omitempty"`
	TaxCode           string                `json:"taxCode,omitempty"`
	Address           string                `json:"address,omitempty"`
	Gender            uint8                 `json:"gender,omitempty"`
	TickVerified      bool                  `json:"tickVerified,omitempty"`
	BackgroundImage   string                `json:"backgroundImage,omitempty"`
	ReferralCode      string                `json:"referralCode,omitempty"`
	ReferralBy        string                `json:"referralBy,omitempty"`
	RoleType          enums.ERole           `json:"roleType,omitempty"`
	RoleKey           uint32                `json:"roleKey,omitempty"`
	RoleRealEstate    enums.ERoleRealEstate `json:"roleRealEstate,omitempty"`
	Position          string                `json:"position,omitempty"`
	Workplace         string                `json:"workplace,omitempty"`
	DepartmentID      uint64                `json:"departmentId,omitempty"`
	Birth             *time.Time            `json:"birth,omitempty"`
	StartDate         *time.Time            `json:"startDate,omitempty"`
	FrontIdentify     string                `json:"frontIdentify,omitempty"`
	BackIdentify      string                `json:"backIdentify,omitempty"`
	PlanID            *uint64               `json:"planId,omitempty"`
	PlanAt            *time.Time            `json:"planAt,omitempty"`
	Slogan            string                `json:"slogan,omitempty"`
	Website           string                `json:"website,omitempty"`
	Facebook          string                `json:"facebook,omitempty"`
	Instagram         string                `json:"instagram,omitempty"`
	Twitter           string                `json:"twitter,omitempty"`
	Linkedin          string                `json:"linkedin,omitempty"`
	Youtube           string                `json:"youtube,omitempty"`
	Introduction      string                `json:"introduction,omitempty"`
	Visibility        uint8                 `json:"visibility,omitempty"`
	StatusOnline      enums.EOnlineStatus   `json:"statusOnline,omitempty"`
	ProfileVisibility uint8                 `json:"profileVisibility,omitempty"`
	TabDefault        enums.ESettingTab     `json:"tabDefault,omitempty"`
	Status            _enum.EUserStatus     `json:"status,omitempty"`
	RoleID            *uint64               `json:"roleId,omitempty"`
	RoleName          string                `json:"roleName,omitempty"`
	RoleColor         string                `json:"roleColor,omitempty"`
	RoleBgColor       string                `json:"roleBgColor,omitempty"`

	// transient
	Certifications []models.CertificationEntity `json:"certifications,omitempty"`
	Professions    []models.ProfessionEntity    `json:"professions,omitempty"`
	PurposeUses    []models.PurposeUseEntity    `json:"purposeUses,omitempty"`
	MainAreas      []models.MainAreaEntity      `json:"mainAreas,omitempty"`
	FriendStatus   *_dto.FriendItemDTO          `json:"relationShip,omitempty"`
	Following      bool                         `json:"following,omitempty"`
	Bookmark       bool                         `json:"bookmark,omitempty"`
	RateStats      *_dto.RateStats              `json:"rateStats,omitempty"`
}

type RelationShip struct {
	FriendStatus *_dto.FriendItemDTO `json:"friendStatus,omitempty"`
	IsFollowing  bool                `json:"isFollowing,omitempty"`
}

type JwtTokenProperties struct {
	AuthID    uint64
	ProfileID uint64
	SessionID uint64
	Role      string
	Type      string
}

// type ProfileResponse struct {

// }

// // NewUserInfoResponse tạo mới từ UserProfileEntity
// func MakeUserInfoResponse(e *models.UserProfileEntity) *UserInfoResponse {
// 	if e == nil {
// 		return nil
// 	}
// 	var response UserInfoResponse
// 	copier.Copy(&response, e)
// 	return &response
// 	// return UserInfoResponse{
// 	// 	Phone:    e.Phone,
// 	// 	Email:    e.Email,
// 	// 	Fullname: e.FullName,
// 	// }
// }

// ReportResponse đại diện cho phản hồi báo cáo
type ReportResponse struct {
	ID           uint64                `json:"id"`
	ReasonID     uint64                `json:"reasonId"`
	Reason       ReportReasonResponse  `json:"reason"`
	ReportStatus uint32                `json:"reportStatus"`
	StatusText   string                `json:"statusText"`
	UserID       uint64                `json:"userId"`
	User         ReportUserResponse    `json:"user"`
	ProofDocs    []ReportProofResponse `json:"proofDocs"`
	Response     string                `json:"response,omitempty"`
	CreatedAt    *time.Time            `json:"createdAt"`
	UpdatedAt    *time.Time            `json:"updatedAt"`
}

// ReportReasonResponse đại diện cho phản hồi nguyên nhân báo cáo
type ReportReasonResponse struct {
	ID     uint64 `json:"id"`
	Reason string `json:"reason"`
}

// ReportUserResponse đại diện cho phản hồi thông tin user trong báo cáo
type ReportUserResponse struct {
	ProfileId uint64 `json:"profileId"`
	Phone     string `json:"phone"`
	FullName  string `json:"fullName"`
	Avatar    string `json:"avatar"`
}

// ReportProofResponse đại diện cho phản hồi bằng chứng báo cáo
type ReportProofResponse struct {
	ID       uint64 `json:"id"`
	ReportID uint64 `json:"reportId"`
	FileName string `json:"fileName"`
	FileUrl  string `json:"fileUrl"`
	FileType string `json:"fileType"`
}

// ReportListResponse đại diện cho phản hồi danh sách báo cáo
type ReportListResponse struct {
	Reports []ReportResponse `json:"reports"`
	Total   int64            `json:"total"`
	Page    int              `json:"page"`
	Size    int              `json:"size"`
}

// AdminNoteResponse đại diện cho phản hồi cập nhật ghi chú admin
// @swagger:model
type AdminNoteResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Report  *ReportResponse `json:"report,omitempty"`
}

type ProfileCompletionResponse struct {
	Percentage      float64  `json:"percentage"`
	CompletedFields uint32   `json:"completedFields"`
	TotalFields     uint32   `json:"totalFields"`
	MissingFields   []string `json:"missingFields"`
}
