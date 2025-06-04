package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
)

type RecordHistory struct {
	_models.BaseEntity
	RecordID    uint64               `json:"recordId"`
	RecordType  enums.ERecordType    `json:"recordType"`
	ActionType  enums.EHistoryAction `json:"actionType"`
	ParentID    *uint64              `json:"parentId"`
	ChildIDs    string               `json:"childIds"`
	Date        string               `json:"date"`
	Description string               `json:"description"`
	PerformedBy uint64               `json:"performedBy"`
}
