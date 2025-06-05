package domain

import (
	_models "common/models"
	"crm/internal/enums"
	"time"
)

type RuleEntity struct {
	_models.BaseEntity
	OwnerID      uint64           `json:"ownerId"`
	OwnerType    enums.EOwnerType `json:"ownerType"`
	Rule         enums.ERule      `json:"rule"`
	CycleSeconds int64            `json:"cycleSeconds"`
	Deadline     *time.Time       `json:"deadline"`
	Active       bool             `json:"active"`
	Note         string           `json:"note"`
}
