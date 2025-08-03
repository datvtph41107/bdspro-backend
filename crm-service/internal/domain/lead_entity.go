package domain

import (
	_models "common/models"
	"crm/internal/enums"
	"time"
)

type LeadEntity struct {
	_models.BaseEntity
	ContactID        uint64            `json:"contactId"`
	Source           enums.ESourceLead `json:"customerSource"`
	AssignNote       string            `json:"assignNote"`
	PipelineID       *uint64           `json:"pipelineId"`
	StageID          *uint64           `json:"stageId"`
	Priority         enums.EPriority   `json:"priority"`
	ChargePersonID   *uint64           `json:"chargePersonId"`
	ChargePersonType enums.EOwnerType  `json:"chargePersonType"`
	OwnerID          *uint64           `json:"ownerId" binding:"required"`
	OwnerType        enums.EOwnerType  `json:"ownerType" binding:"required"`
	StageNote        string            `json:"stageNote"`
	Note             string            `json:"note"`

	Stage        *StageEntity        `gorm:"references:ID" json:"stage"`
	Pipeline     *PipelineEntity     `gorm:"-" json:"pipeline"`
	Documents    []DocumentEntity    `gorm:"foreignKey:LeadID" json:"documents"`
	ProductCares []ProductCareEntity `gorm:"foreignKey:LeadID" json:"productCares"`
	Contact      *ContactEntity      `gorm:"foreignKey:ContactID;references:ID" json:"contact"`
	// FullName   string                `json:"fullName" binding:"required"`
	// Phone      string                `json:"phone"`
	// Email      string                `json:"email"`
	// Address    string                `json:"address"`
	// Birthday   *time.Time            `json:"birthday"`

	// Step       enums.EStep           `gorm:"default:10" json:"step"`
	// todo: sắp trễ/ chưa phân công
}

func (e *LeadEntity) TableName() string {
	return "customers"
}

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
	ChargeType enums.EOwnerType  `json:"chargeType"`
	OwnerID    *uint64           `json:"ownerId" binding:"required"`
	OwnerType  enums.EOwnerType  `json:"ownerType" binding:"required"`

	// Step       enums.EStep           `gorm:"default:10" json:"step"`
	// todo: sắp trễ/ chưa phân công
}
