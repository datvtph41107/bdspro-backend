package handler

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	"fmt"
	shared_enum "pb/enums"
	notipb "pb/types/notification"
	"time"

	"notification/infra/mapper"
	"notification/internal/dto"
	"notification/internal/enums"
	"notification/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationHandler struct {
	Usecase               *usecase.NotificationUsecase
	AdminHistoryUsecase   *usecase.AdminHistoryUsecase
	AdminHistoryMapper    *mapper.AdminHistoryMapper
	AccountWarningUsecase *usecase.AccountWarningUsecase
	UserClient            usecase.UserProfileReader
	Mapper                *mapper.NotificationMapper
	SyncProvider          *_utils.SyncUtil
	notipb.UnimplementedGatewayNotificationServiceServer
}

func NewNotificationHandler(
	usecase *usecase.NotificationUsecase,
	adminHistoryUsecase *usecase.AdminHistoryUsecase,
	adminHistoryMapper *mapper.AdminHistoryMapper,
	accountWarningUsecase *usecase.AccountWarningUsecase,
	userClient usecase.UserProfileReader,
	mapper *mapper.NotificationMapper,
	syncProvider *_utils.SyncUtil,
) *NotificationHandler {
	return &NotificationHandler{
		Usecase:               usecase,
		AdminHistoryUsecase:   adminHistoryUsecase,
		AdminHistoryMapper:    adminHistoryMapper,
		AccountWarningUsecase: accountWarningUsecase,
		UserClient:            userClient,
		Mapper:                mapper,
		SyncProvider:          syncProvider,
	}
}

// @Summary Tạo thông báo mới
// @Description API này không bắn thông báo tới thiết bị của người dùng, chỉ tạo thông báo
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param notification body notipb.GatewayNotiNewRequest true "Chi tiết thông báo"
// @Success 200 {object} map[string]string "message"
// @Failure 500 {object} map[string]string "error"
// @Router /new [post]
func (s *NotificationHandler) Create(ctx context.Context, req *notipb.GatewayNotiNewRequest) (*notipb.NotificationDTO, error) {
	dto := &dto.NotiNewRequest{
		Title:   req.Title,
		Message: []string{req.Message},
		Link:    req.Link,
		Type:    enums.TypeEnum(req.Type),
	}
	notification, err := s.Usecase.CreateNotification(ctx, dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method Send not implemented")
	}
	s.invalidateNotificationSync(ctx)
	response := s.Mapper.MapToPB(notification)
	return response, nil
}

// @Summary Số thông báo mới chưa đọc
// @Description Hiển thị số lượng thông báo chưa đọc của người dùng
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param timestamp query int64 false "Timestamp lần sync trước (0 = luôn lấy từ DB)"
// @Success 200 {object} map[string]int "count"
// @Failure 500 {object} map[string]string "error"
// @Router /count [get]
func (s *NotificationHandler) Count(ctx context.Context, req *notipb.CountRequest) (*notipb.CountResponse, error) {
	clientTs := req.GetTimestamp()
	key := s.SyncProvider.GetKey(ctx, _utils.SyncKeyNotificationMe, 0)
	updated := s.SyncProvider.HasUpdated(ctx, key, clientTs)
	println("updated", updated)
	if !updated {
		ts := s.SyncProvider.GetTimestamp(ctx, key)
		return &notipb.CountResponse{
			Count:     0,
			Unread:    0,
			Timestamp: ts,
			Unchanged: true,
		}, nil
	}

	count := s.Usecase.Count(ctx)
	s.SyncProvider.PutTimeRequest(ctx, key, 0)

	return &notipb.CountResponse{
		Count:  count,
		Unread: count,
	}, nil
}

// @Summary Đánh dấu thông báo đã đọc
// @Description Đánh dấu thông báo với ID cụ thể là đã đọc
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của thông báo"
// @Success 200 {object} map[string]string "message"
// @Failure 500 {object} map[string]string "error"
// @Router /read/{id} [put]
func (s *NotificationHandler) Read(ctx context.Context, req *notipb.IdRequest) (*notipb.StatusResponse, error) {
	_, err := s.Usecase.Read(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method Read not implemented")
	}
	s.invalidateNotificationSync(ctx)
	return &notipb.StatusResponse{
		Success: true,
		Message: "Thông báo đã được đánh dấu là đã đọc",
	}, nil
}

