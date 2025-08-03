package client

import (
	"context"
	"crm/env"
	delivery_dto "crm/internal/dto/delivery"
	"pb/clients"
	shared_dto "pb/dto"
	"time"
)

// @bind: crm/internal/interface.INotificationClient
type NotificationClient struct {
	Client *clients.NotificationClient
}

// todo: triển khai nốt thông báo
func NewNotificationClient() *NotificationClient {
	return &NotificationClient{
		Client: clients.MakeNotificationClient(env.ENV_RUNTIME),
	}
}

func (client *NotificationClient) SendNoti(c context.Context, noti *delivery_dto.NotiDTO) error {
	return nil
}

func (client *NotificationClient) SendReminder(c context.Context, customerID uint64) error {
	return nil
}

func (client *NotificationClient) CreateHistory(c context.Context, entity *shared_dto.HistoryDTO) error {
	if client.Client == nil {
		return nil
	}
	return client.Client.CreateHistory(c, entity)
}

func (client *NotificationClient) ExistedFromDate(c context.Context, customerID uint64, fromDate time.Time) (bool, error) {
	return false, nil
}
