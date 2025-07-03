package mapper

import (
	"bdspro/internal/domain"
	pb_bdspro "pb/types/bdspro"
)

type HistoryMapper struct {
}

func NewHistoryMapper() *HistoryMapper {
	return &HistoryMapper{}
}

func (m *HistoryMapper) HistoryToPb(history *domain.RecordHistory) *pb_bdspro.History {
	return &pb_bdspro.History{
		Id:      history.ID,
		Action:  uint64(history.ActionType),
		Content: history.Content,
		// : history.ParentID,
		// ParentId:    history.ParentID,
		// ChildIds:    history.ChildIDs,
		// Date:        history.Date,
		// Description: history.Description,
		// PerformedBy: history.PerformedBy,
	}
}