// @Summary Đánh dấu nhiều thông báo đã đọc
// @Description Đánh dấu nhiều thông báo với danh sách ID là đã đọc
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body notipb.ReadManyRequest true "Danh sách ID thông báo"
// @Success 200 {object} notipb.StatusResponse "message"
// @Failure 500 {object} map[string]string "error"
// @Router /read-many [put]
func (s *NotificationHandler) ReadMany(ctx context.Context, req *notipb.ReadManyRequest) (*notipb.StatusResponse, error) {
	if len(req.Ids) == 0 {
		return &notipb.StatusResponse{
			Success: true,
			Message: "Không có thông báo nào để đánh dấu",
		}, nil
	}
	count, err := s.Usecase.ReadMany(ctx, req.Ids)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to mark notifications as read: %v", err)
	}
	s.invalidateNotificationSync(ctx)
	return &notipb.StatusResponse{
		Success: true,
		Message: fmt.Sprintf("Đã đánh dấu %d thông báo là đã đọc", count),
	}, nil
}

// @Summary Đánh dấu tất cả thông báo là đã đọc
// @Description Đánh dấu tất cả thông báo là đã đọc
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "message"
// @Failure 500 {object} map[string]string "error"
// @Router /read-all [put]
func (s *NotificationHandler) ReadAll(ctx context.Context, req *notipb.EmptyRequest) (*notipb.StatusResponse, error) {
	_, err := s.Usecase.ReadAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method ReadAll not implemented")
	}
	s.invalidateNotificationSync(ctx)
	return &notipb.StatusResponse{
		Success: true,
		Message: "Thông báo đã được đánh dấu là đã đọc",
	}, nil
}

// @Summary Xóa thông báo
// @Description Xóa thông báo với ID cụ thể
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của thông báo"
// @Router /remove/{id} [delete]
func (s *NotificationHandler) Remove(ctx context.Context, req *notipb.IdRequest) (*notipb.StatusResponse, error) {
	_, err := s.Usecase.Remove(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method Remove not implemented")
	}
	s.invalidateNotificationSync(ctx)
	return &notipb.StatusResponse{
		Success: true,
		Message: "Thông báo đã được xóa",
	}, nil
}

// @Summary Lấy danh sách thông báo
// @Description Trả về danh sách thông báo của người dùng
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Trang cần lấy (default = 0)"
// @Param size query int false "Số lượng trên mỗi trang (default = 200)"
// @Param timestamp query int64 false "Timestamp lần sync trước (0 = luôn lấy từ DB)"
// @Success 200 {object} []notipb.NotificationDTO "Danh sách thông báo"
// @Failure 500 {object} map[string]string "error"
// @Router /list [get]
func (s *NotificationHandler) List(ctx context.Context, req *notipb.ListRequest) (*notipb.ListResponse, error) {
	clientTs := req.GetTimestamp()

	key := s.SyncProvider.GetKey(ctx, _utils.SyncKeyNotificationMe, 0)
	updated := s.SyncProvider.HasUpdated(ctx, key, clientTs)
	s.SyncProvider.PutTimeRequest(ctx, key, clientTs)
	if !updated && req.Page == 0 {
		ts := s.SyncProvider.GetTimestamp(ctx, key)
		return &notipb.ListResponse{
			Data:      nil,
			Total:     0,
			Timestamp: ts,
			Unchanged: true,
		}, nil
	}

	query := dto.SearchNotiRequest{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		IsRead: &req.IsRead,
	}
	notifications, total, err := s.Usecase.Search(ctx, &query)
	if err != nil {
		return nil, err
	}
	var responses []*notipb.NotificationDTO
	for _, noti := range notifications {
		responses = append(responses, s.Mapper.MapToPB(&noti))
	}

	resp := &notipb.ListResponse{
		Data:      responses,
		Total:     int32(total),
		Unchanged: false,
	}
	return resp, nil
}

// invalidateNotificationSync cập nhật timestamp khi ds thông báo thay đổi (create/read/remove).
func (s *NotificationHandler) invalidateNotificationSync(ctx context.Context) {
	if s.SyncProvider == nil {
		return
	}
	key := s.SyncProvider.GetKey(ctx, _utils.SyncKeyNotificationMe, 0)
	_ = s.SyncProvider.PutTimestamp(ctx, key, time.Now().UnixMilli())
}

// todo: implement
func (s *NotificationHandler) PushByToken(ctx context.Context, req *notipb.PushTokenRequest) (*notipb.PushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushByToken not implemented")
}

func (s *NotificationHandler) PushByTopic(ctx context.Context, req *notipb.PushTopicRequest) (*notipb.PushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushByTopic not implemented")
}

