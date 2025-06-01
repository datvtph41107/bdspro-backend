package user_model

import (
	_models "common/models"
	"time"
	"user/enums"
)

type CertificationEntity struct {
	_models.BaseEntity
	ApproveUserID  *uint64            `json:"approveUserId"`
	Type           uint8              `gorm:"type:smallint" json:"type"`                      // Loại chứng chỉ
	Title          string             `gorm:"type:varchar(255);not null" json:"title"`        // Tên văn bằng
	Issuer         string             `gorm:"type:varchar(255)" json:"issuer"`                // Cơ quan cấp
	IssueDate      time.Time          `gorm:"type:date" json:"issueDate"`                     // Ngày cấp
	FileURL        string             `gorm:"type:text" json:"fileUrl"`                       // Link tài liệu
	VerifiedStatus enums.VerifyStatus `gorm:"type:smallint;default:10" json:"verifiedStatus"` // Trạng thái xác minh
}

func (CertificationEntity) TableName() string {
	return "certification"
}

type CertificationItem struct {
	ID             uint64             `json:"id"`
	Title          string             `gorm:"type:varchar(255);not null" json:"title"`        // Tên văn bằng
	FileURL        string             `gorm:"type:text" json:"fileUrl"`                       // Link tài liệu
	VerifiedStatus enums.VerifyStatus `gorm:"type:smallint;default:10" json:"verifiedStatus"` // Trạng thái xác minh
}

func (CertificationItem) TableName() string {
	return "certification"
}
