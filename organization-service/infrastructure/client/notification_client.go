package client

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/structpb"
	"organization/internal/dto"
	"organization/internal/enums"
	clients "pb/clients"
	notificationpb "pb/types/notification"
	sharepb "pb/types/shared"
)

// @bind: organization/internal/interface.INotiClient
type NotificationClient struct {
	*clients.NotificationClient
}

func NewNotificationClient(rpcClient *clients.NotificationClient) *NotificationClient {
	return &NotificationClient{NotificationClient: rpcClient}
}

func (c *NotificationClient) SendNotification(ctx context.Context, notification *dto.NotificationDTO) error {
	return c.NotifyToUser(ctx,
		int32(notification.Type),
		notification.UserId,
		notification.AttachData,
		notification.TargetID,
		notification.IsMerge,
	)
}

func (c *NotificationClient) SendBatchNotification(ctx context.Context, notifications []*dto.NotificationDTO) error {
	notiRequests := make([]*sharepb.NotificationRequest, 0)
	for _, notification := range notifications {
		notiRequests = append(notiRequests, &sharepb.NotificationRequest{
			Type:     int32(notification.Type),
			OwnerId:  *notification.UserId,
			TargetId: notification.TargetID,
			// TargetType: int32(notification.Type),
			IsMerge:    notification.IsMerge,
			AttachData: notification.AttachData,
			Link:       fmt.Sprintf("/news-feed/%d", notification.TargetID),
		})
	}
	response, err := c.Client.SendBatch(ctx, &notificationpb.NotiBatchRequest{
		Datas: notiRequests,
	})
	if err != nil {
		return err
	}
	if response.Success {
		return nil
	}
	return nil
}

func (c *NotificationClient) HistoryDeal(ctx context.Context, history *dto.DealHistoryDTO) error {

	pbStruct, err := structpb.NewStruct(history.Metadata)
	if err != nil {
		return err
	}

	domainReq := &notificationpb.CreateDealHistoryRequest{
		DealId:      history.DealID,
		ActorId:     history.ActorID,
		ActorName:   history.ActorName,
		ActorAvatar: history.ActorAvatar,
		ActionType:  enums.DealHistoryEventNameMap[history.ActionType],
		ActionName:  history.ActionName,
		Content:     history.Content,
		Metadata:    pbStruct,
	}
	result, err := c.Client.CreateDealHistory(ctx, domainReq)
	if err != nil {
		return err
	}
	if result.Success {
		return nil
	}
	return nil
}
