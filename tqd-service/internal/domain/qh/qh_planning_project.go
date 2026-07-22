package qh_domain

import (
	_models "common/domain/entity"
	"encoding/json"
	"time"
	"tqd/internal/enums"
)

type QHPlanningProject struct {
	_models.BaseEntity
	Code string `gorm:"type:varchar(255);not null;index" json:"code"`
	Slug string `gorm:"type:varchar(180);index" json:"slug"`
	Name string `gorm:"type:varchar(255);not null" json:"name"`

	PlanningType  uint32 `gorm:"type:int4;not null" json:"planningType"`
	PlanningLevel uint32 `gorm:"type:int4;not null" json:"planningLevel"`

	TotalArea float64 `gorm:"type:float8" json:"totalArea"`

	JurisdictionID *uint64         `gorm:"type:int8" json:"jurisdictionId"`
	Jurisdiction   *QHJurisdiction `gorm:"foreignKey:JurisdictionID" json:"jurisdiction"`

	LegalStatus enums.LegalStatus `gorm:"type:int4;default:800" json:"legalStatus"`

	Authority string `gorm:"type:varchar(255)" json:"authority"`

	ResearchScope string          `gorm:"type:text" json:"researchScope"`
	Indicators    json.RawMessage `gorm:"type:jsonb;default:'[]'" json:"indicators"`

	Summary        string     `gorm:"type:text" json:"summary"`
	ApprovalDate   *time.Time `json:"approvalDate"`
	EffectiveDate  *time.Time `json:"effectiveDate"`
	ExpiryDate     *time.Time `json:"expiryDate"`
	ValidityStatus string     `gorm:"type:varchar(50);not null" json:"validityStatus"`
	CurrentVersion string     `gorm:"type:varchar(50)" json:"currentVersion"`
	Metadata       string     `gorm:"type:jsonb" json:"metadata"`

	Layers []*QHLayer `gorm:"foreignKey:PlanningProjectID" json:"layers"`
}

func (QHPlanningProject) TableName() string {
	return "qh_planning_projects"
}
