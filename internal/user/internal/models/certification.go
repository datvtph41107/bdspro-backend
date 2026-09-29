package models

import (
	_models "common/models"
	"time"
	"user/enums"
)

type CertificationEntity struct {
	_models.BaseEntity
	ApproveUserID  *uint64            `json:"approveUserId"`
	Name           string             `gorm:"type:varchar(255)" json:"name"`                  // Tên chứng chỉ
	FileName       string             `gorm:"type:varchar(255)" json:"fileName"`              // Tên file
	FileURL        string             `gorm:"type:text" json:"fileUrl"`                       // Link tài liệu
	FileType       string             `gorm:"type:varchar(255)" json:"fileType"`              // Loại file
	VerifiedStatus enums.VerifyStatus `gorm:"type:smallint;default:10" json:"verifiedStatus"` // Trạng thái xác minh
	Issuer         string             `gorm:"type:varchar(255)" json:"issuer"`                // Cơ quan cấp
	IssueDate      *time.Time         `gorm:"type:date" json:"issueDate"`                     // Ngày cấp
	// Title          string             `gorm:"type:varchar(255);not null" json:"title"`        // Tên văn bằng
}

func (CertificationEntity) TableName() string {
	return "certification"
}

// type CertificationItem struct {
// 	ID             uint64             `json:"id"`
// 	Name           string             `gorm:"type:varchar(255)" json:"name"`                  // Tên chứng chỉ
// 	FileName       string             `gorm:"type:varchar(255)" json:"fileName"`              // Tên file
// 	FileURL        string             `gorm:"type:text" json:"fileUrl"`                       // Link tài liệu
// 	FileType       string             `gorm:"type:varchar(255)" json:"fileType"`              // Loại file
// 	VerifiedStatus enums.VerifyStatus `gorm:"type:smallint;default:10" json:"verifiedStatus"` // Trạng thái xác minh
// }

// func (CertificationItem) TableName() string {
// 	return "certification"
// }
