package handler

import (
	"auth/internal/dto"
	"auth/internal/usecase"
	"context"
	authpb "pb/types/auth"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminAccessHandler struct {
	AdminAccessUsecase *usecase.AdminAccessUsecase
	authpb.UnimplementedAuthServiceServer
}

func NewAdminAccessHandler(adminAccessUsecase *usecase.AdminAccessUsecase) *AdminAccessHandler {
	return &AdminAccessHandler{
		AdminAccessUsecase: adminAccessUsecase,
	}
}

// @Summary Tạo quyền truy cập admin
// @Description Tạo quyền truy cập hệ thống nội bộ mới
// @Tags Admin Access Control
// @Accept json
// @Produce json
// @Param body body authpb.CreateAdminAccessRequest true "Thông tin quyền truy cập"
// @Success 200 {object} authpb.AdminAccessResponse
// @Failure 400 {object} authpb.AdminAccessResponse
// @Failure 401 {object} authpb.AdminAccessResponse
// @Failure 500 {object} authpb.AdminAccessResponse
// @Router /admin-access [post]
func (h *AdminAccessHandler) CreateAdminAccess(ctx context.Context, req *authpb.CreateAdminAccessRequest) (*authpb.AdminAccessResponse, error) {
	// Convert proto request to DTO
	var effectiveFrom, effectiveTo *time.Time
	if req.EffectiveFrom != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z07:00", req.EffectiveFrom); err == nil {
			effectiveFrom = &t
		}
	}
	if req.EffectiveTo != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z07:00", req.EffectiveTo); err == nil {
			effectiveTo = &t
		}
	}

	dto := &dto.AdminAccessRequest{
		UserID:        req.UserId,
		IPAddress:     req.IpAddress,
		IPRange:       req.IpRange,
		DeviceID:      req.DeviceId,
		DeviceName:    req.DeviceName,
		DeviceType:    req.DeviceType,
		AuthType:      req.AuthType,
		Status:        req.Status,
		EffectiveFrom: effectiveFrom,
		EffectiveTo:   effectiveTo,
		MaxDevices:    int(req.MaxDevices),
		Require2FA:    req.Require2Fa,
	}

	result, err := h.AdminAccessUsecase.CreateAdminAccess(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi tạo quyền truy cập: %v", err)
	}

	// Convert DTO response to proto response
	return &authpb.AdminAccessResponse{
		Id:            result.ID,
		UserId:        result.UserID,
		IpAddress:     result.IPAddress,
		IpRange:       result.IPRange,
		DeviceId:      result.DeviceID,
		DeviceName:    result.DeviceName,
		DeviceType:    result.DeviceType,
		AuthType:      result.AuthType,
		Status:        result.Status,
		EffectiveFrom: result.EffectiveFrom.Format("2006-01-02T15:04:05Z07:00"),
		EffectiveTo:   result.EffectiveTo.Format("2006-01-02T15:04:05Z07:00"),
		MaxDevices:    int32(result.MaxDevices),
		Require2Fa:    result.Require2FA,
		CreatedBy:     result.CreatedBy,
		UpdatedBy:     result.UpdatedBy,
		CreatedAt:     result.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     result.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// @Summary Cập nhật quyền truy cập admin
// @Description Cập nhật quyền truy cập hệ thống nội bộ
// @Tags Admin Access Control
// @Accept json
// @Produce json
// @Param id path int true "ID quyền truy cập"
// @Param body body authpb.UpdateAdminAccessRequest true "Thông tin cập nhật"
// @Success 200 {object} authpb.AdminAccessResponse
// @Failure 400 {object} authpb.AdminAccessResponse
// @Failure 401 {object} authpb.AdminAccessResponse
// @Failure 404 {object} authpb.AdminAccessResponse
// @Failure 500 {object} authpb.AdminAccessResponse
// @Router /admin-access/{id} [put]
func (h *AdminAccessHandler) UpdateAdminAccess(ctx context.Context, req *authpb.UpdateAdminAccessRequest) (*authpb.AdminAccessResponse, error) {
	// Convert proto request to DTO
	var effectiveFrom, effectiveTo *time.Time
	if req.EffectiveFrom != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z07:00", req.EffectiveFrom); err == nil {
			effectiveFrom = &t
		}
	}
	if req.EffectiveTo != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z07:00", req.EffectiveTo); err == nil {
			effectiveTo = &t
		}
	}

	dto := &dto.AdminAccessRequest{
		UserID:        req.UserId,
		IPAddress:     req.IpAddress,
		IPRange:       req.IpRange,
		DeviceID:      req.DeviceId,
		DeviceName:    req.DeviceName,
		DeviceType:    req.DeviceType,
		AuthType:      req.AuthType,
		Status:        req.Status,
		EffectiveFrom: effectiveFrom,
		EffectiveTo:   effectiveTo,
		MaxDevices:    int(req.MaxDevices),
		Require2FA:    req.Require2Fa,
	}

	result, err := h.AdminAccessUsecase.UpdateAdminAccess(ctx, req.Id, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi cập nhật quyền truy cập: %v", err)
	}

	// Convert DTO response to proto response
	return &authpb.AdminAccessResponse{
		Id:            result.ID,
		UserId:        result.UserID,
		IpAddress:     result.IPAddress,
		IpRange:       result.IPRange,
		DeviceId:      result.DeviceID,
		DeviceName:    result.DeviceName,
		DeviceType:    result.DeviceType,
		AuthType:      result.AuthType,
		Status:        result.Status,
		EffectiveFrom: result.EffectiveFrom.Format("2006-01-02T15:04:05Z07:00"),
		EffectiveTo:   result.EffectiveTo.Format("2006-01-02T15:04:05Z07:00"),
		MaxDevices:    int32(result.MaxDevices),
		Require2Fa:    result.Require2FA,
		CreatedBy:     result.CreatedBy,
		UpdatedBy:     result.UpdatedBy,
		CreatedAt:     result.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     result.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// @Summary Xóa quyền truy cập admin
// @Description Xóa quyền truy cập hệ thống nội bộ
// @Tags Admin Access Control
// @Accept json
// @Produce json
// @Param id path int true "ID quyền truy cập"
// @Success 200 {object} authpb.DeleteAdminAccessResponse
// @Failure 400 {object} authpb.DeleteAdminAccessResponse
// @Failure 401 {object} authpb.DeleteAdminAccessResponse
// @Failure 404 {object} authpb.DeleteAdminAccessResponse
// @Failure 500 {object} authpb.DeleteAdminAccessResponse
// @Router /admin-access/{id} [delete]
func (h *AdminAccessHandler) DeleteAdminAccess(ctx context.Context, req *authpb.DeleteAdminAccessRequest) (*authpb.DeleteAdminAccessResponse, error) {
	err := h.AdminAccessUsecase.DeleteAdminAccess(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi xóa quyền truy cập: %v", err)
	}

	return &authpb.DeleteAdminAccessResponse{
		Message: "Quyền truy cập đã được xóa thành công",
	}, nil
}

// @Summary Lấy danh sách quyền truy cập admin
// @Description Lấy danh sách quyền truy cập hệ thống nội bộ
// @Tags Admin Access Control
// @Accept json
// @Produce json
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Param userId query int false "ID người dùng"
// @Param status query string false "Trạng thái"
// @Param authType query string false "Loại xác thực"
// @Param ipAddress query string false "Địa chỉ IP"
// @Success 200 {object} authpb.GetAdminAccessListResponse
// @Failure 400 {object} authpb.GetAdminAccessListResponse
// @Failure 401 {object} authpb.GetAdminAccessListResponse
// @Failure 500 {object} authpb.GetAdminAccessListResponse
// @Router /admin-access [get]
func (h *AdminAccessHandler) GetAdminAccessList(ctx context.Context, req *authpb.GetAdminAccessListRequest) (*authpb.GetAdminAccessListResponse, error) {
	// Convert proto request to DTO
	dto := &dto.AdminAccessListRequest{
		Page:      int(req.Page),
		Size:      int(req.Size),
		Status:    req.Status,
		AuthType:  req.AuthType,
		IPAddress: req.IpAddress,
	}

	if req.UserId > 0 {
		dto.UserID = &req.UserId
	}

	result, err := h.AdminAccessUsecase.GetAdminAccessList(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi lấy danh sách quyền truy cập: %v", err)
	}

	// Convert DTO response to proto response
	items := make([]*authpb.AdminAccessItem, len(result.Data))
	for i, item := range result.Data {
		items[i] = &authpb.AdminAccessItem{
			Id:            item.ID,
			UserId:        item.UserID,
			IpAddress:     item.IPAddress,
			IpRange:       item.IPRange,
			DeviceId:      item.DeviceID,
			DeviceName:    item.DeviceName,
			DeviceType:    item.DeviceType,
			AuthType:      item.AuthType,
			Status:        item.Status,
			EffectiveFrom: item.EffectiveFrom.Format("2006-01-02T15:04:05Z07:00"),
			EffectiveTo:   item.EffectiveTo.Format("2006-01-02T15:04:05Z07:00"),
			MaxDevices:    int32(item.MaxDevices),
			Require2Fa:    item.Require2FA,
			CreatedBy:     item.CreatedBy,
			UpdatedBy:     item.UpdatedBy,
			CreatedAt:     item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &authpb.GetAdminAccessListResponse{
		Data:  items,
		Total: int32(result.Total),
		Page:  int32(result.Page),
		Size:  int32(result.Size),
	}, nil
}

// @Summary Kiểm tra quyền truy cập admin
// @Description Kiểm tra quyền truy cập hệ thống nội bộ
// @Tags Admin Access Control
// @Accept json
// @Produce json
// @Param body body authpb.ValidateAdminAccessRequest true "Thông tin kiểm tra"
// @Success 200 {object} authpb.ValidateAdminAccessResponse
// @Failure 400 {object} authpb.ValidateAdminAccessResponse
// @Failure 500 {object} authpb.ValidateAdminAccessResponse
// @Router /admin-access/validate [post]
func (h *AdminAccessHandler) ValidateAdminAccess(ctx context.Context, req *authpb.ValidateAdminAccessRequest) (*authpb.ValidateAdminAccessResponse, error) {
	// Convert proto request to DTO
	dto := &dto.AdminAccessValidationRequest{
		UserID:    req.UserId,
		IPAddress: req.IpAddress,
		DeviceID:  req.DeviceId,
	}

	result, err := h.AdminAccessUsecase.ValidateAdminAccess(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi kiểm tra quyền truy cập: %v", err)
	}

	// Convert DTO response to proto response
	return &authpb.ValidateAdminAccessResponse{
		IsAllowed:  result.IsAllowed,
		Reason:     result.Reason,
		Require2Fa: result.Require2FA,
		AccessId:   result.AccessID,
		DeviceId:   result.DeviceID,
		IpAddress:  result.IPAddress,
	}, nil
}

// @Summary Tạo log truy cập admin
// @Description Tạo log truy cập hệ thống nội bộ
// @Tags Admin Access Control
// @Accept json
// @Produce json
// @Param body body authpb.CreateAdminAccessLogRequest true "Thông tin log"
// @Success 200 {object} authpb.CreateAdminAccessLogResponse
// @Failure 400 {object} authpb.CreateAdminAccessLogResponse
// @Failure 500 {object} authpb.CreateAdminAccessLogResponse
// @Router /admin-access/log [post]
func (h *AdminAccessHandler) CreateAdminAccessLog(ctx context.Context, req *authpb.CreateAdminAccessLogRequest) (*authpb.CreateAdminAccessLogResponse, error) {
	// Convert proto request to DTO
	dto := &dto.AdminAccessLogRequest{
		UserID:     req.UserId,
		IPAddress:  req.IpAddress,
		DeviceID:   req.DeviceId,
		DeviceName: req.DeviceName,
		Action:     req.Action,
		Status:     req.Status,
		Reason:     req.Reason,
		UserAgent:  req.UserAgent,
	}

	err := h.AdminAccessUsecase.CreateAdminAccessLog(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi tạo log truy cập: %v", err)
	}

	return &authpb.CreateAdminAccessLogResponse{
		Message: "Log truy cập đã được tạo thành công",
	}, nil
}

// @Summary Lấy danh sách log truy cập admin
// @Description Lấy danh sách log truy cập hệ thống nội bộ
// @Tags Admin Access Control
// @Accept json
// @Produce json
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Param userId query int false "ID người dùng"
// @Param ipAddress query string false "Địa chỉ IP"
// @Success 200 {object} authpb.GetAdminAccessLogListResponse
// @Failure 400 {object} authpb.GetAdminAccessLogListResponse
// @Failure 401 {object} authpb.GetAdminAccessLogListResponse
// @Failure 500 {object} authpb.GetAdminAccessLogListResponse
// @Router /admin-access/log [get]
func (h *AdminAccessHandler) GetAdminAccessLogList(ctx context.Context, req *authpb.GetAdminAccessLogListRequest) (*authpb.GetAdminAccessLogListResponse, error) {
	// Convert proto request to DTO
	dto := &dto.AdminAccessLogListRequest{
		Page:      int(req.Page),
		Size:      int(req.Size),
		UserID:    req.UserId,
		IPAddress: req.IpAddress,
	}

	result, err := h.AdminAccessUsecase.GetAdminAccessLogList(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi lấy danh sách log truy cập: %v", err)
	}

	// Convert DTO response to proto response
	items := make([]*authpb.AdminAccessLogItem, len(result.Data))
	for i, item := range result.Data {
		items[i] = &authpb.AdminAccessLogItem{
			Id:         item.ID,
			UserId:     item.UserID,
			IpAddress:  item.IPAddress,
			DeviceId:   item.DeviceID,
			DeviceName: item.DeviceName,
			Action:     item.Action,
			Status:     item.Status,
			Reason:     item.Reason,
			UserAgent:  item.UserAgent,
			CreatedAt:  item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &authpb.GetAdminAccessLogListResponse{
		Data:  items,
		Total: int32(result.Total),
		Page:  int32(result.Page),
		Size:  int32(result.Size),
	}, nil
}
