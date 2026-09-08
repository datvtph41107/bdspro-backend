package client

import (
	_dto "common/domain/dto"
	"context"
	"crm/internal/interface/provider"
	"pb/clients"
)

type NotificationClient struct {
	*clients.NotificationClient
}

func NewNotificationClient(rpcClient *clients.NotificationClient) provider.NotificationProvider {
	return &NotificationClient{NotificationClient: rpcClient}
}

// func (c *NotificationClient) CreateNotification(
// 	ctx context.Context,
// 	avatar string,
// 	title string,
// 	message []string,
// 	notificationType _enum.ENotificationType,
// 	targetId *uint64,
// 	ownerId uint64,
// 	ownerOf _enum.EOwnerOf,
// 	attachData []string,
// ) error {
// 	targetIdValue := uint64(0)
// 	if targetId != nil {
// 		targetIdValue = *targetId
// 	}
// 	_, err := c.Client.CreateToOwner(ctx, &sharepb.NotificationRequest{
// 		Avatar:           avatar,
// 		Title:            title,
// 		Message:          message,
// 		NotificationType: uint32(notificationType),
// 		OwnerId:          ownerId,
// 		TargetId:         targetIdValue,
// 		OwnerOf:          uint32(ownerOf),
// 		AttachData:       attachData,
// 		SendToDevice:     true,
// 	})
// 	return err
// }

// func (c *NotificationClient) RemoveNotification(
// 	ctx context.Context,
// 	ownerOf _enum.EOwnerOf,
// 	ownerId uint64,
// 	targetId uint64,
// 	notificationType _enum.ENotificationType,
// ) error {
// 	_, err := c.Client.RemoveToOwner(ctx, &sharepb.NotificationRequest{
// 		OwnerOf:          uint32(ownerOf),
// 		OwnerId:          ownerId,
// 		TargetId:         targetId,
// 		NotificationType: uint32(notificationType),
// 	})
// 	return err
// }

// func (client *NotificationClient) SendNoti(c context.Context, noti *delivery_dto.NotificationDTO) error {
// 	return nil
// }

// func (client *NotificationClient) SendReminder(c context.Context, customerID uint64) error {
// 	return nil
// }

func (client *NotificationClient) CreateHistory(ctx context.Context, history *_dto.HistoryDTO) error {
	return nil
}

// // func (client *NotificationClient) ExistedFromDate(c context.Context, customerID uint64, fromDate time.Time) (bool, error) {
// // 	return false, nil
// // }

// func (client *NotificationClient) SendNotification(ctx context.Context, noti *base_dto.NotificationDTO) error {
// 	if client.NotificationClient == nil {
// 		return nil
// 	}
// 	return client.NotificationClient.CreateNotification(ctx, noti.Title, noti.Message, int32(noti.Type), "", noti.UserID, noti.AttachData)
// }

// func (client *NotificationClient) SendBatchNotification(ctx context.Context, noti *base_dto.NotificationBatchDTO) error {
// 	if client.NotificationClient == nil {
// 		return nil
// 	}
// 	// Convert to protobuf format
// 	_noti := &_dto.NotificationBatchDTO{}
// 	for _, data := range noti.Datas {
// 		_noti.Datas = append(_noti.Datas, _dto.NotificationDTO{
// 			Title:   data.Title,
// 			Message: data.Message,
// 			Type:    data.Type,
// 		})
// 	}
// 	return client.NotificationClient.SendBatchNotification(ctx, _noti)
// }

// func (client *NotificationClient) CreateCRMHistory(ctx context.Context, history *base_dto.HistoryDTO) error {
// 	if client.NotificationClient == nil {
// 		return nil
// 	}
// 	// Convert base_dto.HistoryDTO to _dto.HistoryDTO
// 	historyDTO := &_dto.HistoryDTO{
// 		Title:      history.Title,
// 		Note:       history.Note,
// 		TargetId:   history.TargetId,
// 		TargetType: _enum.ETargetHistory(history.TargetType),
// 		ActionType: _enum.EHistory(history.ActionType),
// 		OwnerID:    history.OwnerID,
// 		OwnerOf:    _enum.EOwnerOf(history.OwnerOf),
// 	}
// 	return client.NotificationClient.CreateCRMHistory(ctx, historyDTO)
// }

// func (client *NotificationClient) ExistedFromDate(ctx context.Context, customerID uint64, fromDate time.Time) (bool, error) {
// 	return false, nil
// }
