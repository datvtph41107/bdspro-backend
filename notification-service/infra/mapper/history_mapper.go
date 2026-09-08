package mapper

import (
	_utils "common/utils"
	"notification/internal/domain"
	notipb "pb/types/notification"
)

type HistoryMapper struct {
}

func NewHistoryMapper() *HistoryMapper {
	return &HistoryMapper{}
}

func (h *HistoryMapper) HistoryToPbs(histories []domain.HistoryEntity) []*notipb.HistoryDTO {
	result := make([]*notipb.HistoryDTO, len(histories))
	for i, history := range histories {
		result[i] = h.HistoryDomainToPb(&history)
	}
	return result
}

func (h *HistoryMapper) HistoryDomainToPb(history *domain.HistoryEntity) *notipb.HistoryDTO {
	return &notipb.HistoryDTO{
		TargetId:   history.TargetId,
		TargetType: int32(history.TargetType),
		ActionType: int32(history.ActionType),
		Note:       history.Note,
		PreStage:   history.PreStage,
		AfterStage: history.AfterStage,
		OwnerId:    history.OwnerID,
		OwnerType:  int32(history.OwnerType),
		CreatedAt:  _utils.FormatTimeToString(history.CreatedAt),
		ActionUser: "Hệ thống",
		ActionName: "Thao tác",
		// CreatedBy:  history.CreatedBy,
		Title: history.Title,
	}
}
