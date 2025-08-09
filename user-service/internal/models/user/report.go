package user_model

import (
	_models "common/models"
	"user/enums"
)

type Report struct {
	_models.BaseEntity
	ReasonID     uint64              `gorm:"not null"`
	Reason       *ReportReason       `gorm:"foreignKey:ReasonID"`
	ReportStatus enums.ApproveStatus `gorm:"not null;default:10"`
	UserID       uint64              `gorm:"not null"`
	User         *UserProfileEntity  `gorm:"foreignKey:UserID;references:ProfileID"`
	ProofDocs    []ReportProof       `gorm:"foreignKey:ReportID;references:ID"`
	Response     string              `gorm:"type:varchar(255)"`
	Reporter     *UserProfileEntity  `gorm:"foreignKey:CreatedBy;references:ProfileID"`
	AdminNote    string              `gorm:"type:varchar(255)"`
}
