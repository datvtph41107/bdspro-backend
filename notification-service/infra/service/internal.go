package service

import (
	"context"
	"log"
	"noti/internal/domain"
	"noti/internal/dto"
	"noti/internal/enums"
	"noti/internal/usecase"
	pb_noti "pb/types/notification"
	sharepb "pb/types/shared"
)

type InternalNotiService struct {
	Usecase *usecase.NotificationUsecase
	pb_noti.UnimplementedNotificationServiceServer
}

func NewInternalNotiService(uc *usecase.NotificationUsecase) *InternalNotiService {
	return &InternalNotiService{Usecase: uc}
}

func (s *InternalNotiService) Create(ctx context.Context, req *sharepb.NotificationRequest) (*pb_noti.NotificationDTO, error) {
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

	response := s.mapToPB(notification)
	return response, nil
}

func (s *InternalNotiService) SendBatch(ctx context.Context, req *pb_noti.NotiBatchRequest) (*pb_noti.BatchResponse, error) {
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
	return &pb_noti.BatchResponse{
		Data: notificationIds,
	}, nil
}

func (s *InternalNotiService) PushByToken(ctx context.Context, req *pb_noti.PushTokenRequest) (*pb_noti.PushResponse, error) {
	return nil, nil
}

func (s *InternalNotiService) PushByTopic(ctx context.Context, req *pb_noti.PushTopicRequest) (*pb_noti.PushResponse, error) {
	return nil, nil
}

// -- mapping ---
func (s *InternalNotiService) mapToPB(noti *domain.NotificationEntity) *pb_noti.NotificationDTO {
	return &pb_noti.NotificationDTO{
		Id:      noti.ID,
		Title:   noti.Title,
		Message: noti.Message,
		Link:    noti.Link,
		Type:    int32(noti.Type),
	}
}
