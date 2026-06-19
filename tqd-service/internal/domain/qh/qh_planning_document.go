package qh_domain

import (
	_models "common/domain/entity"
	"time"
)

type QHPlanningDocument struct {
	_models.BaseEntity
	PlanningProjectID uint64     `gorm:"not null;index" json:"planningProjectId"`
	DocumentType      string     `gorm:"type:varchar(50);not null" json:"documentType"`
	Code              string     `gorm:"type:varchar(255)" json:"code"`
	Title             string     `gorm:"type:varchar(255);not null" json:"title"`
	Description       string     `gorm:"type:text" json:"description"`
	Filepath          string     `gorm:"type:varchar(255)" json:"filepath"`
	VersionNo         string     `gorm:"type:varchar(50)" json:"versionNo"`
	ValidityStatus    string     `gorm:"type:varchar(50)" json:"validityStatus"`
	IssueDate         *time.Time `json:"issueDate"`
	EffectiveDate     *time.Time `json:"effectiveDate"`
	Metadata          string     `gorm:"type:jsonb" json:"metadata"`
}

func (QHPlanningDocument) TableName() string {
	return "qh_planning_documents"
}
