package handler

import (
	_enum "common/domain/enum"
	"context"
	"log"
	shared_enum "pb/enums"
	notificationpb "pb/types/notification"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"notification/infra/mapper"
	"notification/internal/dto"
	"notification/internal/enums"
	"notification/internal/usecase"
)

type InternalHandler struct {
	notificationpb.UnimplementedInternalNotificationServiceServer
	Usecase                *usecase.NotificationUsecase
	Mapper                 *mapper.NotificationMapper
	AdminHistoryUsecase    *usecase.AdminHistoryUsecase
	AdminHistoryMapper     *mapper.AdminHistoryMapper
	DealHistoryUsecase     *usecase.DealHistoryUsecase
	HistoryAuthUsecase     *usecase.HistoryAuthUsecase
	HistoryAuthMapper      *mapper.HistoryAuthMapper
	PersonConfigUsecase    *usecase.PersonConfigUsecase
	PropertyHistoryUsecase *usecase.PropertyHistoryUseCase
}

func NewInternalHandler(
	uc *usecase.NotificationUsecase,
	mapper *mapper.NotificationMapper,
	adminHistoryUsecase *usecase.AdminHistoryUsecase,
	adminHistoryMapper *mapper.AdminHistoryMapper,
	dealHistoryUsecase *usecase.DealHistoryUsecase,
	historyAuthUsecase *usecase.HistoryAuthUsecase,
	historyAuthMapper *mapper.HistoryAuthMapper,
	personConfigUsecase *usecase.PersonConfigUsecase,
	propertyHistoryUsecase *usecase.PropertyHistoryUseCase,
) *InternalHandler {
	return &InternalHandler{
		Usecase:                uc,
		Mapper:                 mapper,
		AdminHistoryUsecase:    adminHistoryUsecase,
		AdminHistoryMapper:     adminHistoryMapper,
		DealHistoryUsecase:     dealHistoryUsecase,
		HistoryAuthUsecase:     historyAuthUsecase,
		HistoryAuthMapper:      historyAuthMapper,
		PersonConfigUsecase:    personConfigUsecase,
		PropertyHistoryUsecase: propertyHistoryUsecase,
	}
}

