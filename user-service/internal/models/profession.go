package models

import (
	_models "common/models"
	"time"
	"user/enums"
)

type ProfessionEntity struct {
	_models.BaseEntity
	ApproveUserID  *uint64            `json:"approveUserId"`
	Name           string             `gorm:"type:varchar(255);not null" json:"name"`          // tên văn bằng
	Issuer         string             `gorm:"type:varchar(255)" json:"issuer"`                 // cơ quan cấp
	IssueDate      *time.Time         `gorm:"type:date" json:"issueDate"`                      // ngày cấp
	VerifiedStatus enums.VerifyStatus `gorm:"type:smallint;not null;default:10" json:"status"` // trạng thái
	IsActive       bool               `gorm:"type:boolean;not null;default:true" json:"isActive"`
}

func (ProfessionEntity) TableName() string {
	return "profession"
}

// type ProfessionItem struct {
// 	ID             uint64             `json:"id"`
// 	Name           string             `json:"name"`
// 	VerifiedStatus enums.VerifyStatus `json:"verifiedStatus"`
// }

// func (ProfessionItem) TableName() string {
// 	return "profession"
// }
