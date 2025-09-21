package handler

import (
	_utils "common/utils"
	"context"
	userpb "pb/types/user"
	"user/internal/dto"
	"user/internal/usecases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminHandler struct {
	userpb.UnimplementedProfileServiceServer
	AdminUsecase *usecases.AdminUsecase
}

func NewAdminHandler(adminUsecase *usecases.AdminUsecase) *AdminHandler {
	return &AdminHandler{
		AdminUsecase: adminUsecase,
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
// @Success 200 {object} userpb.ListAllUsersResponse
// @Failure 400 {object} userpb.ListAllUsersResponse
// @Failure 401 {object} userpb.ListAllUsersResponse
// @Failure 500 {object} userpb.ListAllUsersResponse
// @Router /v2/user/admin/users [get]
func (h *AdminHandler) ListAllUsers(ctx context.Context, req *userpb.ListAllUsersRequest) (*userpb.ListAllUsersResponse, error) {
	// TODO: Kiểm tra quyền admin
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if !h.isAdmin(profileID) {
	//     return nil, status.Errorf(codes.PermissionDenied, "unauthorized: admin access required")
	// }

	// Convert proto request to DTO
	dtoReq := &dto.ListAllUsersRequest{
		Search:    req.Search,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
		RoleType:  req.RoleType,
	}

	// Set pagination
	if req.Page > 0 {
		dtoReq.Page = req.Page
	} else {
		dtoReq.Page = 1
	}
	if req.Size > 0 {
		dtoReq.Size = req.Size
	} else {
		dtoReq.Size = 20
	}

	// Set status filter
	if req.Status != 0 {
		status := req.Status
		dtoReq.Status = &status
	}

	// Set date filters
	if req.StartDate != "" {
		// TODO: Parse startDate string to time.Time
	}
	if req.EndDate != "" {
		// TODO: Parse endDate string to time.Time
	}

	// Call usecase
	result, err := h.AdminUsecase.ListAllUsers(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}

	// Convert DTO response to proto response
	response := &userpb.ListAllUsersResponse{
		Total:      result.Total,
		Page:       result.Page,
		Size:       result.Size,
		TotalPages: result.TotalPages,
	}

	// Convert data items
	for _, item := range result.Data {
		protoItem := &userpb.AdminUserItem{
			ProfileId:      item.ProfileID,
			FullName:       item.FullName,
			Email:          item.Email,
			Phone:          item.Phone,
			Avatar:         item.Avatar,
			Address:        item.Address,
			Gender:         uint32(item.Gender),
			Status:         uint32(item.Status),
			RoleType:       item.RoleType,
			RoleRealEstate: item.RoleRealEstate,
			Position:       item.Position,
			DepartmentId:   item.DepartmentID,
			TickVerified:   item.TickVerified,
			TotalLogin:     item.TotalLogin,
		}

		// Convert time fields
		protoItem.CreatedAt = _utils.FormatTimeToString(&item.CreatedAt)
		protoItem.UpdatedAt = _utils.FormatTimeToString(&item.UpdatedAt)
		if item.Birth != nil {
			protoItem.Birth = _utils.FormatTimeToString(item.Birth)
		}
		if item.LastLoginAt != nil {
			protoItem.LastLoginAt = _utils.FormatTimeToString(item.LastLoginAt)
		}
		if item.VerifiedAt != nil {
			protoItem.VerifiedAt = _utils.FormatTimeToString(item.VerifiedAt)
		}
		if item.LockedAt != nil {
			protoItem.LockedAt = _utils.FormatTimeToString(item.LockedAt)
		}

		response.Data = append(response.Data, protoItem)
	}

	return response, nil
}
