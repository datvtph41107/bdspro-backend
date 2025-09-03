package handler

import (
	"auth/internal/dto"
	"auth/internal/usecase"
	"context"
	authpb "pb/types/auth"
)

type AdminHandler struct {
	authpb.UnimplementedAuthServiceServer
	AdminUsecase    *usecase.AdminUsecase
	InternalHandler *InternalHandler
}

func NewAdminHandler(
	adminUsecase *usecase.AdminUsecase,
	internalHandler *InternalHandler,
) *AdminHandler {
	return &AdminHandler{
		AdminUsecase:    adminUsecase,
		InternalHandler: internalHandler,
	}
}

// @Summary Create Admin
// @Description Tạo tài khoản admin mới
// @Tags Admin
// @Accept json
// @Produce json
// @Param body body authpb.CreateAdminRequest true "Thông tin tạo admin"
// @Success 200 {object} authpb.AuthSuccessResponse
// @Failure 400 {object} authpb.AuthSuccessResponse
// @Failure 500 {object} authpb.AuthSuccessResponse
// @Router /admin/create [post]
func (h *AuthHandler) CreateAdmin(ctx context.Context, req *authpb.CreateAdminRequest) (*authpb.AuthSuccessResponse, error) {
	if err := h.InternalHandler.HasPermissions(ctx, []string{"ADMIN_TK_TAO"}); err != nil {
		return nil, err
	}

	// Convert proto request to DTO
	dto := &dto.CreateAdminRequest{
		Username: req.Username,
		Password: req.Password,
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
	}

	result, err := h.AdminUsecase.CreateAdmin(ctx, *dto)
	if err != nil {
		return nil, err
	}

	// Convert DTO response to proto response
	return &authpb.AuthSuccessResponse{
		AuthId:    result.AuthID,
		ProfileId: result.ProfileID,
		Role:      result.Role,
	}, nil
}

// @Summary Update Admin
// @Description Cập nhật thông tin admin
// @Tags Admin
// @Accept json
// @Produce json
// @Param body body authpb.UpdateAdminRequest true "Thông tin cập nhật admin"
// @Success 200 {object} authpb.AuthSuccessResponse
// @Failure 400 {object} authpb.AuthSuccessResponse
// @Failure 500 {object} authpb.AuthSuccessResponse
// @Router /admin/update [put]
func (h *AuthHandler) UpdateAdmin(ctx context.Context, req *authpb.UpdateAdminRequest) (*authpb.AuthSuccessResponse, error) {
	if err := h.InternalHandler.HasPermissions(ctx, []string{"ADMIN_TK_SUA"}); err != nil {
		return nil, err
	}

	// Convert proto request to DTO
	dto := &dto.UpdateAdminRequest{
		AuthID:   req.AuthId,
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
	}

	result, err := h.AdminUsecase.UpdateAdmin(ctx, *dto)
	if err != nil {
		return nil, err
	}

	// Convert DTO response to proto response
	return &authpb.AuthSuccessResponse{
		AuthId:    result.AuthID,
		ProfileId: result.ProfileID,
		Role:      result.Role,
	}, nil
}

// @Summary Delete Admin
// @Description Xóa tài khoản admin
// @Tags Admin
// @Accept json
// @Produce json
// @Param body body authpb.DeleteAdminRequest true "Thông tin xóa admin"
// @Success 200 {object} authpb.DeleteAdminResponse
// @Failure 400 {object} authpb.DeleteAdminResponse
// @Failure 500 {object} authpb.DeleteAdminResponse
// @Router /admin/delete [delete]
func (h *AuthHandler) DeleteAdmin(ctx context.Context, req *authpb.DeleteAdminRequest) (*authpb.DeleteAdminResponse, error) {
	if err := h.InternalHandler.HasPermissions(ctx, []string{"ADMIN_TK_XOA"}); err != nil {
		return nil, err
	}

	// Convert proto request to DTO
	dto := &dto.DeleteAdminRequest{
		AuthID: req.AuthId,
	}

	err := h.AdminUsecase.DeleteAdmin(ctx, *dto)
	if err != nil {
		return nil, err
	}

	return &authpb.DeleteAdminResponse{
		Message: "Xóa tài khoản admin thành công",
	}, nil
}

// @Summary Assign Role
// @Description Gán quyền cho admin
// @Tags Admin
// @Accept json
// @Produce json
// @Param body body authpb.AssignRoleRequest true "Thông tin gán quyền"
// @Success 200 {object} authpb.AssignRoleResponse
// @Failure 400 {object} authpb.AssignRoleResponse
// @Failure 500 {object} authpb.AssignRoleResponse
// @Router /admin/assign-role [post]
func (h *AuthHandler) AssignRole(ctx context.Context, req *authpb.AssignRoleRequest) (*authpb.AssignRoleResponse, error) {
	if err := h.InternalHandler.HasPermissions(ctx, []string{"ADMIN_TK_GAN"}); err != nil {
		return nil, err
	}

	// Convert proto request to DTO
	dto := &dto.AssignRoleRequest{
		AuthID:  req.AuthId,
		RoleKey: req.RoleKey,
	}

	err := h.AdminUsecase.AssignRole(ctx, *dto)
	if err != nil {
		return nil, err
	}

	return &authpb.AssignRoleResponse{
		Message: "Gán quyền thành công",
	}, nil
}

// @Summary List Admins
// @Description Lấy danh sách admin
// @Tags Admin
// @Accept json
// @Produce json
// @Param page query int false "Trang hiện tại"
// @Param size query int false "Số lượng mục trên mỗi trang"
// @Success 200 {object} authpb.AdminListResponse
// @Failure 400 {object} authpb.AdminListResponse
// @Failure 500 {object} authpb.AdminListResponse
// @Router /admin/list [get]
func (h *AuthHandler) ListAdmins(ctx context.Context, req *authpb.AdminListRequest) (*authpb.AdminListResponse, error) {
	if err := h.InternalHandler.HasPermissions(ctx, []string{"ADMIN_TK_XEM"}); err != nil {
		return nil, err
	}

	// Convert proto request to DTO
	dto := &dto.AdminListRequest{
		Page: req.Page,
		Size: req.Size,
	}

	result, err := h.AdminUsecase.ListAdmins(ctx, *dto)
	if err != nil {
		return nil, err
	}

	// Convert DTO response to proto response
	var adminItems []*authpb.AdminItem
	for _, item := range result.Data {
		adminItems = append(adminItems, &authpb.AdminItem{
			AuthId:    item.AuthID,
			ProfileId: item.ProfileID,
			Username:  item.Username,
			FullName:  item.FullName,
			Email:     item.Email,
			Phone:     item.Phone,
			Avatar:    item.Avatar,
			Role:      item.Role,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &authpb.AdminListResponse{
		Data:  adminItems,
		Total: result.Total,
	}, nil
}
