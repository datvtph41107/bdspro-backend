package service

import (
	"context"
	"log"
	"noti/infra/mapper"
	"noti/internal/dto"
	"noti/internal/enums"
	"noti/internal/usecase"
	notificationpb "pb/types/notification"
	sharepb "pb/types/shared"
)

type InternalNotiService struct {
	Usecase *usecase.NotificationUsecase
	Mapper  *mapper.NotificationMapper
	notificationpb.UnimplementedNotificationServiceServer
}

func NewInternalNotiService(
	uc *usecase.NotificationUsecase,
	mapper *mapper.NotificationMapper,
) *InternalNotiService {
	return &InternalNotiService{
		Usecase: uc,
		Mapper:  mapper,
	}
}

func (s *InternalNotiService) Create(ctx context.Context, req *sharepb.NotificationRequest) (*notificationpb.NotificationDTO, error) {
	dto := &dto.NotiNewRequest{
		Title:        req.Title,
		Message:      []string{req.Message},
		AttachData:   req.AttachData,
		Link:         req.Link,
		Type:         enums.TypeEnum(req.Type),
		TargetID:     req.TargetId,
		OwnerID:      req.OwnerId,
		OwnerType:    enums.OwnerType(req.TargetType),
		SendToDevice: req.SendToDevice,
		IsMerge:      req.IsMerge,
	}
	log.Println("dto", dto)
	notification, err := s.Usecase.Create(ctx, dto)
	if err != nil {
		return nil, err
	}

	response := s.Mapper.MapToPB(notification)
	return response, nil
}

func (s *InternalNotiService) SendBatch(ctx context.Context, req *notificationpb.NotiBatchRequest) (*notificationpb.BatchResponse, error) {
	notifications := make([]*dto.NotiNewRequest, 0)
	for _, noti := range req.Datas {
		notifications = append(notifications, &dto.NotiNewRequest{
			Title:        noti.Title,
			Message:      []string{noti.Message},
			Type:         enums.TypeEnum(noti.Type),
			AttachData:   noti.AttachData,
			TargetID:     noti.TargetId,
			OwnerID:      noti.OwnerId,
			OwnerType:    enums.OwnerType(noti.TargetType),
			SendToDevice: noti.SendToDevice,
			IsMerge:      noti.IsMerge,
		})
	}
	notificationIds, err := s.Usecase.SendBatch(ctx, notifications)
	if err != nil {
		return nil, err
	}
	return &notificationpb.BatchResponse{
		Data: notificationIds,
	}, nil
}

func (s *InternalNotiService) PushByToken(ctx context.Context, req *notificationpb.PushTokenRequest) (*notificationpb.PushResponse, error) {
	return nil, nil
}

func (s *InternalNotiService) PushByTopic(ctx context.Context, req *notificationpb.PushTopicRequest) (*notificationpb.PushResponse, error) {
	return nil, nil
}
