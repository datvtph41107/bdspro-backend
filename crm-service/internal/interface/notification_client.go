package iusecase

import (
	"context"
	delivery_dto "crm/internal/dto/delivery"
	shared_dto "pb/dto"
	"time"
)

type INotificationClient interface {
	SendReminder(c context.Context, customerID uint64) error
	SendNoti(c context.Context, noti *delivery_dto.NotiDTO) error
	CreateHistory(c context.Context, entity *shared_dto.HistoryDTO) error
	ExistedFromDate(c context.Context, customerID uint64, fromDate time.Time) (bool, error)
}
