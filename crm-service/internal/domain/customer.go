package domain

import (
	_models "common/models"
	"crm/internal/enums"
	"time"
)

type CustomerEntity struct {
	_models.BaseEntity
	ContactID  *uint64               `json:"contactId"`
	FullName   string                `json:"fullName" binding:"required"`
	Phone      string                `json:"phone"`
	Source     enums.ESourceCustomer `json:"customerSource"`
	Email      string                `json:"email"`
	Address    string                `json:"address"`
	Birthday   *time.Time            `json:"birthday"`
	Note       string                `json:"note"`
	AssignNote string                `json:"assignNote"`
	PipelineID *uint64               `json:"pipelineId"`
	StageID    *uint64               `json:"stageId"`
	ChargeId   *uint64               `json:"chargeId"`
	ChargeType enums.EOwnerType      `json:"chargeType"`
	OwnerID    *uint64               `json:"ownerId" binding:"required"`
	OwnerType  enums.EOwnerType      `json:"ownerType" binding:"required"`
	// Step       enums.EStep           `gorm:"default:10" json:"step"`
	// todo: sắp trễ/ chưa phân công
}

func (e *CustomerEntity) TableName() string {
	return "customers"
}
