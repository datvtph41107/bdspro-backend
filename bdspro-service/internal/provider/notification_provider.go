package provider

import (
	"bdspro/internal/dto"
	_dto "common/domain/dto"
	"context"
)

type NotificationProvider interface {
	// _provider.NotificationProvider
	UserInGroup(ctx context.Context, userID uint64, groupID uint64) (bool, error)
	SendNoti(ctx context.Context, noti *_dto.NotificationDTO) error
	SendBatch(ctx context.Context, noti *_dto.NotificationBatchDTO) error
	CreateHistory(ctx context.Context, history *_dto.HistoryDTO) error
	HistoryDeal(ctx context.Context, history *dto.DealHistoryDTO) error
	RegistedEventProperty(ctx context.Context, arg *dto.CreatePropertyActivityDTO) error
}
