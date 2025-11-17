package models

import (
	_models "common/domain/entity"
	_enum "common/domain/enum"
	"time"
)

type KYCEntity struct {
	_models.BaseEntity
	ProfileID    uint64               `gorm:"column:profile_id;not null;index" json:"profileId"`
	FullName     string               `gorm:"column:full_name;size:255;not null" json:"fullName"`
	IdentityCard string               `gorm:"column:identity_card;size:20;not null" json:"identityCard"`
	FrontImage   string               `gorm:"column:front_image;size:500;not null" json:"frontImage"`
	BackImage    string               `gorm:"column:back_image;size:500;not null" json:"backImage"`
	SelfieImage  string               `gorm:"column:selfie_image;size:500" json:"selfieImage"`
	Status       _enum.EApproveStatus `gorm:"column:status;default:10" json:"status"`
	RejectReason string               `gorm:"column:reject_reason;size:500" json:"rejectReason"`
	ReviewedBy   *uint64              `gorm:"column:reviewed_by" json:"reviewedBy"`
	ReviewedAt   *time.Time           `gorm:"column:reviewed_at" json:"reviewedAt"`
}

// TableName đặt tên bảng trong DB
func (KYCEntity) TableName() string {
	return "kyc"
}
