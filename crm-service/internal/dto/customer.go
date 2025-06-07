package dto

import (
	_dto "common/dto"
	"crm/internal/enums"
	"time"
)

type CustomerSearchDTO struct {
	_dto.Page
	Phone          string                  `form:"phone"`
	FullName       string                  `form:"fullName"`
	Steps          []enums.EStep           `form:"steps" parser:"uint64s"`
	StageIDs       []uint64                `form:"stageIds" parser:"uint64s"`
	Sources        []enums.ESourceCustomer `form:"customerSources" parser:"uint64s"`
	ChargeIds      []uint64                `form:"chargeIds" parser:"uint64s"`
	ChargeType     enums.EOwnerType        `form:"chargeType"` // deprecated
	FromDate       *time.Time              `form:"fromDate"`
	ToDate         *time.Time              `form:"toDate"`
	ProductCareIDs []uint64                `form:"productCareIds" parser:"uint64s"` // todo: add product care
	// OwnerID    uint64                  `form:"ownerId"`
	// OwnerType  enums.EOwnerType        `form:"ownerType"`

}

type CustomerDTO struct {
	ContactID  *uint64               `json:"contactId"`
	FullName   string                `json:"fullName"`
	Phone      string                `json:"phone"`
	Source     enums.ESourceCustomer `json:"source"`
	Email      string                `json:"email"`
	Address    string                `json:"address"`
	Birthday   *time.Time            `json:"birthday"`
	Note       string                `json:"note"`
	AssignNote string                `json:"assignNote"`
	PipelineID *uint64               `json:"pipelineId"`
	StageID    *uint64               `json:"stageId"`
	ChargeId   *uint64               `json:"chargeId"`
	ChargeType enums.EOwnerType      `json:"chargeType"`
	OwnerID    *uint64               `json:"ownerId"`
	OwnerType  enums.EOwnerType      `json:"ownerType"`
}

type CustomerWithRuleDTO struct {
	CustomerID     uint64           `json:"customerId" db:"customer_id"`
	FullName       string           `json:"fullName" db:"full_name"`
	Phone          string           `json:"phone" db:"phone"`
	Email          string           `json:"email" db:"email"`
	Address        string           `json:"address" db:"address"`
	Note           string           `json:"note" db:"note"`
	Birthday       *time.Time       `json:"birthday" db:"birthday"`
	PipelineID     *uint64          `json:"pipelineId" db:"pipeline_id"`
	StageID        *uint64          `json:"stageId" db:"stage_id"`
	StageName      string           `json:"stageName" db:"stage_name"`
	ChargeId       uint64           `json:"chargeId" db:"charge_id"`
	RuleID         uint64           `json:"ruleId" db:"rule_id"`
	RuleName       string           `json:"ruleName" db:"rule_name"`
	Condition      enums.ECondition `json:"condition" db:"condition"`
	ConditionValue string           `json:"conditionValue" db:"condition_value"`
	Trigger        enums.ERule      `json:"trigger" db:"trigger"`
	TriggerValue   string           `json:"triggerValue" db:"trigger_value"`
	OwnerID        uint64           `json:"ownerId" db:"owner_id"`
	OwnerType      enums.EOwnerType `json:"ownerType" db:"owner_type"`
	UpdatedAt      time.Time        `json:"updatedAt" db:"updated_at"`
	TriggerToDay   bool             `json:"triggerToDay" db:"trigger_to_day"`
}
