package provider

import (
	"context"

	_dto "common/domain/dto"
	notificationpb "pb/types/notification"
)

// ContactHistoryReader reads contact activity history owned by modules/notification (noti_histories).
type ContactHistoryReader interface {
	ContactHistory(ctx context.Context, contactID uint64, page _dto.Pagable) ([]*notificationpb.HistoryDTO, int64, error)
}
