package client

import (
	"bdspro/env"
	_dto "common/domain/dto"
	"context"
	"pb/clients"
)

// @bind: bdspro/internal/provider.NotificationProvider
type NotificationClient struct {
	*clients.NotificationClient
}

// todo: triển khai nốt thông báo
func NewNotificationClient() *NotificationClient {
	return &NotificationClient{
		NotificationClient: clients.MakeNotificationClient(env.ENV_RUNTIME),
	}
}

func (s *NotificationClient) UserInGroup(ctx context.Context, userID uint64, groupID uint64) (bool, error) {
	return true, nil
}

func (s *NotificationClient) SendNoti(ctx context.Context, noti *_dto.NotificationDTO) error {
	return nil
}

func (s *NotificationClient) SendBatch(ctx context.Context, noti *_dto.NotificationBatchDTO) error {
	return nil
}

func (s *NotificationClient) CreateHistory(ctx context.Context, history *_dto.HistoryDTO) error {
	err := s.NotificationClient.CreateHistory(ctx, history)
	return err
}
