package mapper

import (
	"crm/internal/domain"
	crmpb "pb/types/crm"
	"time"
)

func HistoryToPbs(histories []domain.HistoryEntity) []*crmpb.HistoryDTO {
	result := make([]*crmpb.HistoryDTO, len(histories))
	for i, history := range histories {
		result[i] = HistoryDomainToPb(&history)
	}
	return result
}

// --map--

func HistoryDomainToPb(history *domain.HistoryEntity) *crmpb.HistoryDTO {
	return &crmpb.HistoryDTO{
		TargetId:   history.TargetId,
		TargetType: int32(history.TargetType),
		ActionType: int32(history.ActionType),
		Note:       history.Note,
		PreStage:   history.PreStage,
		AfterStage: history.AfterStage,
		OwnerId:    history.OwnerID,
		OwnerType:  int32(history.OwnerType),
		CreatedAt:  history.CreatedAt.Format(time.RFC3339),
		// CreatedBy:  history.CreatedBy,
		Title: history.Title,
	}
}
