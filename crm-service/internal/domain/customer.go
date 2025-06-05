package domain

import (
	_models "common/models"
	"crm/internal/enums"
)

type CustomerEntity struct {
	_models.BaseEntity
	ContactID  *uint64               `json:"contactId"`
	FullName   string                `json:"fullName" binding:"required"`
	Phone      string                `json:"phone"`
	Source     enums.ESourceCustomer `json:"customerSource"`
	Email      string                `json:"email"`
	Address    string                `json:"address"`
	Step       enums.EStep           `gorm:"default:10" json:"step"`
	Note       string                `json:"note"`
	PipelineID *uint64               `json:"pipelineId"`
	StageID    *uint64               `json:"stageId"`
	ChargeId   uint64                `json:"chargeId"`
	ChargeType enums.EOwnerType      `json:"chargeType"`
	OwnerID    uint64                `json:"ownerId"`
	OwnerType  enums.EOwnerType      `json:"ownerType" binding:"required"`
	// todo: sắp trễ/ chưa phân công
}
