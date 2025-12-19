package client

import (
	"bdspro/env"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	_dto "common/domain/dto"
	"context"
	"pb/clients"
	notificationpb "pb/types/notification"

	"google.golang.org/protobuf/types/known/structpb"
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

func (s *NotificationClient) HistoryDeal(ctx context.Context, history *dto.DealHistoryDTO) error {

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
	result, err := s.NotificationClient.Client.CreateDealHistory(ctx, domainReq)
	if err != nil {
		return err
	}
	if result.Success {
		return nil
	}
	return nil
}
