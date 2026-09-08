package dto

import (
	base_enum "base/enum"
	"crm/internal/enums"
	"time"
)

type LeadSaveDTO struct {
	// contact
	FullName string     `gorm:"column:full_name" json:"fullName" binding:"required"`
	Phone    string     `gorm:"column:phone" json:"phone"`
	Avatar   string     `gorm:"column:avatar" json:"avatar"`
	Email    string     `gorm:"column:email" json:"email"`
	Zalo     string     `gorm:"column:zalo" json:"zalo"`
	Company  string     `gorm:"column:company" json:"company"`
	Address  string     `json:"address"`
	Birthday *time.Time `gorm:"column:birthday" json:"birthday"`
	Note     string     `gorm:"column:note" json:"note"`

	// lead
	Source            enums.ESourceLead `json:"source"`
	AssignNote        string            `json:"assignNote"`
	StageID           *uint64           `json:"stageId"`
	ChargePersonID    *uint64           `json:"chargePersonId"`
	Priority          enums.EPriority   `json:"priority"`
	ProductIDs        []uint64          `json:"productIds" parser:"uint64s"`
	Documents         []DocumentDTO     `json:"documents"`
	RemoveDocumentIDs []uint64          `json:"removeDocumentIds" parser:"uint64s"`

	OwnerID uint64             `json:"ownerId"`
	OwnerOf base_enum.EOwnerOf `json:"ownerOf"`
}

type DocumentDTO struct {
	ID       uint64 `json:"id"`
	FileName string `json:"fileName"`
	FileUrl  string `json:"fileUrl"`
	FileType string `json:"fileType"`
}