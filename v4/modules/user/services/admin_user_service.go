package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	_dto "common/domain/dto"
	_enum "common/domain/enum"
	_errors "common/errors"
	"common/identity"
	_utils "common/utils"
	authpb "pb/types/auth"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"userv4/infra/mapper"
	"userv4/internal/dto"
	"userv4/internal/interface/providers"
	"userv4/internal/job"
	usecase "userv4/internal/usecase/user"
	"userv4/internal/usecase/useradmin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AdminUserService owns admin-users + admin-profile HTTP surface.
type AdminUserService interface {
	ListAllUsers(ctx context.Context, req *userpb.ListAllUsersRequest) (*userpb.ListAllUsersResponse, error)
	LockUser(ctx context.Context, req *userpb.LockUserRequest) error
	UnLockUser(ctx context.Context, req *userpb.UnLockUserRequest) error
	ApproveUser(ctx context.Context, req *userpb.ApproveUserRequest) error
	GetUserStats(ctx context.Context, req *sharepb.GetStatsRequest) (*sharepb.GetStatsResponse, error)
	GetUserGrowthDailyStats(ctx context.Context, req *userpb.GetUserGrowthDailyStatsRequest) (*userpb.GetUserGrowthDailyStatsResponse, error)
	CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error)
	UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) error
	GetUserDetail(ctx context.Context, profileID uint64) (*userpb.GetUserDetailResponse, error)

	CreateAdmin(ctx context.Context, req *userpb.AdminProfile) (*userpb.AdminProfile, error)
	UpdateAdmin(ctx context.Context, req *userpb.AdminProfile) (*userpb.AdminProfile, error)
	DeleteAdmin(ctx context.Context, id uint64) (*userpb.DeleteAdminResponse, error)
	ListAdmins(ctx context.Context, req *userpb.AdminListRequest) (*userpb.AdminListResponse, error)
	GetAdminDetail(ctx context.Context, id uint64) (*userpb.AdminProfile, error)
}

// LocalAdminUserService adapts admin usecases into modules/user.
type LocalAdminUserService struct {
	admin       *usecase.AdminUsecase
	profile     *usecase.ProfileUsecase
	adminMapper *mapper.AdminMapper
	notify      providers.NotificationProvider
	statsJob    *job.UserDashboardStatsJob
	authorizer  useradmin.PermissionAuthorizer
}

func NewLocalAdminUserService(
	admin *usecase.AdminUsecase,
	profile *usecase.ProfileUsecase,
	adminMapper *mapper.AdminMapper,
	notify providers.NotificationProvider,
	statsJob *job.UserDashboardStatsJob,
	authorizer useradmin.PermissionAuthorizer,
) *LocalAdminUserService {
	return &LocalAdminUserService{
		admin:       admin,
		profile:     profile,
		adminMapper: adminMapper,
		notify:      notify,
		statsJob:    statsJob,
		authorizer:  authorizer,
	}
}

var _ AdminUserService = (*LocalAdminUserService)(nil)

func (s *LocalAdminUserService) requirePermission(ctx context.Context, permissionID uint32) error {
	_, err := useradmin.RequirePermission(ctx, s.authorizer, permissionID)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, useradmin.ErrActorRequired):
		return status.Error(codes.Unauthenticated, "trusted user admin actor is required")
	case errors.Is(err, useradmin.ErrPermissionDenied):
		return status.Errorf(codes.PermissionDenied, "permission %d is required", permissionID)
	case errors.Is(err, useradmin.ErrPermissionAuthorityUnavailable):
		return status.Error(codes.Unavailable, "user admin permission authority failed")
	default:
		return status.Error(codes.Internal, "user admin authorization failed")
	}
}

