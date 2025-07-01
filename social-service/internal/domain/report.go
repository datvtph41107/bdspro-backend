package domain

import (
	_models "common/models"
	"social/internal/enums"
)

type Report struct {
	_models.BaseEntity
	Content    string             `json:"content"`
	ReasonID   uint64             `json:"reasonId"`
	UserID     uint64             `json:"userId"`
	Status     enums.ReportStatus `json:"status" gorm:"default:10"`
	TargetID   uint64             `json:"targetId"`
	TargetType enums.TargetType   `json:"targetType"`
	Reason     *ReportReason      `json:"reason" gorm:"foreignKey:ReasonID;references:ID"`
}

func (Report) TableName() string {
	return "report"
}
