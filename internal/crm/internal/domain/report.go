package domain

import (
	_entity "common/domain/entity"
)

type Report struct {
	_entity.BaseEntity
	ReasonID     uint64        `gorm:"not null"`
	Reason       *ReportReason `gorm:"foreignKey:ReasonID"`
	ReportStatus uint8         `gorm:"not null;default:10"` // 10: pending, 20: approved, 30: rejected
	UserID       uint64        `gorm:"not null"`
	ProofDocs    []ReportProof `gorm:"foreignKey:ReportID;references:ID"`
	Response     string        `gorm:"type:varchar(255)"`
	AdminNote    string        `gorm:"type:varchar(255)"`
	OwnerID      uint64        `gorm:"not null"`
	OwnerOf      uint32        `gorm:"not null"` // 10: user, 20: post, 30: comment
	Content      string        `gorm:"type:text"`
}

func (Report) TableName() string {
	return "feedback_reports"
}