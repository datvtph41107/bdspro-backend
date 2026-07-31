package dto

import (
	_dto "common/domain/dto"
	"user/enums"
)

type ProfileSearch struct {
	_dto.Pagable
	// Page int64  `form:"page"`
	// Size int64  `form:"size"`
	Text string `form:"text"`
}

type AccountItemDTO struct {
	ID       uint64 `json:"id"`
	FullName string `json:"fullName"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
}

// Warning/Alert DTOs
type SendWarningRequest struct {
	UserID      uint64 `json:"userId" binding:"required"`
	WarningType int32  `json:"warningType" binding:"required"` // 0: vi phạm, 1: nhắc nhở, 2: hướng dẫn
	Content     string `json:"content" binding:"required"`
	Title       string `json:"title,omitempty"`
}

type SendWarningResponse struct {
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	NotificationID uint64 `json:"notificationId,omitempty"`
}

type ProfilePrivacyRequest struct {
	Visibility        uint8
	ProfileVisibility uint8
	StatusOnline      enums.EOnlineStatus
	PersonConfigs     []*PersonConfigDTO
}
