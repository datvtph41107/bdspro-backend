package service

import (
	"context"
	"noti/internal/domain"
	"noti/internal/dto"
	"noti/internal/enums"
	"noti/internal/usecase"
	pb_noti "pb/types/notification"
)

type InternalNotiService struct {
	Usecase usecase.NotificationUsecase
	pb_noti.UnimplementedNotificationServiceServer
}

func NewInternalNotiService() *InternalNotiService {
	return &InternalNotiService{}
}

func (s *InternalNotiService) Create(ctx context.Context, req *pb_noti.NotiNewRequest) (*pb_noti.NotificationDTO, error) {
	dto := &dto.NotiNewRequest{
		Title:        req.Title,
		Message:      req.Message,
		Link:         req.Link,
		Type:         enums.TypeEnum(req.Type),
		TargetID:     req.TargetId,
		OwnerType:    enums.OwnerType(req.TargetType),
		SendToDevice: req.SendToDevice,
	}
	notification, err := s.Usecase.Create(ctx, dto)
	if err != nil {
		return nil, err
	}

	response := s.mapToPB(notification)
	return response, nil
}

func (s *InternalNotiService) SendBatch(ctx context.Context, req *pb_noti.NotiBatchRequest) (*pb_noti.ListResponse, error) {
	return nil, nil
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
		Message: []string{noti.Message},
		Link:    noti.Link,
		Type:    int32(noti.Type),
	}
}
