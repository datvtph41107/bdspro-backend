package handler

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	"common/identity"
	"context"
	"fmt"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"time"
	"user/infra/mapper"
	"user/internal/interface/providers"
	"user/internal/usecase/useradmin"
	"user/internal/usecases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminHandler struct {
	userpb.UnimplementedAdminProfileServiceServer
	AdminUsecase       *usecases.AdminUsecase
	AdminMapper        *mapper.AdminMapper
	NotificationClient providers.NotificationProvider
	Authorizer         useradmin.PermissionAuthorizer
}

func NewAdminHandler(
	adminUsecase *usecases.AdminUsecase,
	adminMapper *mapper.AdminMapper,
	notificationClient providers.NotificationProvider,
	authorizer useradmin.PermissionAuthorizer,
) *AdminHandler {
	return &AdminHandler{
		AdminUsecase:       adminUsecase,
		AdminMapper:        adminMapper,
		NotificationClient: notificationClient,
		Authorizer:         authorizer,
	}
}

// createAdminHistory tạo lịch sử hoạt động cho admin
func (h *AdminHandler) createAdminHistory(ctx context.Context, actionType _enum.EHistory, targetId uint64, title string, note []string) {
	if h.NotificationClient == nil {
		return
	}

	actor, ok := identity.ActorFromContext(ctx)
	if !ok || actor.ProfileID == 0 {
		return
	}
	profileID := actor.ProfileID

	historyDTO := &_dto.HistoryDTO{
		Title:      title,
		Note:       note,
		TargetId:   targetId,
		TargetType: _enum.TargetHistoryAdmin,
		ActionType: actionType,
		OwnerID:    &profileID,
		// OwnerType:  int32(_enum.EOwnerOfMember),
	}

	// Tạo context với timeout
	contextTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Gọi notification service để tạo history (bỏ qua lỗi để không ảnh hưởng đến API chính)
	_ = h.NotificationClient.CreateHistory(contextTimeout, historyDTO)
}

// @Summary Lấy danh sách tất cả admin (Admin only)
// @Description Lấy danh sách tất cả admin trong hệ thống với phân trang và thông tin role (chỉ dành cho admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param search query string false "Tìm kiếm theo tên, email, phone"
// @Param page query int false "Trang hiện tại (mặc định: 1)"
// @Param size query int false "Số lượng item trên mỗi trang (mặc định: 20, tối đa: 100)"
// @Param sortBy query string false "Sắp xếp theo trường (createdAt, fullName, email)"
// @Param sortOrder query string false "Thứ tự sắp xếp (asc, desc)"
// @Param status query int false "Lọc theo trạng thái (1: active, 0: inactive)"
// @Param roleType query string false "Lọc theo role type"
// @Param startDate query string false "Lọc từ ngày tạo (format: 2006-01-02T15:04:05)"
// @Param endDate query string false "Lọc đến ngày tạo (format: 2006-01-02T15:04:05)"
// @Router /admin/list [get]
func (h *AdminHandler) ListAdmins(ctx context.Context, req *userpb.AdminListRequest) (*userpb.AdminListResponse, error) {
	if err := requireUserAdminPermission(ctx, h.Authorizer, useradmin.PermissionView); err != nil {
		return nil, err
	}

	// Convert proto request to DTO using mapper
	dtoReq := h.AdminMapper.ListAdminRequestToDTO(req)

	// Call usecase
	result, total, err := h.AdminUsecase.ListAdmins(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}

	// Convert DTO response to proto response using mapper (bao gồm thông tin role)
	response := h.AdminMapper.ListAdminResponseToProto(result)

	return &userpb.AdminListResponse{
		Data:  response,
		Total: int32(total),
	}, nil
}

// @Summary Create Admin
// @Description Tạo tài khoản admin mới
// @Tags Admin
// @Accept json
// @Produce json
// @Param body body userpb.AdminProfile true "Thông tin tạo admin"
// @Router /admin/create [post]
func (h *AdminHandler) CreateAdmin(ctx context.Context, req *userpb.AdminProfile) (*userpb.AdminProfile, error) {
	if err := requireUserAdminPermission(ctx, h.Authorizer, useradmin.PermissionManage); err != nil {
		return nil, err
	}
	// Convert proto request to model using mapper
	payload := h.AdminMapper.PbToDomain(req)

	result, err := h.AdminUsecase.CreateAdmin(ctx, payload)
	if err != nil {
		return nil, err
	}

	// Ghi nhật ký tạo admin
	h.createAdminHistory(ctx, _enum.HistoryAdminCreate, result.ID,
		"Tạo tài khoản admin mới",
		[]string{fmt.Sprintf("Tạo admin: %s (%s)", result.FullName, result.Email)})

	// Convert model response to proto response using mapper
	return h.AdminMapper.AdminProfileResponseToProto(result), nil
}

// @Summary Update Admin
// @Description Cập nhật thông tin admin
// @Tags Admin
// @Accept json
// @Produce json
// @Param body body userpb.AdminProfile true "Thông tin cập nhật admin"
// @Success 200 {object} userpb.AdminProfile
// @Failure 400 {object} userpb.AdminProfile
// @Failure 500 {object} authpb.AuthSuccessResponse
// @Router /admin/update [put]
func (h *AdminHandler) UpdateAdmin(ctx context.Context, req *userpb.AdminProfile) (*userpb.AdminProfile, error) {
	if err := requireUserAdminPermission(ctx, h.Authorizer, useradmin.PermissionManage); err != nil {
		return nil, err
	}
	// Convert proto request to DTO using mapper
	dtoReq := h.AdminMapper.PbToDomain(req)

	result, err := h.AdminUsecase.UpdateAdmin(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	// Ghi nhật ký cập nhật admin
	h.createAdminHistory(ctx, _enum.HistoryAdminUpdate, result.ID,
		"Cập nhật thông tin admin",
		[]string{fmt.Sprintf("Cập nhật admin: %s (%s)", result.FullName, result.Email)})

	// Convert model response to proto response using mapper
	return h.AdminMapper.AdminProfileResponseToProto(result), nil
}

// @Summary Delete Admin
// @Description Xóa tài khoản admin
// @Tags Admin
// @Accept json
// @Produce json
// @Param body body sharepb.IdRequest true "Thông tin xóa admin"
// @Router /admin/delete [delete]
func (h *AdminHandler) DeleteAdmin(ctx context.Context, req *sharepb.IdRequest) (*userpb.DeleteAdminResponse, error) {
	if err := requireUserAdminPermission(ctx, h.Authorizer, useradmin.PermissionManage); err != nil {
		return nil, err
	}
	err := h.AdminUsecase.DeleteAdmin(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Ghi nhật ký xóa admin
	h.createAdminHistory(ctx, _enum.HistoryAdminDelete, req.Id,
		"Xóa tài khoản admin",
		[]string{fmt.Sprintf("Xóa admin với ID: %d", req.Id)})

	return &userpb.DeleteAdminResponse{
		Message: "Xóa tài khoản admin thành công",
	}, nil
}

func (h *AdminHandler) GetDetail(ctx context.Context, req *sharepb.IdRequest) (*userpb.AdminProfile, error) {
	if err := requireUserAdminPermission(ctx, h.Authorizer, useradmin.PermissionView); err != nil {
		return nil, err
	}
	// Convert proto request to DTO using mapper

	result, err := h.AdminUsecase.GetDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return h.AdminMapper.AdminProfileResponseToProto(result), nil
}
