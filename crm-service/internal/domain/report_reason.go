package domain

import (
	_models "common/models"
)

type ReportReason struct {
	_models.BaseEntity
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"isActive"`
}

func (ReportReason) TableName() string {
	return "feedback_report_reasons"
}

type ReportProof struct {
	_models.BaseEntity
	ReportID uint64 `gorm:"not null" json:"reportId"`
	FileURL  string `gorm:"type:varchar(500)" json:"fileUrl"`
	FileName string `gorm:"type:varchar(255)" json:"fileName"`
	FileType string `gorm:"type:varchar(50)" json:"fileType"`
}

func (ReportProof) TableName() string {
	return "feedback_report_proofs"
}

// ReportReasonListRequest represents request parameters for getting report reason list
type ReportReasonListRequest struct {
	IsActive *bool `json:"is_active,omitempty"`
	Page     int   `json:"page"`
	Size     int   `json:"size"`
}