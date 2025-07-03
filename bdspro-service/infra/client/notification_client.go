package client

import (
	"bdspro/env"
	delivery_dto "bdspro/internal/dto/delivery"
	"context"
	"pb/clients"
	shared_dto "pb/dto"
)

// @bind: bdspro/internal/usecases/interface.INotificationClient
type NotificationClient struct {
	Client *clients.NotificationClient
}

// todo: triển khai nốt thông báo
func NewNotificationClient() *NotificationClient {
	return &NotificationClient{
		Client: clients.MakeNotificationClient(env.ENV_RELEASE),
	}
}

func (s *NotificationClient) UserInGroup(ctx context.Context, userID uint64, groupID uint64) (bool, error) {
	return true, nil
}

func (s *NotificationClient) SendNoti(ctx context.Context, noti *delivery_dto.NotiDTO) error {
	return nil
}

func (s *NotificationClient) SendBatch(ctx context.Context, noti *delivery_dto.NotiBatchDTO) error {
	return nil
}

func (s *NotificationClient) CreateHistory(ctx context.Context, history *shared_dto.HistoryDTO) error {
	err := s.Client.CreateHistory(ctx, history)
	return err
}