// @Summary Tìm kiếm lịch sử admin
// @Description Tìm kiếm lịch sử các hành động của admin
// @Tags Admin History
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Trang cần lấy (default = 0)"
// @Param size query int false "Số lượng trên mỗi trang (default = 20)"
// @Param adminId query int false "ID của admin"
// @Param ownerId query int false "ID của owner"
// @Param ownerType query int false "Loại owner"
// @Param targetId query int false "ID của target"
// @Param targetType query int false "Loại target"
// @Param actionType query int false "Loại hành động"
// @Param adminRole query string false "Vai trò admin"
// @Param fromDate query string false "Từ ngày (RFC3339 format)"
// @Param toDate query string false "Đến ngày (RFC3339 format)"
// @Param isInternal query bool false "Có phải internal không"
// @Success 200 {object} notipb.AdminHistoryListResponse "Danh sách lịch sử admin"
// @Failure 500 {object} map[string]string "error"
// @Router /admin-history/search [get]
func (s *NotificationHandler) SearchAdminHistory(ctx context.Context, request *notipb.AdminHistorySearchRequest) (*notipb.AdminHistoryListResponse, error) {
	// Convert protobuf request to DTO
	searchDTO := &dto.AdminHistorySearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(request.Page),
			Size: uint32(request.Size),
		},
	}

	if request.AdminId != nil {
		searchDTO.AdminID = request.AdminId
	}
	if request.OwnerId != nil {
		searchDTO.OwnerID = request.OwnerId
	}
	if request.OwnerType != nil {
		ownerType := shared_enum.EOwnerType(*request.OwnerType)
		searchDTO.OwnerType = &ownerType
	}
	if request.TargetId != nil {
		searchDTO.TargetID = request.TargetId
	}
	if request.TargetType != nil {
		targetType := shared_enum.ETargetHistory(*request.TargetType)
		searchDTO.TargetType = &targetType
	}
	if request.ActionType != nil {
		actionType := shared_enum.EHistory(*request.ActionType)
		searchDTO.ActionType = &actionType
	}
	if request.AdminRole != nil {
		searchDTO.AdminRole = request.AdminRole
	}
	if request.FromDate != nil {
		if fromDate, err := time.Parse(time.RFC3339, *request.FromDate); err == nil {
			searchDTO.FromDate = &fromDate
		}
	}
	if request.ToDate != nil {
		if toDate, err := time.Parse(time.RFC3339, *request.ToDate); err == nil {
			searchDTO.ToDate = &toDate
		}
	}
	if request.IsInternal != nil {
		searchDTO.IsInternal = request.IsInternal
	}

	histories, total, err := s.AdminHistoryUsecase.SearchInternal(ctx, *searchDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to search admin history: %v", err)
	}

	dtos := s.AdminHistoryMapper.EntitiesToDTOs(histories)

	// Convert DTOs to protobuf DTOs
	pbDtos := make([]*notipb.AdminHistoryDTO, len(dtos))
	for i, dto := range dtos {
		pbDtos[i] = s.convertToPbDTO(&dto)
	}

	return &notipb.AdminHistoryListResponse{
		Data:  pbDtos,
		Total: int32(total),
	}, nil
}

