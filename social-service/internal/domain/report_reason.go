package domain

import _models "common/models"

type ReportReason struct {
	_models.BaseEntity
	ReasonName string `json:"reasonName"`
	Active     bool   `gorm:"column:active;default:true" json:"active"`
}

func (ReportReason) TableName() string {
	return "report_reason"
}
