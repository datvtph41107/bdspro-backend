package handler

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	_errors "common/errors"
	"common/identity"
	_utils "common/utils"
	"context"
	authpb "pb/types/auth"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"time"
	"user/infra/mapper"
	"user/internal/dto"
	"user/internal/interface/providers"
	"user/internal/job"
	"user/internal/usecase/useradmin"
	"user/internal/usecases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminUserProfileHandler struct {
	userpb.UnimplementedAdminUserProfileServiceServer
	AdminUsecase          *usecases.AdminUsecase
	ProfileUsecase        *usecases.ProfileUsecase
	NotificationClient    providers.NotificationProvider
	UserDashboardStatsJob *job.UserDashboardStatsJob
	Authorizer            useradmin.PermissionAuthorizer
}

func NewAdminUserProfileHandler(
	adminUsecase *usecases.AdminUsecase,
	profileUsecase *usecases.ProfileUsecase,
	notificationClient providers.NotificationProvider,
	UserDashboardStatsJob *job.UserDashboardStatsJob,
	authorizer useradmin.PermissionAuthorizer,
) *AdminUserProfileHandler {
	return &AdminUserProfileHandler{
		AdminUsecase:          adminUsecase,
		ProfileUsecase:        profileUsecase,
		NotificationClient:    notificationClient,
		UserDashboardStatsJob: UserDashboardStatsJob,
		Authorizer:            authorizer,
	}
}

// createAdminHistory tạo lịch sử hoạt động cho admin
func (h *AdminUserProfileHandler) createAdminHistory(ctx context.Context, actionType _enum.EHistory, targetId uint64, title string, note []string) {
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
	}

	// Tạo context với timeout
	contextTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Gọi notification service để tạo history (bỏ qua lỗi để không ảnh hưởng đến API chính)
	_ = h.NotificationClient.CreateHistory(contextTimeout, historyDTO)
}

