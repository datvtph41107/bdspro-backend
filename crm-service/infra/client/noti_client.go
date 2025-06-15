package client

import (
	"context"
	delivery_dto "crm/internal/dto/delivery"
)

// @bind: crm/internal/interface.NotiClient
type NotiClient struct {
}

func NewNotiClient() *NotiClient {
	return &NotiClient{}
}

func (client *NotiClient) SendNoti(c context.Context, noti *delivery_dto.NotiDTO) error {
	return nil
}

func (client *NotiClient) SendReminder(c context.Context, customerID uint64) error {
	return nil
}
