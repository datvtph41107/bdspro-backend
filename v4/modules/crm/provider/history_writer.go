package provider

import (
	"context"

	_dto "common/domain/dto"
)

// HistoryWriter records CRM activity into the history owned by modules/notification (noti_histories).
type HistoryWriter interface {
	CreateHistory(ctx context.Context, history *_dto.HistoryDTO) error
}
