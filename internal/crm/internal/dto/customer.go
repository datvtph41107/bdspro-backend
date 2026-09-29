package dto

import (
	base_enum "base/enum"
	_dto "common/domain/dto"
	"crm/internal/enums"
	"time"
)

type LeadSearchDTO struct {
	_dto.Pagable
	Phone           string              `form:"phone"`
	FullName        string              `form:"fullName"`
	Steps           []enums.EStep       `form:"steps" parser:"uint64s"`
	StageIDs        []uint64            `form:"stageIds" parser:"uint64s"`
	Sources         []enums.ESourceLead `form:"customerSources" parser:"uint64s"`
	ChargePersonIds []uint64            `form:"chargePersonIds" parser:"uint64s"`
	ChargeType      enums.EOwnerOf      `form:"chargeType"` // deprecated
	FromDate        *time.Time          `form:"fromDate"`
	ToDate          *time.Time          `form:"toDate"`
	LastUpdatedAt   *time.Time          `form:"lastUpdatedAt"`
	ProductCareIDs  []uint64            `form:"productCareIds" parser:"uint64s"`
	PipelineID      *uint64             `form:"pipelineId"`
	UnassignedOnly  bool                `form:"unassignedOnly"`

	// Admin sales filters (IV.10.9.1)
	OpportunityStatuses []enums.EOpportunityStatus `form:"opportunityStatuses"`
	CustomerTypes       []enums.ECustomerType      `form:"customerTypes"`
	AdminSources        []enums.EAdminSalesSource  `form:"adminSources"`
	Queue               string                     `form:"queue"` // overdue_followup|unassigned|waiting_quote|waiting_approval|waiting_payment|renewal|upgrade|churn|new|consulting|business|api|mine (mine via chargePersonIds on FE)
	FollowUpFrom        *time.Time                 `form:"followUpFrom"`
	FollowUpTo          *time.Time                 `form:"followUpTo"`
	MinExpectedValue    *float64                   `form:"minExpectedValue"`
	MaxExpectedValue    *float64                   `form:"maxExpectedValue"`
	Probabilities       []enums.EWinProbability    `form:"probabilities"`
	ChurnRisks          []enums.EChurnRisk         `form:"churnRisks"`
	UpgradeOnly         bool                       `form:"upgradeOnly"`
	RenewalOnly         bool                       `form:"renewalOnly"`
	InterestedPlan      string                     `form:"interestedPlan"`

	OwnerID uint64             `form:"ownerId"`
	OwnerOf base_enum.EOwnerOf `form:"ownerOf"`
}

type LeadDTO struct {
	ContactID      *uint64           `json:"contactId"`
	FullName       string            `json:"fullName"`
	Phone          string            `json:"phone"`
	Source         enums.ESourceLead `json:"source"`
	Email          string            `json:"email"`
	Address        string            `json:"address"`
	Birthday       *time.Time        `json:"birthday"`
	Note           string            `json:"note"`
	AssignNote     string            `json:"assignNote"`
	PipelineID     *uint64           `json:"pipelineId"`
	StageID        *uint64           `json:"stageId"`
	ChargePersonId *uint64           `json:"chargePersonId"`
	ChargeType     enums.EOwnerOf    `json:"chargeType"`
	OwnerID        *uint64           `json:"ownerId"`
	OwnerType      enums.EOwnerOf    `json:"ownerType"`
}

type CustomerWithRuleDTO struct {
	// Phone          string             `json:"phone" db:"phone"`
	CustomerID     uint64             `json:"customerId" db:"customer_id"`
	FullName       string             `json:"fullName" db:"full_name"`
	Email          string             `json:"email" db:"email"`
	Address        string             `json:"address" db:"address"`
	Note           string             `json:"note" db:"note"`
	Birthday       *time.Time         `json:"birthday" db:"birthday"`
	PipelineID     *uint64            `json:"pipelineId" db:"pipeline_id"`
	StageID        *uint64            `json:"stageId" db:"stage_id"`
	StageName      string             `json:"stageName" db:"stage_name"`
	ChargePersonId uint64             `json:"chargePersonId" db:"charge_person_id"`
	RuleID         uint64             `json:"ruleId" db:"rule_id"`
	RuleName       string             `json:"ruleName" db:"rule_name"`
	Condition      enums.ECondition   `json:"condition" db:"condition"`
	ConditionValue string             `json:"conditionValue" db:"condition_value"`
	Trigger        enums.ERuleTrigger `json:"trigger" db:"trigger"`
	TriggerValue   string             `json:"triggerValue" db:"trigger_value"`
	OwnerID        uint64             `json:"ownerId" db:"owner_id"`
	OwnerType      enums.EOwnerOf     `json:"ownerType" db:"owner_type"`
	UpdatedAt      time.Time          `json:"updatedAt" db:"updated_at"`
	TriggerToDay   bool               `json:"triggerToDay" db:"trigger_to_day"`
}