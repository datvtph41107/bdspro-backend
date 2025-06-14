package iusecase

import (
	"context"
	delivery_dto "crm/internal/dto/delivery"
)

type NotiClient interface {
	SendReminder(c context.Context, customerID uint64) error
	SendNoti(c context.Context, noti *delivery_dto.NotiDTO) error
}
