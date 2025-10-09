package handler

import (
	"context"
	authpb "pb/types/auth"
	userpb "pb/types/user"
	"user/infra/mapper"
	"user/internal/usecases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminHandler struct {
	userpb.UnimplementedAdminProfileServiceServer
	AdminUsecase *usecases.AdminUsecase
	AdminMapper  *mapper.AdminMapper
}

func NewAdminHandler(
	adminUsecase *usecases.AdminUsecase,
	adminMapper *mapper.AdminMapper) *AdminHandler {
	return &AdminHandler{
		AdminUsecase: adminUsecase,
		AdminMapper:  adminMapper,
	}
}

// @Summary Lấy danh sách tất cả user (Admin only)
// @Description Lấy danh sách tất cả user trong hệ thống với phân trang và tìm kiếm (chỉ dành cho admin)
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
func (h *AdminHandler) ListAdmin(ctx context.Context, req *userpb.ListAdminRequest) (*userpb.ListAdminResponse, error) {
	// TODO: Kiểm tra quyền admin
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if !h.isAdmin(profileID) {
	//     return nil, status.Errorf(codes.PermissionDenied, "unauthorized: admin access required")
	// }

	// Convert proto request to DTO using mapper
	dtoReq := h.AdminMapper.ListAdminRequestToDTO(req)

	// Call usecase
	result, total, err := h.AdminUsecase.ListAdmins(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}

	// Convert DTO response to proto response using mapper
	response := h.AdminMapper.ListAdminResponseToProto(result)

	return &userpb.ListAdminResponse{
		Data:  response,
		Total: uint32(total),
	}, nil
}

// @Summary Create Admin
// @Description Tạo tài khoản admin mới
// @Tags Admin
// @Accept json
// @Produce json
// @Param body body authpb.CreateAdminRequest true "Thông tin tạo admin"
// @Router /admin/create [post]
func (h *AdminHandler) CreateAdmin(ctx context.Context, req *authpb.CreateAdminRequest) (*authpb.AuthSuccessResponse, error) {
	// Convert proto request to model using mapper
	payload := h.AdminMapper.CreateAdminRequestToModel(req)

	result, err := h.AdminUsecase.CreateAdmin(ctx, payload)
	if err != nil {
		return nil, err
	}

	// Convert model response to proto response using mapper
	return h.AdminMapper.AdminProfileResponseToProto(result), nil
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
func (h *AdminHandler) UpdateAdmin(ctx context.Context, req *authpb.UpdateAdminRequest) (*authpb.AuthSuccessResponse, error) {
	// Convert proto request to DTO using mapper
	dtoReq := h.AdminMapper.UpdateAdminRequestToDTO(req)

	result, err := h.AdminUsecase.UpdateAdmin(ctx, *dtoReq)
	if err != nil {
		return nil, err
	}

	// Convert model response to proto response using mapper
	return h.AdminMapper.AdminProfileResponseToProto(result), nil
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
func (h *AdminHandler) DeleteAdmin(ctx context.Context, req *authpb.DeleteAdminRequest) (*authpb.DeleteAdminResponse, error) {
	// Convert proto request to DTO using mapper
	dtoReq := h.AdminMapper.DeleteAdminRequestToDTO(req)

	err := h.AdminUsecase.DeleteAdmin(ctx, *dtoReq)
	if err != nil {
		return nil, err
	}

	return &authpb.DeleteAdminResponse{
		Message: "Xóa tài khoản admin thành công",
	}, nil
}

// // @Summary Assign Role
// // @Description Gán quyền cho admin
// // @Tags Admin
// // @Accept json
// // @Produce json
// // @Param body body authpb.AssignRoleRequest true "Thông tin gán quyền"
// // @Success 200 {object} authpb.AssignRoleResponse
// // @Failure 400 {object} authpb.AssignRoleResponse
// // @Failure 500 {object} authpb.AssignRoleResponse
// // @Router /admin/assign-role [post]
// func (h *AdminHandler) AssignRole(ctx context.Context, req *authpb.AssignRoleRequest) (*authpb.AssignRoleResponse, error) {
// 	// Convert proto request to DTO using mapper
// 	dtoReq := h.AdminMapper.AssignRoleRequestToDTO(req)

// 	err := h.AdminUsecase.AssignRole(ctx, *dtoReq)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &authpb.AssignRoleResponse{
// 		Message: "Gán quyền thành công",
// 	}, nil
// }
