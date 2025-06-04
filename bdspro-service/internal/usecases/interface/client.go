package iusecase

import (
	delivery_dto "bdspro/internal/dto/delivery"
	"context"
)

type IClientUC interface {
	UserInGroup(ctx context.Context, userID uint64, groupID uint64) (bool, error)
	SendNoti(ctx context.Context, noti *delivery_dto.NotiDTO) error
	SendBatch(ctx context.Context, noti *delivery_dto.NotiBatchDTO) error
}
