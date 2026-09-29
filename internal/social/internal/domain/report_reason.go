package domain

import (
	_models "common/models"
	"social/internal/enums"
)

type ReportReason struct {
	_models.BaseEntity
	ReasonName string           `json:"reasonName"`
	Active     bool             `gorm:"column:active;default:true" json:"active"`
	TargetType enums.TargetType `json:"targetType"`
}

func (ReportReason) TableName() string {
	return "report_reason"
}
