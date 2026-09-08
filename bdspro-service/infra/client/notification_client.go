package client

import (
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/provider"
	_dto "common/domain/dto"
	_errors "common/errors"
	"context"
	"log"
	"pb/clients"
	notificationpb "pb/types/notification"
	sharepb "pb/types/shared"

	"google.golang.org/protobuf/types/known/structpb"
)

type NotificationClient struct {
	*clients.NotificationClient
}

// todo: triển khai nốt thông báo
func NewNotificationClient(
	rpcClient *clients.NotificationClient,
) provider.NotificationProvider {
	return &NotificationClient{
		NotificationClient: rpcClient,
	}
}

func (s *NotificationClient) RegistedEventProperty(ctx context.Context, req *dto.CreatePropertyActivityDTO) error {
	if err := s.CreatePropertyHistoryActivity(ctx,
		req.SubjectID,
		uint32(req.SubjectType),
		uint32(req.Action),
		req.Description,
		req.ActorID,
	); err != nil {
		log.Printf("create PropertyHistory failed: %v", err)
		return _errors.InternalServerException("create PropertyHistory failed")
	}

	return nil
}

func (s *NotificationClient) UserInGroup(ctx context.Context, userID uint64, groupID uint64) (bool, error) {
	return true, nil
}

func (s *NotificationClient) SendNoti(ctx context.Context, noti *_dto.NotificationDTO) error {
	_, err := s.NotificationClient.Client.Create(ctx, &sharepb.NotificationRequest{
		Avatar:       noti.Avatar,
		Title:        noti.Title,
		Message:      noti.Message,
		Type:         int32(noti.NotificationType),
		OwnerId:      noti.OwnerID,
		TargetId:     noti.TargetID,
		OwnerOf:      uint32(noti.OwnerOf),
		AttachData:   noti.AttachData,
		SendToDevice: noti.SendToDevice,
		IsMerge:      noti.IsMerge,
	})
	return err
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
