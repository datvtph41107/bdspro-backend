package entity

import (
	"mdn4/modules/crm/enums"
	"time"
)

type LeadContactEntity struct {
	ID         uint64            `json:"id"`
	ContactID  *uint64           `json:"contactId"`
	Source     enums.ESourceLead `json:"customerSource"`
	FullName   string            `json:"fullName" binding:"required"`
	Phone      string            `json:"phone"`
	Email      string            `json:"email"`
	Address    string            `json:"address"`
	Birthday   *time.Time        `json:"birthday"`
	Note       string            `json:"note"`
	AssignNote string            `json:"assignNote"`
	PipelineID *uint64           `json:"pipelineId"`
	StageID    *uint64           `json:"stageId"`
	ChargeId   *uint64           `json:"chargeId"`
	ChargeType enums.EOwnerOf    `json:"chargeType"`
	OwnerID    *uint64           `json:"ownerId" binding:"required"`
	OwnerType  enums.EOwnerOf    `json:"ownerType" binding:"required"`

	// Step       enums.EStep           `gorm:"default:10" json:"step"`
	// todo: sắp trễ/ chưa phân công
}
