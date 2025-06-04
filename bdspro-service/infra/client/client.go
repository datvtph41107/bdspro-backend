package client

import (
	delivery_dto "bdspro/internal/dto/delivery"
	"context"
)

// @bind: bdspro/internal/usecases/interface.IClientUC
type ServiceClient struct{}

func NewServiceClient() *ServiceClient {
	return &ServiceClient{}
}

func (s *ServiceClient) UserInGroup(ctx context.Context, userID uint64, groupID uint64) (bool, error) {
	return true, nil
}

func (s *ServiceClient) SendNoti(ctx context.Context, noti *delivery_dto.NotiDTO) error {
	return nil
}

func (s *ServiceClient) SendBatch(ctx context.Context, noti *delivery_dto.NotiBatchDTO) error {
	return nil
}
