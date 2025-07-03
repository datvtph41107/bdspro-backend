package domain

import (
	"bdspro/internal/enums"
	_models "common/models"

	"github.com/lib/pq"
)

type RecordHistory struct {
	_models.BaseEntity
	RecordID   *uint64              `json:"recordId"`
	RecordType enums.ERecordType    `json:"recordType"`
	ActionType enums.EHistoryAction `json:"actionType"`
	Content    pq.StringArray       `json:"content"`
	// ParentID   *uint64              `json:"parentId"`
	// Date        string               `json:"date"`
	// Description string               `json:"description"`
	// PerformedBy uint64               `json:"performedBy"`
	// ChildIDs    pq.Int64Array        `json:"childIds"`
}

func (RecordHistory) TableName() string {
	return "record_history"
}
