package dto

import (
	_dto "common/domain/dto"
	"time"
)

type ReportCreateRequest struct {
	ReasonID     uint64   `json:"reasonId" binding:"required" validate:"required"`
	UserID       uint64   `json:"userId" binding:"required" validate:"required"`
	OwnerID      uint64   `json:"ownerId" binding:"required" validate:"required"`
	OwnerOf      uint32   `json:"ownerOf" binding:"required" validate:"required"`
	RegistryName string   `json:"registryName,omitempty"`
	Content      string   `json:"content,omitempty"`
	ProofDocs    []string `json:"proofDocs,omitempty"`
}

type ReportUpdateRequest struct {
	ReasonID  uint64   `json:"reasonId,omitempty"`
	Content   string   `json:"content,omitempty"`
	ProofDocs []string `json:"proofDocs,omitempty"`
}

type ReportListRequest struct {
	_dto.Pagable
	OwnerId  *uint64 `json:"ownerId,omitempty" form:"ownerId"`
	OwnerOf  uint32  `json:"ownerOf,omitempty" form:"ownerOf"`
	Status   *uint32 `json:"status,omitempty" form:"status"`
	ReasonID *uint64 `json:"reasonId,omitempty" form:"reasonId"`
	Registry string  `json:"registry,omitempty" form:"registry"`
}

type ReportResponseRequest struct {
	Status           uint32 `json:"status" binding:"required" validate:"required"`
	Response         string `json:"response,omitempty"`
	AdminNote        string `json:"adminNote,omitempty"`        // Ghi chú nội bộ admin
	SendNotification bool   `json:"sendNotification,omitempty"` // Có gửi thông báo cho user không
}

type AdminNoteRequest struct {
	ReportID   uint64 `json:"reportId" binding:"required" validate:"required"`
	AdminNote  string `json:"adminNote" binding:"required" validate:"required"`
	IsInternal bool   `json:"isInternal,omitempty"` // Ghi chú nội bộ hay công khai
}

type AdminReportActionRequest struct {
	Response         *string `json:"response,omitempty"`
	AdminNote        *string `json:"adminNote,omitempty"`
	SendNotification bool    `json:"sendNotification,omitempty"`
}

type ReportResponse struct {
	ID           uint64        `json:"id"`
	ReasonID     uint64        `json:"reasonId"`
	Reason       *ReportReason `json:"reason,omitempty"`
	ReportStatus uint8         `json:"reportStatus"`
	UserID       uint64        `json:"userId"`
	ProofDocs    []ReportProof `json:"proofDocs,omitempty"`
	Response     string        `json:"response,omitempty"`
	AdminNote    string        `json:"adminNote,omitempty"`
	OwnerID      uint64        `json:"ownerId"`
	OwnerOf      uint32        `json:"ownerOf"`
	Content      string        `json:"content,omitempty"`
	CreatedAt    *time.Time    `json:"createdAt"`
	UpdatedAt    *time.Time    `json:"updatedAt"`
}

type ReportReason struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
}

type ReportProof struct {
	ID       uint64 `json:"id"`
	FileURL  string `json:"fileUrl"`
	FileName string `json:"fileName"`
	FileType string `json:"fileType"`
}

type UserProfileEntity struct {
	ProfileID uint64 `json:"profileId"`
	FullName  string `json:"fullName"`
	Avatar    string `json:"avatar"`
}