func (s *InternalHandler) Create(ctx context.Context, req *sharepb.NotificationRequest) (*notificationpb.NotificationDTO, error) {
	dto := &dto.NotiNewRequest{
		Title:        req.Title,
		Message:      req.Message,
		AttachData:   req.AttachData,
		Link:         req.Link,
		Type:         enums.TypeEnum(req.Type),
		TargetID:     req.TargetId,
		OwnerID:      req.OwnerId,
		OwnerOf:      _enum.EOwnerOf(req.OwnerOf),
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
			Avatar:       noti.Avatar,
			Title:        noti.Title,
			Message:      noti.Message,
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
	dto := &dto.PushTokenRequest{
		Token: req.Token,
		Title: req.Title,
		Body:  req.Body,
	}
	return s.Usecase.PushByToken(ctx, dto)
}

func (s *InternalHandler) PushByTopic(ctx context.Context, req *notificationpb.PushTopicRequest) (*notificationpb.PushResponse, error) {
	dto := &dto.PushTopicRequest{
		Topic: req.Topic,
		Title: req.Title,
		Body:  req.Body,
	}
	return s.Usecase.PushByTopic(ctx, dto)
}

func (s *InternalHandler) PushByUserId(ctx context.Context, req *notificationpb.PushByUserIdRequest) (*notificationpb.PushResponse, error) {
	dto := &dto.PushByUserIdRequest{
		UserID: req.UserId,
		Title:  req.Title,
		Body:   req.Body,
		Data:   req.Data,
	}
	return s.Usecase.PushByUserId(ctx, dto)
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

func (s *InternalHandler) CreateHistoryAuth(ctx context.Context, req *notificationpb.HistoryAuthCreateRequest) (*notificationpb.HistoryAuthDTO, error) {
	if s.HistoryAuthUsecase == nil {
		return nil, status.Error(codes.FailedPrecondition, "history auth usecase is not configured")
	}

	dto := &dto.HistoryAuthCreateDTO{
		UserID:         req.UserId,
		OrganizationID: req.OrganizationId,
		ActionType:     req.ActionType,
		ActionName:     req.ActionName,
		Description:    req.Description,
		Success:        &req.Success,
		Reason:         req.Reason,
		IPAddress:      req.IpAddress,
		UserAgent:      req.UserAgent,
		PerformedBy:    req.PerformedBy,
		SessionID:      req.SessionId,
		Channel:        req.Channel,
		DeviceID:       req.DeviceId,
		Location:       req.Location,
		AdditionalNote: req.AdditionalNote,
		SourceService:  req.SourceService,
	}

	if req.Metadata != nil {
		dto.Metadata = req.Metadata.AsMap()
	}

	entity, err := s.HistoryAuthUsecase.CreateHistoryAuth(ctx, dto)
	if err != nil {
		return nil, err
	}

	if s.HistoryAuthMapper == nil {
		return nil, status.Error(codes.FailedPrecondition, "history auth mapper is not configured")
	}

	return s.HistoryAuthMapper.DomainToPb(entity), nil
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

func (s *InternalHandler) CreateToOwner(ctx context.Context, req *sharepb.NotificationRequest) (*notificationpb.NotificationDTO, error) {
	dto := &dto.NotiNewRequest{
		Avatar:           req.Avatar,
		Title:            req.Title,
		Message:          req.Message,
		AttachData:       req.AttachData,
		NotificationType: _enum.ENotificationType(req.NotificationType),
		Type:             enums.TypeEnum(req.Type),
		TargetID:         req.TargetId,
		OwnerID:          req.OwnerId,
		OwnerOf:          _enum.EOwnerOf(req.OwnerOf),
		SendToDevice:     req.SendToDevice,
		IsMerge:          req.IsMerge,
	}
	log.Println("dto", dto)
	notification, err := s.Usecase.CreateToOwner(ctx, dto)
	if err != nil {
		return nil, err
	}

	response := s.Mapper.MapToPB(notification)
	return response, nil
}

func (s *InternalHandler) RemoveToOwner(ctx context.Context, req *sharepb.NotificationRequest) (*notificationpb.NotificationDTO, error) {
	dto := &dto.NotiNewRequest{
		OwnerID:          req.OwnerId,
		TargetID:         req.TargetId,
		NotificationType: _enum.ENotificationType(req.NotificationType),
		OwnerOf:          _enum.EOwnerOf(req.OwnerOf),
	}
	err := s.Usecase.RemoveToOwner(ctx, dto)
	if err != nil {
		return nil, err
	}
	return &notificationpb.NotificationDTO{}, nil
}

// BatchPersonConfig cập nhật danh sách cấu hình cá nhân
func (s *InternalHandler) BatchPersonConfig(ctx context.Context, req *notificationpb.BatchPersonConfigRequest) (*notificationpb.BatchPersonConfigResponse, error) {
	if s.PersonConfigUsecase == nil {
		return &notificationpb.BatchPersonConfigResponse{
			Success: false,
		}, nil
	}

	configDTOs := make([]*dto.PersonConfigDTO, len(req.Configs))
	for i, cfg := range req.Configs {
		configDTOs[i] = &dto.PersonConfigDTO{
			ID:        cfg.Id,
			UserID:    cfg.UserId,
			Key:       cfg.Key,
			Checked:   cfg.Checked,
			IsDefault: cfg.IsDefault,
			Channel:   _enum.EChannelNotification(cfg.Channel),
		}
	}

	entities, err := s.PersonConfigUsecase.BatchUpsert(ctx, configDTOs)
	if err != nil {
		return nil, err
	}

	response := &notificationpb.BatchPersonConfigResponse{
		Success: true,
		Data:    make([]*sharepb.PersonConfigV3Proto, len(entities)),
	}

	for i, entity := range entities {
		response.Data[i] = &sharepb.PersonConfigV3Proto{
			Id:        entity.ID,
			UserId:    entity.UserID,
			Key:       entity.Key,
			Checked:   entity.Checked,
			IsDefault: entity.IsDefault,
			Channel:   uint32(entity.Channel),
		}
	}

	return response, nil
}

// GetPersonConfigs trả về danh sách cấu hình cá nhân theo userId
func (s *InternalHandler) GetPersonConfigs(ctx context.Context, req *notificationpb.GetPersonConfigsRequest) (*notificationpb.GetPersonConfigsResponse, error) {
	if req == nil || req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "userId là bắt buộc")
	}

	if s.PersonConfigUsecase == nil {
		return &notificationpb.GetPersonConfigsResponse{
			Success: false,
			Data:    []*sharepb.PersonConfigV3Proto{},
		}, nil
	}

	entities, err := s.PersonConfigUsecase.ListByUserID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	response := &notificationpb.GetPersonConfigsResponse{
		Success: true,
		Data:    make([]*sharepb.PersonConfigV3Proto, 0, len(entities)),
	}

	for _, entity := range entities {
		response.Data = append(response.Data, &sharepb.PersonConfigV3Proto{
			Id:        entity.ID,
			UserId:    entity.UserID,
			Key:       entity.Key,
			Checked:   entity.Checked,
			IsDefault: entity.IsDefault,
			Channel:   uint32(entity.Channel),
		})
	}

	return response, nil
}

func (s *InternalHandler) CreatePropertyHistoryActivity(
	ctx context.Context,
	req *notificationpb.CreatePropertyHistoryRequest,
) (*notificationpb.CreatePropertyHistoryResponse, error) {
	args := &dto.CreatePropertyHistoryDTO{
		SubjectID:   req.SubjectId,
		Action:      _enum.PropertyAction(req.Action),
		Description: req.Description,
		ActorID:     req.ActorId,
	}

	property, err := s.PropertyHistoryUsecase.Create(ctx, args)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &notificationpb.CreatePropertyHistoryResponse{
		Success: true,
		Id:      property.ID,
	}, nil
}
