package usecase

import (
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"user/internal/domain/access"
	"user/internal/dto"
	"user/internal/interface/providers"
	"user/internal/interface/repo"
)

type AdminAccessUsecase struct {
	AdminAccessRepo repo.AdminAccessRepository
	AuthMethodRepo  repo.AuthMethodRepository
	UserProvider    providers.ProfileProvider
}

func NewAdminAccessUsecase(
	adminAccessRepo repo.AdminAccessRepository,
	authMethodRepo repo.AuthMethodRepository,
	userProvider providers.ProfileProvider,
) *AdminAccessUsecase {
	return &AdminAccessUsecase{
		AdminAccessRepo: adminAccessRepo,
		AuthMethodRepo:  authMethodRepo,
		UserProvider:    userProvider,
	}
}

// CreateAdminAccess tạo quyền truy cập mới
func (u *AdminAccessUsecase) CreateAdminAccess(ctx context.Context, req dto.AdminAccessRequest) (*dto.AdminAccessResponse, error) {
	// Kiểm tra quyền admin
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if profileID == 0 {
	// 	return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	// }

	// Kiểm tra user tồn tại
	// user, err := u.AuthMethodRepo.FindByID(ctx, req.UserID)
	// if err != nil || user == nil {
	// 	return nil, _errors.ReturnError(int32(404), "Không tìm thấy người dùng")
	// }

	// Kiểm tra giới hạn thiết bị
	if req.DeviceID != "" {
		count, err := u.AdminAccessRepo.CountByUserID(ctx, req.UserID)
		if err != nil {
			return nil, _errors.ReturnError(int32(500), "Lỗi kiểm tra giới hạn thiết bị")
		}
		if count >= int64(req.MaxDevices) {
			return nil, _errors.ReturnError(int32(400), fmt.Sprintf("Tài khoản đã gắn với số thiết bị tối đa (%d)", req.MaxDevices))
		}
	}

	// Kiểm tra giới hạn IP
	if req.IPAddress != "" {
		count, err := u.AdminAccessRepo.CountByIPAddress(ctx, req.IPAddress)
		if err != nil {
			return nil, _errors.ReturnError(int32(500), "Lỗi kiểm tra giới hạn IP")
		}
		if count >= 3 {
			return nil, _errors.ReturnError(int32(400), "IP này đã được gắn với tối đa 3 tài khoản")
		}
	}

	// Tạo quyền truy cập
	access := &access.AdminAccessDomain{
		UserID:        req.UserID,
		IPAddress:     req.IPAddress,
		IPRange:       req.IPRange,
		DeviceID:      req.DeviceID,
		DeviceName:    req.DeviceName,
		DeviceType:    req.DeviceType,
		AuthType:      req.AuthType,
		Status:        "active",
		EffectiveFrom: req.EffectiveFrom,
		EffectiveTo:   req.EffectiveTo,
		MaxDevices:    req.MaxDevices,
		Require2FA:    req.Require2FA,
		// CreatedBy:     profileID,
		// UpdatedBy:     profileID,
	}

	if access.Status == "" {
		access.Status = "active"
	}
	if access.MaxDevices == 0 {
		access.MaxDevices = 3
	}

	createdAccess, err := u.AdminAccessRepo.Create(ctx, access)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi tạo quyền truy cập")
	}

	return u.convertToResponse(createdAccess), nil
}

// UpdateAdminAccess cập nhật quyền truy cập
func (u *AdminAccessUsecase) UpdateAdminAccess(ctx context.Context, id uint64, req dto.AdminAccessRequest) (*dto.AdminAccessResponse, error) {
	// Kiểm tra quyền admin
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if profileID == 0 {
	// 	return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	// }

	// Tìm quyền truy cập hiện tại
	existingAccess, err := u.AdminAccessRepo.FindByID(ctx, id)
	if err != nil || existingAccess == nil {
		return nil, _errors.ReturnError(int32(404), "Không tìm thấy quyền truy cập")
	}

	// Cập nhật thông tin
	existingAccess.IPAddress = req.IPAddress
	existingAccess.IPRange = req.IPRange
	existingAccess.DeviceID = req.DeviceID
	existingAccess.DeviceName = req.DeviceName
	existingAccess.DeviceType = req.DeviceType
	existingAccess.AuthType = req.AuthType
	existingAccess.Status = req.Status
	existingAccess.EffectiveFrom = req.EffectiveFrom
	existingAccess.EffectiveTo = req.EffectiveTo
	existingAccess.MaxDevices = req.MaxDevices
	existingAccess.Require2FA = req.Require2FA
	// existingAccess.UpdatedBy = profileID

	updatedAccess, err := u.AdminAccessRepo.Update(ctx, existingAccess)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi cập nhật quyền truy cập")
	}

	return u.convertToResponse(updatedAccess), nil
}

// DeleteAdminAccess xóa quyền truy cập
func (u *AdminAccessUsecase) DeleteAdminAccess(ctx context.Context, id uint64) error {
	// Kiểm tra quyền admin
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}

	// Tìm quyền truy cập
	access, err := u.AdminAccessRepo.FindByID(ctx, id)
	if err != nil || access == nil {
		return _errors.ReturnError(int32(404), "Không tìm thấy quyền truy cập")
	}

	// Soft delete
	err = u.AdminAccessRepo.SoftDelete(ctx, id)
	if err != nil {
		return _errors.ReturnError(int32(500), "Lỗi xóa quyền truy cập")
	}

	return nil
}

