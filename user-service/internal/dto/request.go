package dto

import (
	"time"

	_dto "common/domain/dto"
	"user/enums"
	"user/internal/models"

	"github.com/go-playground/validator/v10"
)

// OtpRequest đại diện cho yêu cầu gửi OTP
// @swagger:model
type OtpRequest struct {
	Phone    string `json:"phone" binding:"required" validate:"required,e164" example:"0901230456"`
	Fullname string `json:"fullname,omitempty"`
}

// OtpVerifyRequest mở rộng từ OtpRequest để xác minh OTP
// @swagger:model
type OtpVerifyRequest struct {
	OtpRequest
	Otp      string `json:"otp" binding:"required" validate:"required"`
	AuthID   int64  `json:"authId,omitempty"`
	Version  string `json:"version,omitempty"`
	Platform string `json:"platform,omitempty"`
	OS       string `json:"os,omitempty"`
	Fullname string `json:"fullname,omitempty"`
}

// RefreshTokenRequest đại diện cho yêu cầu làm mới token
// @swagger:model
type RefreshTokenRequest struct {
	AccessToken  string `json:"accessToken" binding:"required" validate:"required"`
	RefreshToken string `json:"refreshToken" binding:"required" validate:"required"`
}

// RestorePhoneRequest đại diện cho yêu cầu khôi phục số điện thoại
// @swagger:model
type RestorePhoneRequest struct {
	Phone string `json:"phone" binding:"required" validate:"required,e164"`
}

type SessionQRRequest struct {
	ClientID   string `json:"clientId"`   // ID của trình duyệt hoặc phiên web
	UserAgent  string `json:"userAgent"`  // Chuỗi nhận diện trình duyệt
	CsrfToken  string `json:"csrfToken"`  // Token chống CSRF (nếu cần)
	Expiration int    `json:"expiration"` // Thời gian hết hạn QR code (tính bằng giây)
	SessionKey string `json:"sessionKey"`
	SessionId  uint64 `json:"sessionId"`
}

type SessionQRResponse struct {
	SessionKey string `json:"sessionKey"`
	SessionId  uint64 `json:"sessionId"`
}
type SessionQRConfirm struct {
	SessionId uint64 `json:"sessionId"` // ID của trình duyệt hoặc phiên web
}

// UserInfoRequest đại diện cho yêu cầu cập nhật thông tin người dùng
// @swagger:model
type UserInfoRequest struct {
	Phone             string                `json:"phone,omitempty"`
	Phone2            string                `json:"phone2,omitempty"`
	Email             string                `json:"email,omitempty"`
	TaxCode           string                `json:"taxCode,omitempty"`
	FullName          string                `json:"fullname,omitempty"`
	Address           string                `json:"address,omitempty"`
	Avatar            string                `json:"avatar,omitempty"`
	BackgroundImage   string                `json:"backgroundImage,omitempty"`
	Position          string                `json:"position,omitempty"`
	Workplace         string                `json:"workplace,omitempty"`
	RoleTitle         string                `json:"roleTitle,omitempty"`
	DepartmentID      uint64                `json:"departmentId,omitempty"`
	Birth             *time.Time            `json:"birth,omitempty" time_format:"2006-01-02T15:04:05"`
	StartDate         *time.Time            `json:"startDate,omitempty" time_format:"2006-01-02T15:04:05"`
	FrontIdentify     string                `json:"frontIdentify,omitempty"`
	BackIdentify      string                `json:"backIdentify,omitempty"`
	Slogan            string                `json:"slogan,omitempty"`
	Website           string                `json:"website,omitempty"`
	Facebook          string                `json:"facebook,omitempty"`
	Instagram         string                `json:"instagram,omitempty"`
	Twitter           string                `json:"twitter,omitempty"`
	Linkedin          string                `json:"linkedin,omitempty"`
	Youtube           string                `json:"youtube,omitempty"`
	Introduction      string                `json:"introduction,omitempty"`
	Visibility        uint8                 `json:"visibility,omitempty"`
	ProfileVisibility uint8                 `json:"profileVisibility,omitempty"`
	StatusOnline      enums.EOnlineStatus   `json:"statusOnline,omitempty"`
	TabDefault        enums.ESettingTab     `json:"tabDefault,omitempty"`
	RoleRealEstate    enums.ERoleRealEstate `json:"roleRealEstate,omitempty"`
	ProvinceID        *uint64               `json:"provinceId,omitempty"`
	WardID            *uint64               `json:"wardId,omitempty"`
	Gender            uint32                `json:"gender,omitempty"`

	Certifications         []CertificationBatchSave `json:"certifications,omitempty"`
	CertificationRemoveIds []uint64                 `json:"certificationRemoveIds,omitempty"`
	Professions            []ProfessionBatchSave    `json:"professions,omitempty"`
	ProfessionRemoveIds    []uint64                 `json:"professionRemoveIds,omitempty"`
	PurposeUseIDs          []uint64                 `json:"purposeUseIds,omitempty"`
	MainAreaIds            []uint64                 `json:"mainAreaIds,omitempty"`
	Tags                   []models.TagEntity       `json:"tags,omitempty"`
	UpdateTags             bool                     `json:"updateTags,omitempty"`
}

