package client

import (
	base_dto "base/dto"
	base_provider "base/provider"
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	"context"
	"pb/clients"
	"time"

	"crm/env"
	delivery_dto "crm/internal/dto/delivery"
)

type NotificationClient struct {
	*clients.NotificationClient
}

// todo: triển khai nốt thông báo
func NewNotificationClient() base_provider.NotificationProvider {
	return &NotificationClient{
		NotificationClient: clients.MakeNotificationClient(env.ENV_RUNTIME),
	}
}

func (client *NotificationClient) SendNoti(c context.Context, noti *delivery_dto.NotificationDTO) error {
	return nil
}

func (client *NotificationClient) SendReminder(c context.Context, customerID uint64) error {
	return nil
}

// func (client *NotificationClient) CreateHistory(c context.Context, entity *shared_dto.HistoryDTO) error {

// }

// func (client *NotificationClient) ExistedFromDate(c context.Context, customerID uint64, fromDate time.Time) (bool, error) {
// 	return false, nil
// }

func (client *NotificationClient) SendNotification(ctx context.Context, noti *base_dto.NotificationDTO) error {
	if client.NotificationClient == nil {
		return nil
	}
	return client.NotificationClient.CreateNotification(ctx, noti.Title, noti.Message, int32(noti.Type), "", noti.UserID)
}

func (client *NotificationClient) SendBatchNotification(ctx context.Context, noti *base_dto.NotificationBatchDTO) error {
	if client.NotificationClient == nil {
		return nil
	}
	// Convert to protobuf format
	_noti := &_dto.NotificationBatchDTO{}
	for _, data := range noti.Datas {
		_noti.Datas = append(_noti.Datas, _dto.NotificationDTO{
			Title:   data.Title,
			Message: data.Message,
			Type:    data.Type,
		})
	}
	return client.NotificationClient.SendBatchNotification(ctx, _noti)
}

func (client *NotificationClient) CreateCRMHistory(ctx context.Context, history *base_dto.HistoryDTO) error {
	if client.NotificationClient == nil {
		return nil
	}
	// Convert base_dto.HistoryDTO to _dto.HistoryDTO
	historyDTO := &_dto.HistoryDTO{
		Title:      history.Title,
		Note:       history.Note,
		TargetId:   history.TargetId,
		TargetType: _enum.ETargetHistory(history.TargetType),
		ActionType: _enum.EHistory(history.ActionType),
		OwnerID:    history.OwnerID,
		OwnerOf:    _enum.EOwnerOf(history.OwnerOf),
	}
	return client.NotificationClient.CreateCRMHistory(ctx, historyDTO)
}

func (client *NotificationClient) ExistedFromDate(ctx context.Context, customerID uint64, fromDate time.Time) (bool, error) {
	return false, nil
}
