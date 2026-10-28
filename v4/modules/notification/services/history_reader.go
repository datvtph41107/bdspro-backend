package services

import (
	"context"

	_dto "common/domain/dto"
	_utils "common/utils"
	"mdn4/modules/notification/entity"
	"mdn4/modules/notification/repository"
	shared_enum "pb/enums"
	notificationpb "pb/types/notification"
)

// LocalHistoryReader serves activity history (noti_histories) to in-process consumers
// (parity notification-service HistoryUsecase.ContactHistory + HistoryMapper).
// @bind: mdn4/modules/crm/provider.ContactHistoryReader
type LocalHistoryReader struct {
	repo repository.NotificationRepository
}

func NewLocalHistoryReader(repo repository.NotificationRepository) *LocalHistoryReader {
	return &LocalHistoryReader{repo: repo}
}

// ContactHistory lists CRM activity of one contact (action types in HistoryContactAction).
func (r *LocalHistoryReader) ContactHistory(ctx context.Context, contactID uint64, page _dto.Pagable) ([]*notificationpb.HistoryDTO, int64, error) {
	rows, total, err := r.repo.ListActivityHistory(ctx, contactID, shared_enum.HistoryContactAction, page)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*notificationpb.HistoryDTO, 0, len(rows))
	for _, row := range rows {
		if row != nil {
			out = append(out, activityHistoryToPb(row))
		}
	}
	return out, total, nil
}

func activityHistoryToPb(h *entity.ActivityHistory) *notificationpb.HistoryDTO {
	return &notificationpb.HistoryDTO{
		TargetId:   h.TargetID,
		TargetType: int32(h.TargetType),
		ActionType: int32(h.ActionType),
		Note:       h.Note,
		PreStage:   h.PreStage,
		AfterStage: h.AfterStage,
		OwnerId:    h.OwnerID,
		OwnerType:  int32(h.OwnerType),
		CreatedAt:  _utils.FormatTimeToString(h.CreatedAt),
		ActionUser: "Hệ thống",
		ActionName: "Thao tác",
		Title:      h.Title,
	}
}
