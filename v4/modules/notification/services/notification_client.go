package services

import (
	"context"
	"errors"

	notificationpb "pb/types/notification"
	"userv4/infra/client"
	"userv4/modules/notification/dto"
)

// ClientNotificationService adapts remote GatewayNotificationService into modules/notification.
type ClientNotificationService struct {
	client *client.NotificationClient
}

func NewClientNotificationService(c *client.NotificationClient) *ClientNotificationService {
	return &ClientNotificationService{client: c}
}

func (s *ClientNotificationService) gateway() (notificationpb.GatewayNotificationServiceClient, error) {
	if s.client == nil || s.client.NotificationClient == nil || s.client.GatewayClient == nil {
		return nil, errors.New("notification gateway client not available")
	}
	return s.client.GatewayClient, nil
}

func (s *ClientNotificationService) Count(ctx context.Context, req dto.CountRequest) (*notificationpb.CountResponse, error) {
	gw, err := s.gateway()
	if err != nil {
		return nil, err
	}
	return gw.Count(ctx, &notificationpb.CountRequest{Timestamp: req.Timestamp})
}

func (s *ClientNotificationService) List(ctx context.Context, req dto.ListRequest) (*notificationpb.ListResponse, error) {
	gw, err := s.gateway()
	if err != nil {
		return nil, err
	}
	return gw.List(ctx, &notificationpb.ListRequest{
		Page:      req.Page,
		Size:      req.Size,
		Sort:      req.Sort,
		IsRead:    req.IsRead,
		Timestamp: req.Timestamp,
	})
}

func (s *ClientNotificationService) Read(ctx context.Context, id uint64) (*notificationpb.StatusResponse, error) {
	gw, err := s.gateway()
	if err != nil {
		return nil, err
	}
	return gw.Read(ctx, &notificationpb.IdRequest{Id: id})
}

func (s *ClientNotificationService) ReadMany(ctx context.Context, ids []uint64) (*notificationpb.StatusResponse, error) {
	gw, err := s.gateway()
	if err != nil {
		return nil, err
	}
	return gw.ReadMany(ctx, &notificationpb.ReadManyRequest{Ids: ids})
}

func (s *ClientNotificationService) ReadAll(ctx context.Context) (*notificationpb.StatusResponse, error) {
	gw, err := s.gateway()
	if err != nil {
		return nil, err
	}
	return gw.ReadAll(ctx, &notificationpb.EmptyRequest{})
}

func (s *ClientNotificationService) Remove(ctx context.Context, id uint64) (*notificationpb.StatusResponse, error) {
	gw, err := s.gateway()
	if err != nil {
		return nil, err
	}
	return gw.Remove(ctx, &notificationpb.IdRequest{Id: id})
}
