package domain

import _models "common/models"

// AdminOpportunityEvent nhật ký thao tác cơ hội admin (owner_of=ADMIN).
// opportunity_id = customers.id (LeadEntity).
type AdminOpportunityEvent struct {
	_models.BaseEntity
	OpportunityID uint64 `gorm:"column:opportunity_id;not null;index" json:"opportunityId"`
	ActorID       uint64 `gorm:"column:actor_id;not null" json:"actorId"`
	Action        string `gorm:"column:action;type:varchar(64);not null" json:"action"`
	BeforeJSON    string `gorm:"column:before_json;type:text" json:"beforeJson"`
	AfterJSON     string `gorm:"column:after_json;type:text" json:"afterJson"`
	Note          string `gorm:"column:note;type:text" json:"note"`
}

func (AdminOpportunityEvent) TableName() string { return "admin_opportunity_events" }