type QuickSetupRequest struct {
	RoleRealEstate enums.ERoleRealEstate `json:"roleRealEstate,omitempty"`
	MainAreaIds    []uint64              `json:"mainAreaIds,omitempty"`
	PurposeUseIDs  []uint64              `json:"purposeUseIds,omitempty"`
	MainAreaNews   []_dto.ItemDTO        `json:"mainAreaNews,omitempty"`
}

type ProfessionBatchSave struct {
	Name      string     `json:"name,omitempty"`
	Issuer    string     `json:"issuer,omitempty"`
	IssueDate *time.Time `json:"issueDate,omitempty" time_format:"2006-01-02T15:04:05"`
	ID        uint64     `json:"id,omitempty"`
	FileName  string     `json:"fileName,omitempty"`
	FileURL   string     `json:"fileUrl,omitempty"`
	FileType  string     `json:"fileType,omitempty"`
	IsActive  bool       `json:"isActive,omitempty"`
}

type CertificationBatchSave struct {
	// Title     string    `json:"title,omitempty"`
	// Type      uint8     `json:"type,omitempty"`
	Name      string     `json:"name,omitempty"`
	Issuer    string     `json:"issuer,omitempty"`
	FileURL   string     `json:"fileUrl,omitempty"`
	FileType  string     `json:"fileType,omitempty"`
	FileName  string     `json:"fileName,omitempty"`
	IssueDate *time.Time `json:"issueDate,omitempty" time_format:"2006-01-02T15:04:05"`
	ID        uint64     `json:"id,omitempty"`
}

type TagRequest struct {
	ID      uint64        `json:"id,omitempty"`
	Name    string        `json:"name,omitempty"`
	TagType enums.TagType `json:"tagType,omitempty"`
}

// Hàm validate để kiểm tra dữ liệu đầu vào
func ValidateStruct(s interface{}) error {
	validate := validator.New()
	return validate.Struct(s)
}

type ProfileIDsRequest struct {
	IDs []uint64 `form:"ids"`
}

type SwitchOrganizationRequest struct {
	OrganizationID uint32 `json:"organizationId" binding:"required" validate:"required"`
}

// ReportCreateRequest đại diện cho yêu cầu tạo báo cáo
// @swagger:model
type ReportCreateRequest struct {
	ReasonID  uint64   `json:"reasonId" binding:"required" validate:"required"`
	UserID    uint64   `json:"userId" binding:"required" validate:"required"`
	ProofDocs []string `json:"proofDocs,omitempty"`
}

// ReportUpdateRequest đại diện cho yêu cầu cập nhật báo cáo
// @swagger:model
type ReportUpdateRequest struct {
	ReasonID  uint64   `json:"reasonId,omitempty"`
	ProofDocs []string `json:"proofDocs,omitempty"`
}

// ReportListRequest đại diện cho yêu cầu lấy danh sách báo cáo
// @swagger:model
type ReportListRequest struct {
	_dto.Pagable
	UserID   *uint64 `json:"userId,omitempty" form:"userId"`
	Status   *uint32 `json:"status,omitempty" form:"status"`
	ReasonID *uint64 `json:"reasonId,omitempty" form:"reasonId"`
}

// ReportResponseRequest đại diện cho yêu cầu phản hồi báo cáo
// @swagger:model
type ReportResponseRequest struct {
	Status           uint32 `json:"status" binding:"required" validate:"required"`
	Response         string `json:"response,omitempty"`
	AdminNote        string `json:"adminNote,omitempty"`        // Ghi chú nội bộ admin
	SendNotification bool   `json:"sendNotification,omitempty"` // Có gửi thông báo cho user không
}

// AdminNoteRequest đại diện cho yêu cầu cập nhật ghi chú admin
// @swagger:model
type AdminNoteRequest struct {
	ReportID   uint64 `json:"reportId" binding:"required" validate:"required"`
	AdminNote  string `json:"adminNote" binding:"required" validate:"required"`
	IsInternal bool   `json:"isInternal,omitempty"` // Ghi chú nội bộ hay công khai
}
