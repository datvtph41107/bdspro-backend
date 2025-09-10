package handler

import (
	"context"
	"log"
	shared_enum "pb/enums"
	notificationpb "pb/types/notification"
	sharepb "pb/types/shared"

	"notification/infra/mapper"
	"notification/internal/dto"
	"notification/internal/enums"
	"notification/internal/usecase"
)

type InternalHandler struct {
	notificationpb.UnimplementedInternalNotificationServiceServer
	Usecase             *usecase.NotificationUsecase
	Mapper              *mapper.NotificationMapper
	AdminHistoryUsecase *usecase.AdminHistoryUsecase
	AdminHistoryMapper  *mapper.AdminHistoryMapper
	DealHistoryUsecase  usecase.DealHistoryUsecase
}

func NewInternalHandler(
	uc *usecase.NotificationUsecase,
	mapper *mapper.NotificationMapper,
	adminHistoryUsecase *usecase.AdminHistoryUsecase,
	adminHistoryMapper *mapper.AdminHistoryMapper,
	dealHistoryUsecase usecase.DealHistoryUsecase,
) *InternalHandler {
	return &InternalHandler{
		Usecase:             uc,
		Mapper:              mapper,
		AdminHistoryUsecase: adminHistoryUsecase,
		AdminHistoryMapper:  adminHistoryMapper,
		DealHistoryUsecase:  dealHistoryUsecase,
	}
}

func (s *InternalHandler) Create(ctx context.Context, req *sharepb.NotificationRequest) (*notificationpb.NotificationDTO, error) {
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

func (s *InternalHandler) SendBatch(ctx context.Context, req *notificationpb.NotiBatchRequest) (*notificationpb.BatchResponse, error) {
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

func (s *InternalHandler) PushByToken(ctx context.Context, req *notificationpb.PushTokenRequest) (*notificationpb.PushResponse, error) {
	return nil, nil
}

func (s *InternalHandler) PushByTopic(ctx context.Context, req *notificationpb.PushTopicRequest) (*notificationpb.PushResponse, error) {
	return nil, nil
}

// CreateAdminHistory tạo lịch sử admin (internal)
func (s *InternalHandler) CreateAdminHistory(c context.Context, request *notificationpb.AdminHistoryCreateRequest) (*notificationpb.AdminHistoryDTO, error) {
	// Convert protobuf request to DTO
	dto := &dto.AdminHistoryCreateDTO{
		TargetID:   request.TargetId,
		TargetType: shared_enum.ETargetHistory(request.TargetType),
		ActionType: shared_enum.EHistory(request.ActionType),
		Title:      request.Title,
		Note:       request.Note,
		PreStage:   request.PreStage,
		AfterStage: request.AfterStage,
		AdminID:    request.AdminId,
		OwnerID:    request.OwnerId,
		OwnerType:  shared_enum.EOwnerType(request.OwnerType),
		IsInternal: request.IsInternal,
		AdminRole:  request.AdminRole,
		IPAddress:  request.IpAddress,
		UserAgent:  request.UserAgent,
	}

	entity, err := s.AdminHistoryUsecase.CreateAdminHistory(c, dto)
	if err != nil {
		return nil, err
	}

	responseDTO := s.AdminHistoryMapper.AdminHistoryToDTO(entity)
	return s.convertToPbDTO(responseDTO), nil
}

// CreateAdminHistoryBatch tạo hàng loạt lịch sử admin (internal)
func (s *InternalHandler) CreateAdminHistoryBatch(c context.Context, request *notificationpb.AdminHistoryBatchRequest) (*notificationpb.BatchResponse, error) {
	// Convert protobuf requests to DTOs
	dtos := make([]*dto.AdminHistoryCreateDTO, len(request.Datas))
	for i, data := range request.Datas {
		dtos[i] = &dto.AdminHistoryCreateDTO{
			TargetID:   data.TargetId,
			TargetType: shared_enum.ETargetHistory(data.TargetType),
			ActionType: shared_enum.EHistory(data.ActionType),
			Title:      data.Title,
			Note:       data.Note,
			PreStage:   data.PreStage,
			AfterStage: data.AfterStage,
			AdminID:    data.AdminId,
			OwnerID:    data.OwnerId,
			OwnerType:  shared_enum.EOwnerType(data.OwnerType),
			IsInternal: data.IsInternal,
			AdminRole:  data.AdminRole,
			IPAddress:  data.IpAddress,
			UserAgent:  data.UserAgent,
		}
	}

	ids, err := s.AdminHistoryUsecase.CreateAdminHistoryBatch(c, dtos)
	if err != nil {
		return nil, err
	}

	return &notificationpb.BatchResponse{
		Data: ids,
	}, nil
}

// convertToPbDTO chuyển đổi AdminHistoryDTO thành protobuf AdminHistoryDTO
func (s *InternalHandler) convertToPbDTO(dto *dto.AdminHistoryDTO) *notificationpb.AdminHistoryDTO {
	return &notificationpb.AdminHistoryDTO{
		Id:         dto.ID,
		TargetId:   dto.TargetId,
		TargetType: int32(dto.TargetType),
		ActionType: int32(dto.ActionType),
		Title:      dto.Title,
		Note:       dto.Note,
		PreStage:   dto.PreStage,
		AfterStage: dto.AfterStage,
		AdminId:    dto.AdminID,
		OwnerId:    dto.OwnerID,
		OwnerType:  int32(dto.OwnerType),
		IsInternal: dto.IsInternal,
		AdminRole:  dto.AdminRole,
		IpAddress:  dto.IPAddress,
		UserAgent:  dto.UserAgent,
		CreatedAt:  dto.CreatedAt,
		UpdatedAt:  dto.UpdatedAt,
	}
}

func (s *InternalHandler) CreateDealHistory(ctx context.Context, req *notificationpb.CreateDealHistoryRequest) (*notificationpb.CreateDealHistoryResponse, error) {
	domainReq := mapper.RequestToDealHistoryRequest(req)

	// Create deal history
	response, err := s.DealHistoryUsecase.CreateDealHistory(ctx, domainReq)
	if err != nil {
		return &notificationpb.CreateDealHistoryResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &notificationpb.CreateDealHistoryResponse{
		Success: true,
		Message: "Deal history created successfully",
		Id:      response.ID,
	}, nil
}