func (s *LocalAdminUserService) createAdminHistory(ctx context.Context, actionType _enum.EHistory, targetId uint64, title string, note []string) {
	if s.notify == nil {
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
	contextTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_ = s.notify.CreateHistory(contextTimeout, historyDTO)
}

func (s *LocalAdminUserService) ListAllUsers(ctx context.Context, req *userpb.ListAllUsersRequest) (*userpb.ListAllUsersResponse, error) {
	if err := s.requirePermission(ctx, useradmin.PermissionView); err != nil {
		return nil, err
	}
	sort := req.GetSort()
	if sort == "" && req.SortBy != "" {
		sortOrder := req.SortOrder
		if sortOrder == "" {
			sortOrder = "desc"
		}
		sort = req.SortBy + "," + sortOrder
	}
	dtoReq := &dto.ListUsersRequest{
		Pagable:   _dto.Pagable{Page: req.Page, Size: req.Size, Sort: sort},
		Search:    req.Search,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
		RoleType:  req.RoleType,
		RoleID:    &req.RoleId,
		Status:    &req.Status,
		Bookmark:  req.Bookmark,
	}
	if req.RoleId != 0 {
		roleId := req.RoleId
		dtoReq.RoleID = &roleId
	}
	if req.Status != 0 {
		st := req.Status
		dtoReq.Status = &st
	}

	result, total, err := s.profile.ListAllUsers(ctx, dtoReq)
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
	response := &userpb.ListAllUsersResponse{
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
		Total:      int32(total),
	}
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
			Position:     item.Position,
			DepartmentId: item.DepartmentID,
			Bookmark:     item.Bookmark,
			TickVerified: item.TickVerified,
		}
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

func (s *LocalAdminUserService) LockUser(ctx context.Context, req *userpb.LockUserRequest) error {
	if err := s.requirePermission(ctx, useradmin.PermissionManage); err != nil {
		return err
	}
	lockType := _enum.EUserStatus(req.LockType)
	if lockType != _enum.EUserStatusPermanentlyLocked && lockType != _enum.EUserStatusTemporaryLocked {
		return _errors.ReturnError(int32(400), "Loại khóa tài khoản không hợp lệ: 20(tạm thời) hoặc 30 (vĩnh viễn)")
	}
	if err := s.admin.LockUser(ctx, &dto.LockUserRequest{
		ProfileID: req.ProfileId,
		Reason:    req.Reason,
		LockType:  lockType,
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to block user: %v", err)
	}
	return nil
}

func (s *LocalAdminUserService) UnLockUser(ctx context.Context, req *userpb.UnLockUserRequest) error {
	if err := s.requirePermission(ctx, useradmin.PermissionManage); err != nil {
		return err
	}
	if err := s.admin.UnLockUser(ctx, &dto.UnLockUserRequest{
		ProfileID: req.ProfileId,
		Reason:    req.Reason,
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to unblock user: %v", err)
	}
	return nil
}

func (s *LocalAdminUserService) ApproveUser(ctx context.Context, req *userpb.ApproveUserRequest) error {
	if err := s.requirePermission(ctx, useradmin.PermissionManage); err != nil {
		return err
	}
	if err := s.admin.ApproveUser(ctx, &dto.ApproveUserRequest{
		ProfileID: req.ProfileId,
		Note:      req.Note,
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to approve user: %v", err)
	}
	return nil
}

func (s *LocalAdminUserService) GetUserStats(ctx context.Context, req *sharepb.GetStatsRequest) (*sharepb.GetStatsResponse, error) {
	if err := s.requirePermission(ctx, useradmin.PermissionView); err != nil {
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
	stats, stats2, err := s.statsJob.GetStats(ctx, firstDate, secondDate)
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

func (s *LocalAdminUserService) GetUserGrowthDailyStats(ctx context.Context, req *userpb.GetUserGrowthDailyStatsRequest) (*userpb.GetUserGrowthDailyStatsResponse, error) {
	if req == nil {
		req = &userpb.GetUserGrowthDailyStatsRequest{}
	}
	resp, err := s.admin.GetUserGrowthDailyStats(ctx, &dto.GetUserGrowthDailyStatsRequest{
		LastDates: req.GetLastDates(),
		Gap:       req.GetGap(),
	})
	if err != nil {
		return nil, err
	}
	data := make([]*userpb.UserGrowthDailyItem, 0, len(resp.Data))
	for _, item := range resp.Data {
		data = append(data, &userpb.UserGrowthDailyItem{
			Date:        item.Date,
			NewUsers:    item.NewUsers,
			ActiveUsers: item.ActiveUsers,
		})
	}
	return &userpb.GetUserGrowthDailyStatsResponse{
		Data:      data,
		LastDates: resp.LastDates,
		Gap:       resp.Gap,
	}, nil
}

func (s *LocalAdminUserService) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	if err := s.requirePermission(ctx, useradmin.PermissionManage); err != nil {
		return nil, err
	}
	var birthTime *time.Time
	if req.Birth != "" {
		birthTime = _utils.ParseStringToTime(req.Birth)
	}
	result, err := s.admin.CreateUser(ctx, &dto.CreateUserRequest{
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
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}
	s.createAdminHistory(ctx, _enum.HistoryUserCreate, result.ProfileID, "Tạo user mới", []string{
		"Username: " + result.Username,
		"Email: " + result.Email,
		"Phone: " + result.Phone,
	})
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

func (s *LocalAdminUserService) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	if err := s.requirePermission(ctx, useradmin.PermissionManage); err != nil {
		return nil, err
	}
	var birthTime *time.Time
	if req.Birth != "" {
		birthTime = _utils.ParseStringToTime(req.Birth)
	}
	var roleID *uint64
	if req.RoleId > 0 {
		roleID = &req.RoleId
	}
	var gender *uint32
	if req.Gender > 0 {
		g := uint32(req.Gender)
		gender = &g
	}
	result, err := s.admin.UpdateUser(ctx, &dto.UpdateUserRequest{
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
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}
	s.createAdminHistory(ctx, _enum.HistoryUserUpdate, result.ProfileID, "Cập nhật thông tin user", []string{
		"Email: " + result.Email,
		"Phone: " + result.Phone,
	})
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

func (s *LocalAdminUserService) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) error {
	if err := s.requirePermission(ctx, useradmin.PermissionManage); err != nil {
		return err
	}
	if err := s.admin.DeleteUser(ctx, &dto.DeleteUserRequest{
		ProfileID: req.ProfileId,
		Reason:    req.Reason,
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}
	s.createAdminHistory(ctx, _enum.HistoryUserDelete, req.ProfileId, "Xóa user", []string{
		"Lý do: " + req.Reason,
	})
	return nil
}

func (s *LocalAdminUserService) GetUserDetail(ctx context.Context, profileID uint64) (*userpb.GetUserDetailResponse, error) {
	if err := s.requirePermission(ctx, useradmin.PermissionView); err != nil {
		return nil, err
	}
	result, err := s.admin.GetUserDetail(ctx, profileID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user detail: %v", err)
	}
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
	if result.Birth != nil {
		response.Birth = _utils.FormatTimeToString(result.Birth)
	}
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

func (s *LocalAdminUserService) CreateAdmin(ctx context.Context, req *userpb.AdminProfile) (*userpb.AdminProfile, error) {
	if err := s.requirePermission(ctx, useradmin.PermissionManage); err != nil {
		return nil, err
	}
	payload := s.adminMapper.PbToDomain(req)
	result, err := s.admin.CreateAdmin(ctx, payload)
	if err != nil {
		return nil, err
	}
	s.createAdminHistory(ctx, _enum.HistoryAdminCreate, result.ID,
		"Tạo tài khoản admin mới",
		[]string{fmt.Sprintf("Tạo admin: %s (%s)", result.FullName, result.Email)})
	return s.adminMapper.AdminProfileResponseToProto(result), nil
}

func (s *LocalAdminUserService) UpdateAdmin(ctx context.Context, req *userpb.AdminProfile) (*userpb.AdminProfile, error) {
	if err := s.requirePermission(ctx, useradmin.PermissionManage); err != nil {
		return nil, err
	}
	dtoReq := s.adminMapper.PbToDomain(req)
	result, err := s.admin.UpdateAdmin(ctx, dtoReq)
	if err != nil {
		return nil, err
	}
	s.createAdminHistory(ctx, _enum.HistoryAdminUpdate, result.ID,
		"Cập nhật thông tin admin",
		[]string{fmt.Sprintf("Cập nhật admin: %s (%s)", result.FullName, result.Email)})
	return s.adminMapper.AdminProfileResponseToProto(result), nil
}

func (s *LocalAdminUserService) DeleteAdmin(ctx context.Context, id uint64) (*userpb.DeleteAdminResponse, error) {
	if err := s.requirePermission(ctx, useradmin.PermissionManage); err != nil {
		return nil, err
	}
	if err := s.admin.DeleteAdmin(ctx, id); err != nil {
		return nil, err
	}
	s.createAdminHistory(ctx, _enum.HistoryAdminDelete, id,
		"Xóa tài khoản admin",
		[]string{fmt.Sprintf("Xóa admin với ID: %d", id)})
	return &userpb.DeleteAdminResponse{Message: "Xóa tài khoản admin thành công"}, nil
}

func (s *LocalAdminUserService) ListAdmins(ctx context.Context, req *userpb.AdminListRequest) (*userpb.AdminListResponse, error) {
	if err := s.requirePermission(ctx, useradmin.PermissionView); err != nil {
		return nil, err
	}
	dtoReq := s.adminMapper.ListAdminRequestToDTO(req)
	result, total, err := s.admin.ListAdmins(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}
	return &userpb.AdminListResponse{
		Data:  s.adminMapper.ListAdminResponseToProto(result),
		Total: int32(total),
	}, nil
}

func (s *LocalAdminUserService) GetAdminDetail(ctx context.Context, id uint64) (*userpb.AdminProfile, error) {
	if err := s.requirePermission(ctx, useradmin.PermissionView); err != nil {
		return nil, err
	}
	result, err := s.admin.GetDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.adminMapper.AdminProfileResponseToProto(result), nil
}