// @Summary Lấy danh sách tất cả user (Admin only)
// @Description Lấy danh sách tất cả user trong hệ thống với phân trang và tìm kiếm (chỉ dành cho admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param search query string false "Tìm kiếm theo tên, email, phone"
// @Param page query int false "Trang hiện tại (mặc định: 1)"
// @Param size query int false "Số lượng item trên mỗi trang (mặc định: 20, tối đa: 100)"
// @Param sort query string false "Sắp xếp (format: createdAt,desc hoặc lastLogin,desc)"
// @Param sortBy query string false "Sắp xếp theo trường (createdAt, fullName, email, lastLogin)"
// @Param sortOrder query string false "Thứ tự sắp xếp (asc, desc)"
// @Param status query int false "Lọc theo trạng thái (1: active, 0: inactive)"
// @Param roleId query int false "Lọc theo role ID"
// @Param roleType query string false "Lọc theo role type"
// @Param startDate query string false "Lọc từ ngày tạo (format: 2006-01-02T15:04:05)"
// @Param endDate query string false "Lọc đến ngày tạo (format: 2006-01-02T15:04:05)"
// @Success 200 {object} userpb.ListAllUsersResponse
// @Failure 400 {object} userpb.ListAllUsersResponse
// @Failure 401 {object} userpb.ListAllUsersResponse
// @Failure 500 {object} userpb.ListAllUsersResponse
// @Router /v2/user/admin/users [get]
func (s *AdminUserProfileHandler) ListAllUsers(ctx context.Context, req *userpb.ListAllUsersRequest) (*userpb.ListAllUsersResponse, error) {
	if err := requireUserAdminPermission(ctx, s.Authorizer, useradmin.PermissionView); err != nil {
		return nil, err
	}
	// Convert proto request to DTO
	sort := req.GetSort()
	if sort == "" && req.SortBy != "" {
		sortOrder := req.SortOrder
		if sortOrder == "" {
			sortOrder = "desc"
		}
		sort = req.SortBy + "," + sortOrder
	}

	dtoReq := &dto.ListUsersRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
			Sort: sort,
		},
		Search:    req.Search,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
		RoleType:  req.RoleType,
		RoleID:    &req.RoleId,
		Status:    &req.Status,
		Bookmark:  req.Bookmark,
	}

	// Set roleId filter
	if req.RoleId != 0 {
		roleId := req.RoleId
		dtoReq.RoleID = &roleId
	}

	// Set bookmark filter
	if req.Bookmark != nil {
		dtoReq.Bookmark = req.Bookmark
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
	result, total, err := s.ProfileUsecase.ListAllUsers(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	size := req.Size
	if size < 1 {
		size = 20
	}
	var totalPages uint32
	if size > 0 && total > 0 {
		totalPages = (total + size - 1) / size
	}

	// Convert DTO response to proto response
	response := &userpb.ListAllUsersResponse{
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
		Total:      int32(total),
	}

	// Convert data items
	kycMapper := mapper.NewKYCMapper()
	for _, item := range result {
		protoItem := &userpb.AdminUserItem{
			ProfileId:    item.ProfileID,
			FullName:     item.FullName,
			Email:        item.Email,
			Phone:        item.Phone,
			Avatar:       item.Avatar,
			Address:      item.Address,
			Gender:       uint32(item.Gender),
			Status:       uint32(item.Status),
			WarningLevel: uint32(item.WarningLevel),
			// RoleType:  item.RoleType,
			// RoleRealEstate: item.RoleRealEstate,
			Position:     item.Position,
			DepartmentId: item.DepartmentID,
			Bookmark:     item.Bookmark,
			TickVerified: item.TickVerified,
		}

		// Set role information if available
		if item.RoleName != "" {
			roleID := uint64(0)
			if item.RoleID != nil {
				roleID = *item.RoleID
			}
			protoItem.Role = &authpb.Role{
				Id:       roleID,
				RoleName: item.RoleName,
				RoleKey:  item.RoleKeyData,
			}

			if item.RoleColor != "" && item.RoleBgColor != "" {
				protoItem.Role.Color = &sharepb.Color{
					ContentColor:    item.RoleColor,
					BackgroundColor: item.RoleBgColor,
				}
			}
		}

		// Convert time fields
		protoItem.CreatedAt = _utils.FormatTimeToString(item.CreatedAt)
		protoItem.UpdatedAt = _utils.FormatTimeToString(item.UpdatedAt)
		if item.Birth != nil {
			protoItem.Birth = _utils.FormatTimeToString(item.Birth)
		}
		if item.LastLoginAt != nil {
			protoItem.LastLoginAt = _utils.FormatTimeToString(item.LastLoginAt)
		}
		protoItem.TotalLogin = item.TotalLogin
		if item.KYC != nil {
			protoItem.Kyc = kycMapper.ToProtoWithOwnerCheck(item.KYC, true)
		}

		response.Data = append(response.Data, protoItem)
	}

	return response, nil
}

// @Summary Khóa tài khoản user (Admin only)
// @Description Khóa tài khoản user với lý do cụ thể. Cập nhật status thành "Khóa tạm thời" trong bảng user_profile (chỉ dành cho admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param profileId path int true "ID của user cần khóa"
// @Param request body userpb.LockUserRequest true "Thông tin khóa tài khoản"
// @Success 200 {object} sharepb.Empty
// @Failure 400 {object} sharepb.Empty
// @Failure 401 {object} sharepb.Empty
// @Failure 500 {object} sharepb.Empty
// @Router /v2/user/admin/users/{profileId}/lock [post]
func (s *AdminUserProfileHandler) LockUser(ctx context.Context, req *userpb.LockUserRequest) (*sharepb.Empty, error) {
	if err := requireUserAdminPermission(ctx, s.Authorizer, useradmin.PermissionManage); err != nil {
		return nil, err
	}

	lockType := _enum.EUserStatus(req.LockType)
	if lockType != _enum.EUserStatusPermanentlyLocked && lockType != _enum.EUserStatusTemporaryLocked {
		return nil, _errors.ReturnError(int32(400), "Loại khóa tài khoản không hợp lệ: 20(tạm thời) hoặc 30 (vĩnh viễn)")
	}
	// Convert proto request to DTO
	dtoReq := &dto.LockUserRequest{
		ProfileID: req.ProfileId,
		Reason:    req.Reason,
		LockType:  lockType,
	}

	// Call usecase
	err := s.AdminUsecase.LockUser(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to block user: %v", err)
	}

	return &sharepb.Empty{}, nil
}

// @Summary Mở khóa tài khoản user (Admin only)
// @Description Mở khóa tài khoản user với lý do cụ thể. Cập nhật status thành "Hoạt động" trong bảng user_profile (chỉ dành cho admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param profileId path int true "ID của user cần mở khóa"
// @Param request body userpb.UnLockUserRequest true "Thông tin mở khóa tài khoản"
// @Success 200 {object} sharepb.Empty
// @Failure 400 {object} sharepb.Empty
// @Failure 401 {object} sharepb.Empty
// @Failure 500 {object} sharepb.Empty
// @Router /v2/user/admin/users/{profileId}/unlock [post]
func (s *AdminUserProfileHandler) UnLockUser(ctx context.Context, req *userpb.UnLockUserRequest) (*sharepb.Empty, error) {
	if err := requireUserAdminPermission(ctx, s.Authorizer, useradmin.PermissionManage); err != nil {
		return nil, err
	}

	// Convert proto request to DTO
	dtoReq := &dto.UnLockUserRequest{
		ProfileID: req.ProfileId,
		Reason:    req.Reason,
	}

	// Call usecase
	err := s.AdminUsecase.UnLockUser(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unblock user: %v", err)
	}

	return &sharepb.Empty{}, nil
}

// @Summary Duyệt tài khoản user (Admin only)
// @Description Duyệt tài khoản user với ghi chú. Cập nhật status thành "Hoạt động" trong bảng user_profile (chỉ dành cho admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param profileId path int true "ID của user cần duyệt"
// @Param request body userpb.ApproveUserRequest true "Thông tin duyệt tài khoản"
// @Success 200 {object} sharepb.Empty
// @Failure 400 {object} sharepb.Empty
// @Failure 401 {object} sharepb.Empty
// @Failure 500 {object} sharepb.Empty
// @Router /v2/user/admin/users/{profileId}/approve [post]
func (s *AdminUserProfileHandler) ApproveUser(ctx context.Context, req *userpb.ApproveUserRequest) (*sharepb.Empty, error) {
	if err := requireUserAdminPermission(ctx, s.Authorizer, useradmin.PermissionManage); err != nil {
		return nil, err
	}

	// Convert proto request to DTO
	dtoReq := &dto.ApproveUserRequest{
		ProfileID: req.ProfileId,
		Note:      req.Note,
	}

	// Call usecase
	err := s.AdminUsecase.ApproveUser(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to approve user: %v", err)
	}

	return &sharepb.Empty{}, nil
}

// @Summary Lấy thống kê số lượng user
// @Description Lấy thống kê số lượng user theo 2 ngày
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param firstDate query string true "Ngày đầu tiên (YYYY-MM-DD)"
// @Param secondDate query string true "Ngày thứ hai (YYYY-MM-DD)"
// @Security BearerAuth
// @Success 200 {object} sharepb.GetStatsResponse "Thành công"
// @Router /v2/user/stats [get]
func (s *AdminUserProfileHandler) GetUserStats(ctx context.Context, req *sharepb.GetStatsRequest) (*sharepb.GetStatsResponse, error) {
	if err := requireUserAdminPermission(ctx, s.Authorizer, useradmin.PermissionView); err != nil {
		return nil, err
	}
	firstDate, err := time.Parse(time.DateOnly, req.FirstDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid firstDate format")
	}

	secondDate, err := time.Parse(time.DateOnly, req.SecondDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid secondDate format")
	}

	stats, stats2, err := s.UserDashboardStatsJob.GetStats(ctx, firstDate, secondDate)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &sharepb.GetStatsResponse{
		FirstDate:       stats.CalculateTime.Format(time.DateOnly),
		SecondDate:      stats2.CalculateTime.Format(time.DateOnly),
		FirstDateCount:  uint32(stats.Count),
		SecondDateCount: uint32(stats2.Count),
	}, nil
}

// @Summary Tạo user mới (Admin only)
// @Description Tạo user mới giống với tài khoản admin (chỉ dành cho admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body userpb.CreateUserRequest true "Thông tin user mới"
// @Success 200 {object} userpb.CreateUserResponse
// @Failure 400 {object} userpb.CreateUserResponse
// @Failure 401 {object} userpb.CreateUserResponse
// @Failure 500 {object} userpb.CreateUserResponse
// @Router /v2/user/admin/user/new [post]
func (h *AdminUserProfileHandler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	if err := requireUserAdminPermission(ctx, h.Authorizer, useradmin.PermissionManage); err != nil {
		return nil, err
	}
	// Parse birth date if provided
	var birthTime *time.Time
	if req.Birth != "" {
		parsedBirth := _utils.ParseStringToTime(req.Birth)
		birthTime = parsedBirth
	}

	// Convert proto request to DTO
	dtoReq := &dto.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		RoleID:   req.RoleId,
		Gender:   req.Gender,
		Address:  req.Address,
		Birth:    birthTime,
	}

	// Call usecase
	result, err := h.AdminUsecase.CreateUser(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	// Create history log
	h.createAdminHistory(ctx, _enum.HistoryUserCreate, result.ProfileID, "Tạo user mới", []string{
		"Username: " + result.Username,
		"Email: " + result.Email,
		"Phone: " + result.Phone,
	})

	// Convert DTO response to proto response
	return &userpb.CreateUserResponse{
		ProfileId: result.ProfileID,
		AuthId:    result.AuthID,
		Username:  result.Username,
		FullName:  result.FullName,
		Email:     result.Email,
		Phone:     result.Phone,
		Avatar:    result.Avatar,
		RoleId:    result.RoleID,
		CreatedAt: _utils.FormatTimeToString(&result.CreatedAt),
	}, nil
}

// @Summary Cập nhật thông tin user (Admin only)
// @Description Cập nhật thông tin user (chỉ dành cho admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param profileId path int true "ID của user cần cập nhật"
// @Param request body userpb.UpdateUserRequest true "Thông tin cập nhật"
// @Success 200 {object} userpb.UpdateUserResponse
// @Failure 400 {object} userpb.UpdateUserResponse
// @Failure 401 {object} userpb.UpdateUserResponse
// @Failure 404 {object} userpb.UpdateUserResponse
// @Failure 500 {object} userpb.UpdateUserResponse
// @Router /v2/user/admin/users/{profileId} [put]
func (h *AdminUserProfileHandler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	if err := requireUserAdminPermission(ctx, h.Authorizer, useradmin.PermissionManage); err != nil {
		return nil, err
	}
	// Parse birth date if provided
	var birthTime *time.Time
	if req.Birth != "" {
		parsedBirth := _utils.ParseStringToTime(req.Birth)
		birthTime = parsedBirth
	}

	// Convert proto request to DTO
	var roleID *uint64
	if req.RoleId > 0 {
		roleID = &req.RoleId
	}

	var gender *uint32
	if req.Gender > 0 {
		g := uint32(req.Gender)
		gender = &g
	}

	dtoReq := &dto.UpdateUserRequest{
		ProfileID: req.ProfileId,
		FullName:  req.FullName,
		Email:     req.Email,
		Phone:     req.Phone,
		Avatar:    req.Avatar,
		RoleID:    roleID,
		Gender:    gender,
		Address:   req.Address,
		Birth:     birthTime,
		Status:    _enum.EUserStatus(req.Status),
	}

	// Call usecase
	result, err := h.AdminUsecase.UpdateUser(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	// Create history log
	h.createAdminHistory(ctx, _enum.HistoryUserUpdate, result.ProfileID, "Cập nhật thông tin user", []string{
		"Email: " + result.Email,
		"Phone: " + result.Phone,
	})

	// Convert DTO response to proto response
	return &userpb.UpdateUserResponse{
		ProfileId: result.ProfileID,
		FullName:  result.FullName,
		Email:     result.Email,
		Phone:     result.Phone,
		Avatar:    result.Avatar,
		RoleId:    result.RoleID,
		Status:    uint32(result.Status),
		UpdatedAt: _utils.FormatTimeToString(result.UpdatedAt),
	}, nil
}

// @Summary Xóa user (Admin only)
// @Description Xóa user khỏi hệ thống (chỉ dành cho admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param profileId path int true "ID của user cần xóa"
// @Param request body userpb.DeleteUserRequest true "Thông tin xóa"
// @Success 200 {object} sharepb.Empty
// @Failure 400 {object} sharepb.Empty
// @Failure 401 {object} sharepb.Empty
// @Failure 404 {object} sharepb.Empty
// @Failure 500 {object} sharepb.Empty
// @Router /v2/user/admin/users/{profileId} [delete]
func (h *AdminUserProfileHandler) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*sharepb.Empty, error) {
	if err := requireUserAdminPermission(ctx, h.Authorizer, useradmin.PermissionManage); err != nil {
		return nil, err
	}
	// Convert proto request to DTO
	dtoReq := &dto.DeleteUserRequest{
		ProfileID: req.ProfileId,
		Reason:    req.Reason,
	}

	// Call usecase
	err := h.AdminUsecase.DeleteUser(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	// Create history log
	h.createAdminHistory(ctx, _enum.HistoryUserDelete, req.ProfileId, "Xóa user", []string{
		"Lý do: " + req.Reason,
	})

	return &sharepb.Empty{}, nil
}

// @Summary Lấy chi tiết user (Admin only)
// @Description Lấy chi tiết thông tin user (chỉ dành cho admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param profileId path int true "ID của user cần lấy chi tiết"
// @Success 200 {object} userpb.GetUserDetailResponse
// @Failure 400 {object} userpb.GetUserDetailResponse
// @Failure 401 {object} userpb.GetUserDetailResponse
// @Failure 404 {object} userpb.GetUserDetailResponse
// @Failure 500 {object} userpb.GetUserDetailResponse
// @Router /v2/user/admin/users/{profileId} [get]
func (h *AdminUserProfileHandler) GetUserDetail(ctx context.Context, req *userpb.GetUserDetailRequest) (*userpb.GetUserDetailResponse, error) {
	if err := requireUserAdminPermission(ctx, h.Authorizer, useradmin.PermissionView); err != nil {
		return nil, err
	}
	// Call usecase
	result, err := h.AdminUsecase.GetUserDetail(ctx, req.ProfileId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user detail: %v", err)
	}

	// Convert DTO response to proto response
	response := &userpb.GetUserDetailResponse{
		ProfileId:       result.ProfileID,
		FullName:        result.FullName,
		Email:           result.Email,
		Phone:           result.Phone,
		Phone2:          result.Phone2,
		Avatar:          result.Avatar,
		Address:         result.Address,
		Gender:          uint32(result.Gender),
		CreatedAt:       _utils.FormatTimeToString(&result.CreatedAt),
		UpdatedAt:       _utils.FormatTimeToString(&result.UpdatedAt),
		Status:          uint32(result.Status),
		Position:        result.Position,
		DepartmentId:    result.DepartmentID,
		TickVerified:    result.TickVerified,
		BackgroundImage: result.BackgroundImage,
		TaxCode:         result.TaxCode,
		FrontIdentify:   result.FrontIdentify,
		BackIdentify:    result.BackIdentify,
		Slogan:          result.Slogan,
		Website:         result.Website,
		Facebook:        result.Facebook,
		Instagram:       result.Instagram,
		Twitter:         result.Twitter,
		Linkedin:        result.Linkedin,
		Youtube:         result.Youtube,
		Introduction:    result.Introduction,
	}

	// Set birth if available
	if result.Birth != nil {
		response.Birth = _utils.FormatTimeToString(result.Birth)
	}

	// Set role information if available
	if result.RoleName != "" {
		roleID := uint64(0)
		if result.RoleID != nil {
			roleID = *result.RoleID
		}
		response.Role = &authpb.Role{
			Id:       roleID,
			RoleName: result.RoleName,
			RoleKey:  result.RoleKey,
		}

		if result.RoleColor != "" && result.RoleBgColor != "" {
			response.Role.Color = &sharepb.Color{
				ContentColor:    result.RoleColor,
				BackgroundColor: result.RoleBgColor,
			}
		}
	}

	return response, nil
}