// GetAdminAccessList lấy danh sách quyền truy cập
func (u *AdminAccessUsecase) GetAdminAccessList(ctx context.Context, req dto.AdminAccessListRequest) (*dto.AdminAccessListResponse, error) {
	// Kiểm tra quyền admin
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if profileID == 0 {
	// 	return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	// }

	// Set default pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	// Build filters
	filters := make(map[string]interface{})
	if req.UserID != nil {
		filters["user_id"] = *req.UserID
	}
	if req.Status != "" {
		filters["status"] = req.Status
	}
	if req.AuthType != "" {
		filters["auth_type"] = req.AuthType
	}
	if req.IPAddress != "" {
		filters["ip_address"] = req.IPAddress
	}

	// Lấy danh sách
	accesses, total, err := u.AdminAccessRepo.FindAll(ctx, req.Page, req.Size, filters)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi lấy danh sách quyền truy cập")
	}

	// Lấy danh sách userID từ accesses
	userIDs := make([]uint64, 0, len(accesses))
	for _, access := range accesses {
		if access.UserID > 0 {
			userIDs = append(userIDs, access.UserID)
		}
	}

	// Lấy thông tin username từ user service
	userMap := make(map[uint64]*dto.ProfileDTO)
	if len(userIDs) > 0 {
		userMap, _ = u.UserProvider.GetMapProfileByIds(ctx, userIDs)
	}

	// Convert to response
	responses := make([]*dto.AdminAccessResponse, len(accesses))
	for i, access := range accesses {
		responses[i] = u.convertToResponse(access)
		// Thêm username từ userMap
		if user, ok := userMap[access.UserID]; ok && user != nil {
			responses[i].Username = user.FullName
		}
	}

	return &dto.AdminAccessListResponse{
		Data:  responses,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

// ValidateAdminAccess kiểm tra quyền truy cập
func (u *AdminAccessUsecase) ValidateAdminAccess(ctx context.Context, req dto.AdminAccessValidationRequest) (*dto.AdminAccessValidationResponse, error) {
	// Kiểm tra quyền truy cập
	result, err := u.AdminAccessRepo.ValidateAccess(ctx, req.UserID, req.IPAddress, req.DeviceID)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi kiểm tra quyền truy cập")
	}

	// Tạo log truy cập
	log := &access.AdminAccessLogDomain{
		UserID:    req.UserID,
		IPAddress: req.IPAddress,
		DeviceID:  req.DeviceID,
		Action:    "login",
		Status:    "success",
		Reason:    result.Reason,
	}

	if !result.IsAllowed {
		log.Status = "failed"
		log.Action = "access_denied"
	}

	u.AdminAccessRepo.CreateLog(ctx, log)

	return &dto.AdminAccessValidationResponse{
		IsAllowed:  result.IsAllowed,
		Reason:     result.Reason,
		Require2FA: result.Require2FA,
		AccessID:   result.AccessID,
		DeviceID:   result.DeviceID,
		IPAddress:  result.IPAddress,
	}, nil
}

// CreateAdminAccessLog tạo log truy cập
func (u *AdminAccessUsecase) CreateAdminAccessLog(ctx context.Context, req dto.AdminAccessLogRequest) error {
	log := &access.AdminAccessLogDomain{
		UserID:     req.UserID,
		IPAddress:  req.IPAddress,
		DeviceID:   req.DeviceID,
		DeviceName: req.DeviceName,
		Action:     req.Action,
		Status:     req.Status,
		Reason:     req.Reason,
		UserAgent:  req.UserAgent,
	}

	return u.AdminAccessRepo.CreateLog(ctx, log)
}

// GetAdminAccessLogList lấy danh sách log truy cập
func (u *AdminAccessUsecase) GetAdminAccessLogList(ctx context.Context, req dto.AdminAccessLogListRequest) (*dto.AdminAccessLogListResponse, error) {
	// Kiểm tra quyền admin
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}

	// Set default pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	var logs []*access.AdminAccessLogDomain
	var total int64
	var err error

	// Lấy log theo UserID hoặc IPAddress
	if req.UserID > 0 {
		logs, total, err = u.AdminAccessRepo.FindLogsByUserID(ctx, req.UserID, req.Page, req.Size)
	} else if req.IPAddress != "" {
		logs, total, err = u.AdminAccessRepo.FindLogsByIPAddress(ctx, req.IPAddress, req.Page, req.Size)
	} else {
		return nil, _errors.ReturnError(int32(400), "Cần cung cấp UserID hoặc IPAddress")
	}

	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi lấy danh sách log")
	}

	// Convert to response
	responses := make([]*dto.AdminAccessLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = &dto.AdminAccessLogResponse{
			ID:         log.ID,
			UserID:     log.UserID,
			IPAddress:  log.IPAddress,
			DeviceID:   log.DeviceID,
			DeviceName: log.DeviceName,
			Action:     log.Action,
			Status:     log.Status,
			Reason:     log.Reason,
			UserAgent:  log.UserAgent,
			CreatedAt:  log.CreatedAt,
		}
	}

	return &dto.AdminAccessLogListResponse{
		Data:  responses,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

// convertToResponse convert domain to response
func (u *AdminAccessUsecase) convertToResponse(access *access.AdminAccessDomain) *dto.AdminAccessResponse {
	return &dto.AdminAccessResponse{
		ID:            access.ID,
		UserID:        access.UserID,
		IPAddress:     access.IPAddress,
		IPRange:       access.IPRange,
		DeviceID:      access.DeviceID,
		DeviceName:    access.DeviceName,
		DeviceType:    access.DeviceType,
		AuthType:      access.AuthType,
		Status:        access.Status,
		EffectiveFrom: access.EffectiveFrom,
		EffectiveTo:   access.EffectiveTo,
		MaxDevices:    access.MaxDevices,
		Require2FA:    access.Require2FA,
		CreatedBy:     access.CreatedBy,
		UpdatedBy:     access.UpdatedBy,
		CreatedAt:     access.CreatedAt,
		UpdatedAt:     access.UpdatedAt,
	}
}
