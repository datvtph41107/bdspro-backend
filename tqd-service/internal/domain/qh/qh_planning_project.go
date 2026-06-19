package qh_domain

import (
	_models "common/domain/entity"
	"time"
)

type QHPlanningProject struct {
	_models.BaseEntity
	Code            string     `gorm:"type:varchar(255);not null;index" json:"code"`
	Name            string     `gorm:"type:varchar(255);not null" json:"name"`
	PlanningType    string     `gorm:"type:varchar(50);not null" json:"planningType"`
	PlanningLevel   string     `gorm:"type:varchar(50);not null" json:"planningLevel"`
	Authority       string     `gorm:"type:varchar(255)" json:"authority"`
	JurisdictionID  *string    `gorm:"type:uuid" json:"jurisdictionId"`
	Summary         string     `gorm:"type:text" json:"summary"`
	ApprovalDate    *time.Time `json:"approvalDate"`
	EffectiveDate   *time.Time `json:"effectiveDate"`
	ExpiryDate      *time.Time `json:"expiryDate"`
	ValidityStatus  string     `gorm:"type:varchar(50);not null" json:"validityStatus"`
	CurrentVersion  string     `gorm:"type:varchar(50)" json:"currentVersion"`
	Metadata        string     `gorm:"type:jsonb" json:"metadata"`
}

func (QHPlanningProject) TableName() string {
	return "qh_planning_projects"
}