// convertToPbDTO chuyển đổi AdminHistoryDTO thành protobuf AdminHistoryDTO
func (s *NotificationHandler) convertToPbDTO(dto *dto.AdminHistoryDTO) *notipb.AdminHistoryDTO {
	pbDTO := &notipb.AdminHistoryDTO{
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

	// Thêm thông tin user nếu có
	if dto.AdminName != nil {
		pbDTO.AdminName = dto.AdminName
	}
	if dto.AdminAvatar != nil {
		pbDTO.AdminAvatar = dto.AdminAvatar
	}

	return pbDTO
}

// @Summary Gửi cảnh báo tài khoản
// @Description Gửi cảnh báo, nhắc nhở hoặc hướng dẫn khắc phục lỗi đến tài khoản người dùng
// @Tags Account Warning
// @Accept json
// @Produce json
// @Param body body notificationpb.SendAccountWarningRequest true "Thông tin cảnh báo"
// @Success 200 {object} notificationpb.SendAccountWarningResponse
// @Failure 400 {object} notificationpb.SendAccountWarningResponse
// @Failure 401 {object} notificationpb.SendAccountWarningResponse
// @Failure 500 {object} notificationpb.SendAccountWarningResponse
// @Router /account-warning/send [post]
func (s *NotificationHandler) SendAccountWarning(ctx context.Context, req *notipb.SendAccountWarningRequest) (*notipb.SendAccountWarningResponse, error) {
	// Convert proto request to DTO
	var expiresAt *time.Time
	if req.ExpiresAt != nil {
		if t, err := time.Parse("2006-01-02T15:04:05Z07:00", *req.ExpiresAt); err == nil {
			expiresAt = &t
		}
	}

	dto := &dto.SendAccountWarningRequest{
		TargetID:    req.TargetId,
		TargetType:  req.TargetType,
		WarningType: req.WarningType,
		Title:       req.Title,
		Content:     req.Content,
		Severity:    req.Severity,
		SendEmail:   req.SendEmail,
		RelatedID:   req.RelatedId,
		RelatedType: req.RelatedType,
		ExpiresAt:   expiresAt,
		TemplateID:  req.TemplateId, // ID mẫu cảnh báo nếu sử dụng
	}

	result, err := s.AccountWarningUsecase.SendAccountWarning(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi gửi cảnh báo: %v", err)
	}

	// Convert DTO response to proto response
	return &notipb.SendAccountWarningResponse{
		Id:          result.ID,
		TargetId:    result.TargetID,
		TargetType:  result.TargetType,
		WarningType: result.WarningType,
		Title:       result.Title,
		Content:     result.Content,
		Severity:    result.Severity,
		Status:      result.Status,
		SentBy:      result.SentBy,
		SentAt:      result.SentAt.Format("2006-01-02T15:04:05Z07:00"),
		EmailSent:   result.EmailSent,
		Message:     result.Message,
	}, nil
}

// @Summary Lấy danh sách cảnh báo tài khoản
// @Description Lấy danh sách cảnh báo tài khoản với filter
// @Tags Account Warning
// @Accept json
// @Produce json
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Param targetId query int false "ID người nhận"
// @Param targetType query string false "Loại đối tượng"
// @Param warningType query string false "Loại cảnh báo"
// @Param severity query string false "Mức độ nghiêm trọng"
// @Param status query string false "Trạng thái"
// @Param sentBy query int false "ID người gửi"
// @Param isRead query bool false "Đã đọc"
// @Success 200 {object} notificationpb.GetAccountWarningListResponse
// @Failure 400 {object} notificationpb.GetAccountWarningListResponse
// @Failure 401 {object} notificationpb.GetAccountWarningListResponse
// @Failure 500 {object} notificationpb.GetAccountWarningListResponse
// @Router /account-warning/list [get]
func (s *NotificationHandler) GetAccountWarningList(ctx context.Context, req *notipb.GetAccountWarningListRequest) (*notipb.GetAccountWarningListResponse, error) {
	// Convert proto request to DTO
	dto := &dto.GetAccountWarningListRequest{
		Page:        int(req.Page),
		Size:        int(req.Size),
		TargetType:  req.TargetType,
		WarningType: req.WarningType,
		Severity:    req.Severity,
		Status:      req.Status,
	}

	if req.TargetId != nil {
		dto.TargetID = req.TargetId
	}
	if req.SentBy != nil {
		dto.SentBy = req.SentBy
	}
	if req.IsRead != nil {
		dto.IsRead = req.IsRead
	}

	result, err := s.AccountWarningUsecase.GetAccountWarningList(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi lấy danh sách cảnh báo: %v", err)
	}

	// Convert DTO response to proto response
	items := make([]*notipb.AccountWarningItem, len(result.Data))
	for i, item := range result.Data {
		items[i] = &notipb.AccountWarningItem{
			Id:          item.ID,
			TargetId:    item.TargetID,
			TargetType:  item.TargetType,
			WarningType: item.WarningType,
			Title:       item.Title,
			Content:     item.Content,
			Severity:    item.Severity,
			Status:      item.Status,
			IsRead:      item.IsRead,
			SentBy:      item.SentBy,
			SentAt:      item.SentAt.Format("2006-01-02T15:04:05Z07:00"),
			EmailSent:   item.EmailSent,
			RelatedId:   item.RelatedID,
			RelatedType: item.RelatedType,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		if item.ReadAt != nil {
			readAtStr := item.ReadAt.Format("2006-01-02T15:04:05Z07:00")
			items[i].ReadAt = &readAtStr
		}
		if item.EmailSentAt != nil {
			emailSentAtStr := item.EmailSentAt.Format("2006-01-02T15:04:05Z07:00")
			items[i].EmailSentAt = &emailSentAtStr
		}
		if item.ExpiresAt != nil {
			expiresAtStr := item.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")
			items[i].ExpiresAt = &expiresAtStr
		}
	}

	return &notipb.GetAccountWarningListResponse{
		Data:  items,
		Total: int32(result.Total),
		Page:  int32(result.Page),
		Size:  int32(result.Size),
	}, nil
}

// @Summary Đánh dấu cảnh báo đã đọc
// @Description Đánh dấu cảnh báo đã được đọc
// @Tags Account Warning
// @Accept json
// @Produce json
// @Param body body notificationpb.MarkWarningAsReadRequest true "Thông tin đánh dấu"
// @Success 200 {object} notificationpb.MarkWarningAsReadResponse
// @Failure 400 {object} notificationpb.MarkWarningAsReadResponse
// @Failure 401 {object} notificationpb.MarkWarningAsReadResponse
// @Failure 403 {object} notificationpb.MarkWarningAsReadResponse
// @Failure 404 {object} notificationpb.MarkWarningAsReadResponse
// @Failure 500 {object} notificationpb.MarkWarningAsReadResponse
// @Router /account-warning/read [put]
func (s *NotificationHandler) MarkWarningAsRead(ctx context.Context, req *notipb.MarkWarningAsReadRequest) (*notipb.MarkWarningAsReadResponse, error) {
	// Convert proto request to DTO
	dto := &dto.MarkWarningAsReadRequest{
		WarningID: req.WarningId,
		TargetID:  req.TargetId,
	}

	result, err := s.AccountWarningUsecase.MarkWarningAsRead(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi đánh dấu đã đọc: %v", err)
	}

	// Convert DTO response to proto response
	return &notipb.MarkWarningAsReadResponse{
		Success: result.Success,
		Message: result.Message,
	}, nil
}

// @Summary Đánh dấu cảnh báo đã xác nhận
// @Description Đánh dấu cảnh báo đã được xác nhận
// @Tags Account Warning
// @Accept json
// @Produce json
// @Param body body notificationpb.MarkWarningAsAcknowledgedRequest true "Thông tin xác nhận"
// @Success 200 {object} notificationpb.MarkWarningAsAcknowledgedResponse
// @Failure 400 {object} notificationpb.MarkWarningAsAcknowledgedResponse
// @Failure 401 {object} notificationpb.MarkWarningAsAcknowledgedResponse
// @Failure 403 {object} notificationpb.MarkWarningAsAcknowledgedResponse
// @Failure 404 {object} notificationpb.MarkWarningAsAcknowledgedResponse
// @Failure 500 {object} notificationpb.MarkWarningAsAcknowledgedResponse
// @Router /account-warning/acknowledge [put]
func (s *NotificationHandler) MarkWarningAsAcknowledged(ctx context.Context, req *notipb.MarkWarningAsAcknowledgedRequest) (*notipb.MarkWarningAsAcknowledgedResponse, error) {
	// Convert proto request to DTO
	dto := &dto.MarkWarningAsAcknowledgedRequest{
		WarningID: req.WarningId,
		TargetID:  req.TargetId,
	}

	result, err := s.AccountWarningUsecase.MarkWarningAsAcknowledged(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi đánh dấu đã xác nhận: %v", err)
	}

	// Convert DTO response to proto response
	return &notipb.MarkWarningAsAcknowledgedResponse{
		Success: result.Success,
		Message: result.Message,
	}, nil
}

// @Summary Tạo mẫu cảnh báo
// @Description Tạo mẫu cảnh báo mới
// @Tags Account Warning Template
// @Accept json
// @Produce json
// @Param body body notificationpb.CreateWarningTemplateRequest true "Thông tin mẫu"
// @Success 200 {object} notificationpb.CreateWarningTemplateResponse
// @Failure 400 {object} notificationpb.CreateWarningTemplateResponse
// @Failure 401 {object} notificationpb.CreateWarningTemplateResponse
// @Failure 500 {object} notificationpb.CreateWarningTemplateResponse
// @Router /account-warning/template [post]
func (s *NotificationHandler) CreateWarningTemplate(ctx context.Context, req *notipb.CreateWarningTemplateRequest) (*notipb.CreateWarningTemplateResponse, error) {
	// Convert proto request to DTO
	dto := &dto.CreateWarningTemplateRequest{
		Name:        req.Name,
		WarningType: req.WarningType,
		Title:       req.Title,
		Content:     req.Content,
		Severity:    req.Severity,
		IsActive:    req.IsActive,
	}

	result, err := s.AccountWarningUsecase.CreateWarningTemplate(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi tạo mẫu cảnh báo: %v", err)
	}

	// Convert DTO response to proto response
	return &notipb.CreateWarningTemplateResponse{
		Id:          result.ID,
		Name:        result.Name,
		WarningType: result.WarningType,
		Title:       result.Title,
		Content:     result.Content,
		Severity:    result.Severity,
		IsActive:    result.IsActive,
		CreatedBy:   result.CreatedBy,
		UpdatedBy:   result.UpdatedBy,
		CreatedAt:   result.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   result.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// @Summary Lấy danh sách mẫu cảnh báo
// @Description Lấy danh sách mẫu cảnh báo
// @Tags Account Warning Template
// @Accept json
// @Produce json
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Param warningType query string false "Loại cảnh báo"
// @Success 200 {object} notificationpb.GetWarningTemplateListResponse
// @Failure 400 {object} notificationpb.GetWarningTemplateListResponse
// @Failure 401 {object} notificationpb.GetWarningTemplateListResponse
// @Failure 500 {object} notificationpb.GetWarningTemplateListResponse
// @Router /account-warning/template/list [get]
func (s *NotificationHandler) GetWarningTemplateList(ctx context.Context, req *notipb.GetWarningTemplateListRequest) (*notipb.GetWarningTemplateListResponse, error) {
	// Convert proto request to DTO
	dto := &dto.GetWarningTemplateListRequest{
		Page:        int(req.Page),
		Size:        int(req.Size),
		WarningType: req.WarningType,
	}

	result, err := s.AccountWarningUsecase.GetWarningTemplateList(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi lấy danh sách mẫu: %v", err)
	}

	// Convert DTO response to proto response
	items := make([]*notipb.WarningTemplateItem, len(result.Data))
	for i, item := range result.Data {
		items[i] = &notipb.WarningTemplateItem{
			Id:          item.ID,
			Name:        item.Name,
			WarningType: item.WarningType,
			Title:       item.Title,
			Content:     item.Content,
			Severity:    item.Severity,
			IsActive:    item.IsActive,
			CreatedBy:   item.CreatedBy,
			UpdatedBy:   item.UpdatedBy,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &notipb.GetWarningTemplateListResponse{
		Data:  items,
		Total: int32(result.Total),
		Page:  int32(result.Page),
		Size:  int32(result.Size),
	}, nil
}

// @Summary Lấy chi tiết mẫu cảnh báo
// @Description Lấy chi tiết mẫu cảnh báo theo ID
// @Tags Account Warning Template
// @Accept json
// @Produce json
// @Param id path int true "ID mẫu cảnh báo"
// @Success 200 {object} notificationpb.GetWarningTemplateDetailResponse
// @Failure 401 {object} notificationpb.GetWarningTemplateDetailResponse
// @Failure 404 {object} notificationpb.GetWarningTemplateDetailResponse
// @Failure 500 {object} notificationpb.GetWarningTemplateDetailResponse
// @Router /account-warning/template/{id} [get]
func (s *NotificationHandler) GetWarningTemplateDetail(ctx context.Context, req *notipb.GetWarningTemplateDetailRequest) (*notipb.GetWarningTemplateDetailResponse, error) {
	result, err := s.AccountWarningUsecase.GetWarningTemplateDetail(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi lấy chi tiết mẫu: %v", err)
	}

	return &notipb.GetWarningTemplateDetailResponse{
		Id:          result.ID,
		Name:        result.Name,
		WarningType: result.WarningType,
		Title:       result.Title,
		Content:     result.Content,
		Severity:    result.Severity,
		IsActive:    result.IsActive,
		CreatedBy:   result.CreatedBy,
		UpdatedBy:   result.UpdatedBy,
		CreatedAt:   result.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   result.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// @Summary Lấy danh sách log cảnh báo
// @Description Lấy danh sách log cảnh báo
// @Tags Account Warning Log
// @Accept json
// @Produce json
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Param warningId query int false "ID cảnh báo"
// @Param targetId query int false "ID người nhận"
// @Success 200 {object} notificationpb.GetWarningLogListResponse
// @Failure 400 {object} notificationpb.GetWarningLogListResponse
// @Failure 401 {object} notificationpb.GetWarningLogListResponse
// @Failure 500 {object} notificationpb.GetWarningLogListResponse
// @Router /account-warning/log/list [get]
func (s *NotificationHandler) GetWarningLogList(ctx context.Context, req *notipb.GetWarningLogListRequest) (*notipb.GetWarningLogListResponse, error) {
	// Convert proto request to DTO
	dto := &dto.GetWarningLogListRequest{
		Page:      int(req.Page),
		Size:      int(req.Size),
		WarningID: req.WarningId,
		TargetID:  req.TargetId,
	}

	result, err := s.AccountWarningUsecase.GetWarningLogList(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi lấy danh sách log: %v", err)
	}

	// Convert DTO response to proto response
	items := make([]*notipb.WarningLogItem, len(result.Data))
	for i, item := range result.Data {
		items[i] = &notipb.WarningLogItem{
			Id:          item.ID,
			WarningId:   item.WarningID,
			TargetId:    item.TargetID,
			Action:      item.Action,
			Status:      item.Status,
			Reason:      item.Reason,
			PerformedBy: item.PerformedBy,
			PerformedAt: item.PerformedAt.Format("2006-01-02T15:04:05Z07:00"),
			IpAddress:   item.IPAddress,
			UserAgent:   item.UserAgent,
		}
	}

	return &notipb.GetWarningLogListResponse{
		Data:  items,
		Total: int32(result.Total),
		Page:  int32(result.Page),
		Size:  int32(result.Size),
	}, nil
}

// System Warning APIs - Admin

// @Summary Lấy danh sách cảnh báo hệ thống (Admin)
// @Description Admin lấy danh sách tất cả cảnh báo hệ thống
// @Tags System Warning - Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Param sort query string false "Sắp xếp"
// @Param warningType query string false "Loại cảnh báo"
// @Param severity query string false "Mức độ nghiêm trọng"
// @Param targetId query int false "ID người nhận"
// @Param targetType query string false "Loại người nhận"
// @Param isRead query bool false "Đã đọc chưa"
// @Success 200 {object} notificationpb.GetSystemWarningListResponse
// @Failure 401 {object} notificationpb.GetSystemWarningListResponse
// @Failure 500 {object} notificationpb.GetSystemWarningListResponse
// @Router /admin/system-warning [get]
func (s *NotificationHandler) GetAdminSystemWarningList(ctx context.Context, req *notipb.GetSystemWarningListRequest) (*notipb.GetSystemWarningListResponse, error) {
	// Convert proto request to DTO
	dto := &dto.GetAccountWarningListRequest{
		Page:        int(req.Page),
		Size:        int(req.Size),
		WarningType: "system", // Chỉ lấy cảnh báo hệ thống
	}

	if req.TargetId != nil {
		dto.TargetID = req.TargetId
	}
	if req.TargetType != nil {
		dto.TargetType = *req.TargetType
	}
	if req.Severity != nil {
		dto.Severity = *req.Severity
	}
	if req.IsRead != nil {
		dto.IsRead = req.IsRead
	}

	result, err := s.AccountWarningUsecase.GetAccountWarningList(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi lấy danh sách cảnh báo hệ thống: %v", err)
	}

	// Convert DTO response to proto response
	items := make([]*notipb.SystemWarningItem, len(result.Data))
	for i, item := range result.Data {
		items[i] = &notipb.SystemWarningItem{
			Id:          item.ID,
			TargetId:    item.TargetID,
			TargetType:  item.TargetType,
			WarningType: item.WarningType,
			Title:       item.Title,
			Content:     item.Content,
			Severity:    item.Severity,
			Status:      item.Status,
			IsRead:      item.IsRead,
			SentBy:      item.SentBy,
			SentAt:      item.SentAt.Format("2006-01-02T15:04:05Z07:00"),
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		if item.ReadAt != nil {
			readAtStr := item.ReadAt.Format("2006-01-02T15:04:05Z07:00")
			items[i].ReadAt = &readAtStr
		}
	}

	return &notipb.GetSystemWarningListResponse{
		Data:  items,
		Total: int32(result.Total),
	}, nil
}

// @Summary Tạo cảnh báo hệ thống (Admin)
// @Description Admin tạo cảnh báo hệ thống gửi tới user
// @Tags System Warning - Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body notificationpb.CreateSystemWarningRequest true "Thông tin cảnh báo"
// @Success 200 {object} notificationpb.CreateSystemWarningResponse
// @Failure 400 {object} notificationpb.CreateSystemWarningResponse
// @Failure 401 {object} notificationpb.CreateSystemWarningResponse
// @Failure 500 {object} notificationpb.CreateSystemWarningResponse
// @Router /admin/system-warning [post]
func (s *NotificationHandler) CreateSystemWarning(ctx context.Context, req *notipb.CreateSystemWarningRequest) (*notipb.CreateSystemWarningResponse, error) {
	// Convert proto request to DTO
	dto := &dto.SendAccountWarningRequest{
		TargetID:    req.TargetId,
		TargetType:  "user", // System warning luôn gửi cho user
		WarningType: "system",
		Title:       req.Title,
		Content:     req.Content,
		Severity:    req.Severity,
		SendEmail:   false,
	}

	result, err := s.AccountWarningUsecase.SendAccountWarning(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi tạo cảnh báo hệ thống: %v", err)
	}

	// Convert DTO response to proto response
	return &notipb.CreateSystemWarningResponse{
		Id:          result.ID,
		WarningType: result.WarningType,
		Title:       result.Title,
		Content:     result.Content,
		Severity:    result.Severity,
		TargetId:    result.TargetID,
		SentBy:      result.SentBy,
		SentAt:      result.SentAt.Format("2006-01-02T15:04:05Z07:00"),
		Message:     result.Message,
	}, nil
}

// @Summary Xóa cảnh báo hệ thống (Admin)
// @Description Admin xóa cảnh báo hệ thống
// @Tags System Warning - Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID cảnh báo"
// @Success 200 {object} notificationpb.DeleteSystemWarningResponse
// @Failure 400 {object} notificationpb.DeleteSystemWarningResponse
// @Failure 401 {object} notificationpb.DeleteSystemWarningResponse
// @Failure 404 {object} notificationpb.DeleteSystemWarningResponse
// @Failure 500 {object} notificationpb.DeleteSystemWarningResponse
// @Router /admin/system-warning/{id} [delete]
func (s *NotificationHandler) DeleteSystemWarning(ctx context.Context, req *notipb.DeleteSystemWarningRequest) (*notipb.DeleteSystemWarningResponse, error) {
	// Xóa cảnh báo (soft delete)
	err := s.AccountWarningUsecase.AccountWarningRepo.SoftDelete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi xóa cảnh báo hệ thống: %v", err)
	}

	return &notipb.DeleteSystemWarningResponse{
		Success: true,
		Message: "Đã xóa cảnh báo thành công",
	}, nil
}

// System Warning APIs - User

// @Summary Lấy cảnh báo hệ thống của tôi (User)
// @Description User lấy danh sách cảnh báo hệ thống của mình
// @Tags System Warning - User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Param sort query string false "Sắp xếp"
// @Param warningType query string false "Loại cảnh báo"
// @Param isRead query bool false "Đã đọc chưa"
// @Success 200 {object} notificationpb.GetMySystemWarningsResponse
// @Failure 401 {object} notificationpb.GetMySystemWarningsResponse
// @Failure 500 {object} notificationpb.GetMySystemWarningsResponse
// @Router /system-warning/me [get]
func (s *NotificationHandler) GetMySystemWarnings(ctx context.Context, req *notipb.GetMySystemWarningsRequest) (*notipb.GetMySystemWarningsResponse, error) {
	// Lấy profileID từ context
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "không có quyền truy cập")
	}

	// Convert proto request to DTO
	dto := &dto.GetAccountWarningListRequest{
		Page:        int(req.Page),
		Size:        int(req.Size),
		TargetID:    &profileID,
		TargetType:  "user",
		WarningType: "system", // Chỉ lấy cảnh báo hệ thống
	}

	if req.IsRead != nil {
		dto.IsRead = req.IsRead
	}

	result, err := s.AccountWarningUsecase.GetAccountWarningList(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi lấy danh sách cảnh báo: %v", err)
	}

	// Convert DTO response to proto response
	items := make([]*notipb.SystemWarningItem, len(result.Data))
	for i, item := range result.Data {
		items[i] = &notipb.SystemWarningItem{
			Id:          item.ID,
			TargetId:    item.TargetID,
			TargetType:  item.TargetType,
			WarningType: item.WarningType,
			Title:       item.Title,
			Content:     item.Content,
			Severity:    item.Severity,
			Status:      item.Status,
			IsRead:      item.IsRead,
			SentBy:      item.SentBy,
			SentAt:      item.SentAt.Format("2006-01-02T15:04:05Z07:00"),
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		if item.ReadAt != nil {
			readAtStr := item.ReadAt.Format("2006-01-02T15:04:05Z07:00")
			items[i].ReadAt = &readAtStr
		}
	}

	return &notipb.GetMySystemWarningsResponse{
		Data:  items,
		Total: int32(result.Total),
	}, nil
}
