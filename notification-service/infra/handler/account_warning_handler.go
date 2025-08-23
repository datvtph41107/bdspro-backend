package handler

import (
	"context"
	"noti/internal/dto"
	"noti/internal/usecase"
	notificationpb "pb/types/notification"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountWarningHandler struct {
	AccountWarningUsecase *usecase.AccountWarningUsecase
	notificationpb.UnimplementedGatewayNotificationServiceServer
}

func NewAccountWarningHandler(accountWarningUsecase *usecase.AccountWarningUsecase) *AccountWarningHandler {
	return &AccountWarningHandler{
		AccountWarningUsecase: accountWarningUsecase,
	}
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
func (h *AccountWarningHandler) SendAccountWarning(ctx context.Context, req *notificationpb.SendAccountWarningRequest) (*notificationpb.SendAccountWarningResponse, error) {
	// Convert proto request to DTO
	var expiresAt *time.Time
	if req.ExpiresAt != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z07:00", req.ExpiresAt); err == nil {
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
	}

	result, err := h.AccountWarningUsecase.SendAccountWarning(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi gửi cảnh báo: %v", err)
	}

	// Convert DTO response to proto response
	return &notificationpb.SendAccountWarningResponse{
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
func (h *AccountWarningHandler) GetAccountWarningList(ctx context.Context, req *notificationpb.GetAccountWarningListRequest) (*notificationpb.GetAccountWarningListResponse, error) {
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

	result, err := h.AccountWarningUsecase.GetAccountWarningList(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi lấy danh sách cảnh báo: %v", err)
	}

	// Convert DTO response to proto response
	items := make([]*notificationpb.AccountWarningItem, len(result.Data))
	for i, item := range result.Data {
		items[i] = &notificationpb.AccountWarningItem{
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

	return &notificationpb.GetAccountWarningListResponse{
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
func (h *AccountWarningHandler) MarkWarningAsRead(ctx context.Context, req *notificationpb.MarkWarningAsReadRequest) (*notificationpb.MarkWarningAsReadResponse, error) {
	// Convert proto request to DTO
	dto := &dto.MarkWarningAsReadRequest{
		WarningID: req.WarningId,
		TargetID:  req.TargetId,
	}

	result, err := h.AccountWarningUsecase.MarkWarningAsRead(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi đánh dấu đã đọc: %v", err)
	}

	// Convert DTO response to proto response
	return &notificationpb.MarkWarningAsReadResponse{
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
func (h *AccountWarningHandler) MarkWarningAsAcknowledged(ctx context.Context, req *notificationpb.MarkWarningAsAcknowledgedRequest) (*notificationpb.MarkWarningAsAcknowledgedResponse, error) {
	// Convert proto request to DTO
	dto := &dto.MarkWarningAsAcknowledgedRequest{
		WarningID: req.WarningId,
		TargetID:  req.TargetId,
	}

	result, err := h.AccountWarningUsecase.MarkWarningAsAcknowledged(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi đánh dấu đã xác nhận: %v", err)
	}

	// Convert DTO response to proto response
	return &notificationpb.MarkWarningAsAcknowledgedResponse{
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
func (h *AccountWarningHandler) CreateWarningTemplate(ctx context.Context, req *notificationpb.CreateWarningTemplateRequest) (*notificationpb.CreateWarningTemplateResponse, error) {
	// Convert proto request to DTO
	dto := &dto.CreateWarningTemplateRequest{
		Name:        req.Name,
		WarningType: req.WarningType,
		Title:       req.Title,
		Content:     req.Content,
		Severity:    req.Severity,
		IsActive:    req.IsActive,
	}

	result, err := h.AccountWarningUsecase.CreateWarningTemplate(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi tạo mẫu cảnh báo: %v", err)
	}

	// Convert DTO response to proto response
	return &notificationpb.CreateWarningTemplateResponse{
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
func (h *AccountWarningHandler) GetWarningTemplateList(ctx context.Context, req *notificationpb.GetWarningTemplateListRequest) (*notificationpb.GetWarningTemplateListResponse, error) {
	// Convert proto request to DTO
	dto := &dto.GetWarningTemplateListRequest{
		Page:        int(req.Page),
		Size:        int(req.Size),
		WarningType: req.WarningType,
	}

	result, err := h.AccountWarningUsecase.GetWarningTemplateList(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi lấy danh sách mẫu: %v", err)
	}

	// Convert DTO response to proto response
	items := make([]*notificationpb.WarningTemplateItem, len(result.Data))
	for i, item := range result.Data {
		items[i] = &notificationpb.WarningTemplateItem{
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

	return &notificationpb.GetWarningTemplateListResponse{
		Data:  items,
		Total: int32(result.Total),
		Page:  int32(result.Page),
		Size:  int32(result.Size),
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
func (h *AccountWarningHandler) GetWarningLogList(ctx context.Context, req *notificationpb.GetWarningLogListRequest) (*notificationpb.GetWarningLogListResponse, error) {
	// Convert proto request to DTO
	dto := &dto.GetWarningLogListRequest{
		Page:      int(req.Page),
		Size:      int(req.Size),
		WarningID: req.WarningId,
		TargetID:  req.TargetId,
	}

	result, err := h.AccountWarningUsecase.GetWarningLogList(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi lấy danh sách log: %v", err)
	}

	// Convert DTO response to proto response
	items := make([]*notificationpb.WarningLogItem, len(result.Data))
	for i, item := range result.Data {
		items[i] = &notificationpb.WarningLogItem{
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

	return &notificationpb.GetWarningLogListResponse{
		Data:  items,
		Total: int32(result.Total),
		Page:  int32(result.Page),
		Size:  int32(result.Size),
	}, nil
}
