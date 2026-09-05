package handler

import (
	_enum "common/domain/enum"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"log"
	authpb "pb/types/auth"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"strconv"
	"strings"
	"time"
	"user/infra/mapper"
	"user/internal/domain/access"
	"user/internal/dto"
	"user/internal/interface/providers"
	"user/internal/interface/repo"
	"user/internal/models"
	"user/internal/usecase"
	"user/internal/usecases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InternalHandler struct {
	userpb.UnimplementedInternalUserServiceServer
	profileUsecase *usecases.ProfileUsecase
	adminUsecase   *usecases.AdminUsecase

	roleRepo            repo.RoleRepository
	redisService        providers.CacheProvider
	roleMapper          *mapper.RoleMapper
	authMethodMapper    *mapper.AuthMethodMapper
	authMethodUsecase   *usecase.AuthAdminUsecase
	userInfoUsecase     *usecase.AuthUserProfileUsecase
	userInfoMapper      *mapper.UserInfoMapper
	authSecurityUsecase *usecase.AuthSecurityUsecase
	authSecurityMapper  *mapper.AuthSecurityMapper
	deviceRepo          repo.DeviceRepository
}

func NewInternalHandler(
	profileUsecase *usecases.ProfileUsecase,
	adminUsecase *usecases.AdminUsecase,

	roleRepo repo.RoleRepository,
	redisService providers.CacheProvider,
	roleMapper *mapper.RoleMapper,
	authMethodMapper *mapper.AuthMethodMapper,
	authMethodUsecase *usecase.AuthAdminUsecase,
	userInfoUsecase *usecase.AuthUserProfileUsecase,
	userInfoMapper *mapper.UserInfoMapper,
	authSecurityUsecase *usecase.AuthSecurityUsecase,
	authSecurityMapper *mapper.AuthSecurityMapper,
	deviceRepo repo.DeviceRepository,
) *InternalHandler {
	return &InternalHandler{
		profileUsecase: profileUsecase,
		adminUsecase:   adminUsecase,

		roleRepo:            roleRepo,
		redisService:        redisService,
		roleMapper:          roleMapper,
		authMethodMapper:    authMethodMapper,
		authMethodUsecase:   authMethodUsecase,
		userInfoUsecase:     userInfoUsecase,
		userInfoMapper:      userInfoMapper,
		authSecurityUsecase: authSecurityUsecase,
		authSecurityMapper:  authSecurityMapper,
		deviceRepo:          deviceRepo,
	}
}

// CreateProfile tạo profile mới (dùng bởi auth-service)
func (h *InternalHandler) CreateProfile(ctx context.Context, req *userpb.CreateProfileRequest) (*userpb.CreateProfileResponse, error) {
	// Tạo UserProfileEntity từ request
	profile := &models.UserProfileEntity{
		// ProfileID:    req.ProfileId,
		FullName:     req.FullName,
		Email:        req.Email,
		Phone:        req.Phone,
		Avatar:       req.Avatar,
		ReferralCode: req.ReferralCode,
	}
	// Set tất cả visibility thành public mặc định
	models.SetAllVisibilityPublic(profile)

	// Gọi repo để tạo profile
	createdProfile, err := h.profileUsecase.ProfileRepo.CreateUserProfile(ctx, profile)
	if err != nil {
		log.Printf("Failed to create profile: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to create profile: %v", err)
	}

	// Trả về response
	return &userpb.CreateProfileResponse{
		ProfileId:    createdProfile.ProfileID,
		FullName:     createdProfile.FullName,
		Email:        createdProfile.Email,
		Phone:        createdProfile.Phone,
		Avatar:       createdProfile.Avatar,
		ReferralCode: createdProfile.ReferralCode,
		PlanId:       req.PlanId,
	}, nil
}

// UpdateProfile cập nhật profile (dùng bởi auth-service)
func (h *InternalHandler) UpdateProfile(ctx context.Context, req *userpb.UpdateProfileRequest) (*sharepb.Empty, error) {
	// Lấy profile hiện tại
	existingProfile, err := h.profileUsecase.ProfileRepo.GetUserDetailByID(ctx, req.ProfileId)
	if err != nil {
		log.Printf("Failed to get profile: %v", err)
		return nil, status.Errorf(codes.NotFound, "Profile not found")
	}

	// Cập nhật các field nếu có
	if req.FullName != nil {
		existingProfile.FullName = *req.FullName
	}
	if req.Email != nil {
		existingProfile.Email = *req.Email
	}
	if req.Phone != nil {
		existingProfile.Phone = *req.Phone
	}
	if req.Avatar != nil {
		existingProfile.Avatar = *req.Avatar
	}

	// Gọi repo để cập nhật
	_, err = h.profileUsecase.ProfileRepo.UpdateUserProfile(ctx, existingProfile)
	if err != nil {
		log.Printf("Failed to update profile: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to update profile: %v", err)
	}

	return &sharepb.Empty{}, nil
}

// SoftDeleteProfile xóa mềm profile (dùng bởi auth-service)
func (h *InternalHandler) SoftDeleteProfile(ctx context.Context, req *userpb.SoftDeleteProfileRequest) (*sharepb.Empty, error) {
	// Lấy thông tin người xóa từ context (nếu có)
	var deletedBy *uint64
	if profileID := _utils.GetProfileIdWithContext(ctx); profileID > 0 {
		deletedBy = &profileID
	}

	reason := "Xóa từ auth service"

	// Xóa profile
	err := h.profileUsecase.ProfileRepo.DeleteUserProfile(ctx, req.ProfileId, deletedBy, reason)
	if err != nil {
		log.Printf("Failed to soft delete profile: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to soft delete profile: %v", err)
	}

	// Xóa TẤT CẢ auth method của user
	err = h.profileUsecase.AuthProvider.DeleteAllAuthMethodsByUserId(ctx, req.ProfileId)
	if err != nil {
		log.Printf("Warning: Failed to delete all auth methods for user %d: %v", req.ProfileId, err)
		// Không return error vì profile đã được xóa
	}

	return &sharepb.Empty{}, nil
}

// HardDeleteProfile xóa vĩnh viễn profile khỏi database (dùng bởi auth-service)
func (h *InternalHandler) HardDeleteProfile(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Empty, error) {
	// Lấy thông tin người xóa từ context (nếu có)
	var deletedBy *uint64
	if profileID := _utils.GetProfileIdWithContext(ctx); profileID > 0 {
		deletedBy = &profileID
	}

	reason := "Hard delete từ auth service"

	// Xóa profile
	err := h.profileUsecase.ProfileRepo.DeleteUserProfile(ctx, req.Id, deletedBy, reason)
	if err != nil {
		log.Printf("Failed to hard delete profile: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to hard delete profile: %v", err)
	}

	// Xóa TẤT CẢ auth method của user
	err = h.profileUsecase.AuthProvider.DeleteAllAuthMethodsByUserId(ctx, req.Id)
	if err != nil {
		log.Printf("Warning: Failed to delete all auth methods for user %d: %v", req.Id, err)
		// Không return error vì profile đã được xóa
	}

	return &sharepb.Empty{}, nil
}

// GetProfileByID lấy profile theo ID
func (h *InternalHandler) GetProfileByID(ctx context.Context, req *userpb.GetProfileByIDRequest) (*userpb.ProfileResponse, error) {
	profile, err := h.profileUsecase.ProfileRepo.GetUserDetailByID(ctx, req.ProfileId)
	if err != nil {
		log.Printf("Failed to get profile by ID: %v", err)
		return nil, status.Errorf(codes.NotFound, "Profile not found")
	}

	return &userpb.ProfileResponse{
		ProfileId:    profile.ProfileID,
		FullName:     profile.FullName,
		Email:        profile.Email,
		Phone:        profile.Phone,
		Avatar:       profile.Avatar,
		ReferralCode: profile.ReferralCode,
		CreatedAt:    _utils.FormatTimeToString(profile.CreatedAt),
		UpdatedAt:    _utils.FormatTimeToString(profile.UpdatedAt),
	}, nil
}

func (h *InternalHandler) GetProfileByIDV3(ctx context.Context, req *sharepb.RequestV3Proto) (*sharepb.UserV3Proto, error) {
	profile, err := h.profileUsecase.ProfileRepo.GetUserDetailByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Profile not found")
	}

	return &sharepb.UserV3Proto{
		ProfileId:         profile.ProfileID,
		FullName:          profile.FullName,
		Avatar:            profile.Avatar,
		RoleRealEstate:    uint32(profile.RoleRealEstate),
		Visibility:        uint32(profile.Visibility),
		ProfileVisibility: uint32(profile.ProfileVisibility),
		StatusOnline:      uint32(profile.StatusOnline),
	}, nil
}

func (h *InternalHandler) LockUser(ctx context.Context, req *sharepb.RequestV3Proto) (*sharepb.Empty, error) {
	err := h.adminUsecase.LockUser(ctx, &dto.LockUserRequest{
		ProfileID: req.Id,
		Reason:    req.Text,
		LockType:  _enum.EUserStatusTemporaryLocked,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to lock account: %v", err)
	}

	return &sharepb.Empty{}, nil
}

// UpdateLastSeen cập nhật thời gian hoạt động gần nhất của user (dùng bởi relay-service khi disconnect)
func (h *InternalHandler) UpdateLastSeen(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Empty, error) {
	err := h.profileUsecase.ProfileRepo.UpdateLastSeen(ctx, req.Id)
	if err != nil {
		log.Printf("Failed to update last seen for profile %d: %v", req.Id, err)
		return nil, status.Errorf(codes.Internal, "Failed to update last seen: %v", err)
	}

	return &sharepb.Empty{}, nil
}

func (h *InternalHandler) HasPermissions(ctx context.Context, keys []string) error {
	funcStart := time.Now()
	profileId := _utils.GetProfileIdWithContext(ctx)

	// -------------------
	// 1. Lấy roleIds từ Redis
	stage1Start := time.Now()
	roleKey := "USER_ROLES_" + strconv.FormatUint(profileId, 10)
	roleIdsStr, err := h.redisService.Get(ctx, roleKey)
	var roleIds []string
	if err != nil || roleIdsStr == "" {
		// fallback: lấy từ DB
		log.Printf("[HasPermissions] Cache miss for roleIds for profileId %d, fetching from DB", profileId)
		roleIdsUint, err := h.roleRepo.GetRoleIdsByProfileId(ctx, profileId)
		if err != nil {
			return err
		}
		for _, id := range roleIdsUint {
			roleIds = append(roleIds, strconv.FormatUint(id, 10))
		}
		// lưu lại Redis
		h.redisService.Set(ctx, roleKey, strings.Join(roleIds, ","), time.Duration(0))
	} else {
		log.Printf("[HasPermissions] Cache hit for roleIds for profileId %d", profileId)
		roleIds = strings.Split(roleIdsStr, ",")
	}
	log.Printf("[HasPermissions] Stage 1 (Get roleIds) took %s", time.Since(stage1Start))

	// -------------------
	// 2. Lấy permission_ids từ Redis (sử dụng MGet nếu có thể)
	stage2Start := time.Now()
	permIdSet := make(map[string]struct{})
	missingRoleIds := []uint64{}

	// Thu thập tất cả các rolePermKey cần tìm
	var rolePermKeys []string
	roleIdToPermIdMap := make(map[string]string) // Để dễ dàng tìm lại roleId khi có kết quả MGet

	for _, roleId := range roleIds {
		rolePermKey := "ROLE_PERMS_" + roleId
		rolePermKeys = append(rolePermKeys, rolePermKey)
		roleIdToPermIdMap[rolePermKey] = roleId // Lưu lại roleId tương ứng
	}

	if len(rolePermKeys) > 0 {
		permIdsVals, err := h.redisService.MGet(ctx, rolePermKeys...)
		if err != nil {
			log.Printf("[HasPermissions] Error MGet role perm keys from Redis: %v", err)
			// Xử lý lỗi MGet, có thể fallback về Get từng cái hoặc bỏ qua cache
		} else {
			for i, val := range permIdsVals {
				if val != nil {
					permIdsStr, ok := val.(string)
					if ok && permIdsStr != "" {
						for _, pid := range strings.Split(permIdsStr, ",") {
							permIdSet[pid] = struct{}{}
						}
					} else {
						// Giá trị null hoặc không phải string, coi như cache miss
						id, _ := strconv.ParseUint(roleIdToPermIdMap[rolePermKeys[i]], 10, 64)
						missingRoleIds = append(missingRoleIds, id)
					}
				} else {
					// Cache miss cho khóa này
					id, _ := strconv.ParseUint(roleIdToPermIdMap[rolePermKeys[i]], 10, 64)
					missingRoleIds = append(missingRoleIds, id)
				}
			}
		}
	}
	log.Printf("[HasPermissions] Initial permIds lookup from Redis (MGet) took %s", time.Since(stage2Start))

	// fallback DB cho role chưa có trong Redis
	if len(missingRoleIds) > 0 {
		log.Printf("[HasPermissions] Cache miss for permIds for roles %v, fetching from DB", missingRoleIds)
		dbPermIds, err := h.roleRepo.GetPermissionIdsByRoleIds(ctx, missingRoleIds)
		if err != nil {
			return err
		}
		for rid, permIds := range dbPermIds {
			rolePermKey := "ROLE_PERMS_" + strconv.FormatUint(rid, 10)
			strIds := []string{}
			for _, pid := range permIds {
				pidStr := strconv.FormatUint(pid, 10)
				strIds = append(strIds, pidStr)
				permIdSet[pidStr] = struct{}{}
			}
			// lưu Redis
			h.redisService.Set(ctx, rolePermKey, strings.Join(strIds, ","), time.Duration(0))
		}
		log.Printf("[HasPermissions] Fallback DB for permIds took %s", time.Since(stage2Start))
	}

	// -------------------
	// 3. Lấy permission_keys từ Redis (sử dụng MGet)
	stage3Start := time.Now()
	userPerms := make(map[string]struct{})
	missingPermIds := []uint64{}

	var permKeysToFetch []string
	permIdToPermKeyMap := make(map[string]string) // Để dễ dàng tìm lại pid khi có kết quả MGet

	for pid := range permIdSet {
		permKey := "PERM_KEY_" + pid
		permKeysToFetch = append(permKeysToFetch, permKey)
		permIdToPermKeyMap[permKey] = pid // Lưu lại pid tương ứng
	}

	if len(permKeysToFetch) > 0 {
		keysVals, err := h.redisService.MGet(ctx, permKeysToFetch...)
		if err != nil {
			log.Printf("[HasPermissions] Error MGet perm keys from Redis: %v", err)
			// Xử lý lỗi MGet, có thể fallback về Get từng cái hoặc bỏ qua cache
		} else {
			for i, val := range keysVals {
				if val != nil {
					keyStr, ok := val.(string)
					if ok && keyStr != "" {
						userPerms[keyStr] = struct{}{}
					} else {
						// Giá trị null hoặc không phải string, coi như cache miss
						id, _ := strconv.ParseUint(permIdToPermKeyMap[permKeysToFetch[i]], 10, 64)
						missingPermIds = append(missingPermIds, id)
					}
				} else {
					// Cache miss cho khóa này
					id, _ := strconv.ParseUint(permIdToPermKeyMap[permKeysToFetch[i]], 10, 64)
					missingPermIds = append(missingPermIds, id)
				}
			}
		}
	}
	log.Printf("[HasPermissions] Initial permKeys lookup from Redis (MGet) took %s", time.Since(stage3Start))

	// fallback DB cho permission chưa có Redis
	if len(missingPermIds) > 0 {
		log.Printf("[HasPermissions] Cache miss for permKeys for permIds %v, fetching from DB", missingPermIds)
		keys, err := h.roleRepo.GetPermissionKeysByIds(ctx, missingPermIds)
		if err != nil {
			return err
		}
		for i, key := range keys {
			userPerms[key] = struct{}{}
			// lưu Redis
			rk := "PERM_KEY_" + strconv.FormatUint(missingPermIds[i], 10)
			h.redisService.Set(ctx, rk, key, time.Duration(0))
		}
		log.Printf("[HasPermissions] Fallback DB for permKeys took %s", time.Since(stage3Start))
	}

	// -------------------
	// 4. Kiểm tra các permission cần thiết
	for _, required := range keys {
		if _, ok := userPerms[required]; ok {
			log.Printf("[HasPermissions] Permission check succeeded for profileId %d, took %s total", profileId, time.Since(funcStart))
			return nil
		}
	}
	log.Printf("[HasPermissions] Permission check failed for profileId %d, took %s total", profileId, time.Since(funcStart))

	return _errors.ReturnError(403, "Bạn không có quyền truy cập")
}

func (h *InternalHandler) ResetPermission(ctx context.Context, roleId uint64) error {
	// if err := h.redisService.DeleteAllKeysWithPrefix("USER_ROLES_*"); err != nil {
	// 	return err
	// }
	if err := h.redisService.DeleteAllKeysWithPrefix(ctx, fmt.Sprintf("ROLE_PERMS_%d", roleId)); err != nil {
		return err
	}
	// if err := h.redisService.DeleteAllKeysWithPrefix("PERM_KEY_*"); err != nil {
	// 	return err
	// }
	return nil
}

// ClearUserRolesCache xóa cache role của 1 profile sau khi gán role mới
func (h *InternalHandler) ClearUserRolesCache(ctx context.Context, profileId uint64) {
	if profileId == 0 || h.redisService == nil {
		return
	}
	_ = h.redisService.DeleteAllKeysWithPrefix(ctx, "USER_ROLES_"+strconv.FormatUint(profileId, 10))
}

func (h *InternalHandler) GetRoleById(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Role, error) {
	role, err := h.roleRepo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return h.roleMapper.MapRoleToPb(role), nil
}

// GetAdminRoleIdByGroupKey là compatibility RPC cho AuthInternal production.
// Role vẫn được đọc từ repository User, không tạo thêm owner thứ hai.
func (h *InternalHandler) GetAdminRoleIdByGroupKey(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Role, error) {
	if req == nil || req.Id == 0 {
		return nil, _errors.ReturnError(400, "groupKey là bắt buộc")
	}
	role, err := h.roleRepo.GetAdminRoleIdByGroupKey(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, status.Error(codes.NotFound, "admin role không tồn tại")
	}
	return h.roleMapper.MapRoleToPb(role), nil
}

// GetProfileByIds dùng profile repository canonical của User thay vì đọc
// chéo auth_method như source production cũ.
func (h *InternalHandler) GetProfileByIds(ctx context.Context, req *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error) {
	if req == nil || len(req.Ids) == 0 {
		return &sharepb.GetProfileByIdsResponse{Profiles: []*sharepb.ProfileItem{}}, nil
	}
	profiles, err := h.profileUsecase.ProfileRepo.GetProfileByIds(req.Ids)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get profiles: %v", err)
	}
	items := make([]*sharepb.ProfileItem, 0, len(profiles))
	for _, profile := range profiles {
		if profile == nil {
			continue
		}
		items = append(items, &sharepb.ProfileItem{
			Id:           profile.ProfileID,
			FullName:     profile.FullName,
			Avatar:       profile.Avatar,
			TickVerified: profile.TickVerified,
			Phone:        profile.Phone,
		})
	}
	return &sharepb.GetProfileByIdsResponse{Profiles: items}, nil
}

// CreateAdmin/UpdateAdmin/GetAdminById/DeleteAdmin là facade nội bộ để Auth
// forward contract cũ; credential và profile đều được User usecase xử lý.
func (h *InternalHandler) CreateAdmin(ctx context.Context, req *authpb.AuthMethod) (*authpb.AuthMethod, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không hợp lệ")
	}
	entity := h.authMethodMapper.MapAuthMethodPbToDomain(req)
	result, err := h.authMethodUsecase.CreateAdmin(ctx, entity)
	if err != nil {
		return nil, err
	}
	return h.authMethodMapper.MapAuthMethodPb(result), nil
}

func (h *InternalHandler) UpdateAdmin(ctx context.Context, req *authpb.AuthMethod) (*authpb.AuthMethod, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không hợp lệ")
	}
	entity := h.authMethodMapper.MapAuthMethodPbToDomain(req)
	result, err := h.authMethodUsecase.UpdateAdmin(ctx, entity)
	if err != nil {
		return nil, err
	}
	return h.authMethodMapper.MapAuthMethodPb(result), nil
}

func (h *InternalHandler) DeleteAdmin(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, _errors.ReturnError(400, "admin id là bắt buộc")
	}
	if err := h.authMethodUsecase.DeleteAdmin(ctx, req.Id); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Message: "Admin đã được xóa thành công"}, nil
}

func (h *InternalHandler) GetAdminById(ctx context.Context, req *sharepb.IdRequest) (*authpb.AuthMethod, error) {
	if req == nil || req.Id == 0 {
		return nil, _errors.ReturnError(400, "admin id là bắt buộc")
	}
	result, err := h.authMethodUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return h.authMethodMapper.MapAuthMethodPb(result), nil
}

func (h *InternalHandler) GetRolesByIds(ctx context.Context, req *sharepb.IdRequest) (*authpb.RoleListResponse, error) {
	roles, err := h.roleRepo.GetListByIDs(ctx, req.Ids)
	if err != nil {
		return nil, err
	}

	result := h.roleMapper.MapRoleListToPb(roles)
	return &authpb.RoleListResponse{
		Data: result,
	}, nil
}

func (h *InternalHandler) GetRolesByRoleKeys(ctx context.Context, req *userpb.GetRolesByRoleKeysRequest) (*authpb.RoleListResponse, error) {
	roles, err := h.roleRepo.FindByRoleKeys(ctx, req.RoleKeys)
	if err != nil {
		return nil, err
	}

	result := h.roleMapper.MapRoleListToPb(roles)
	return &authpb.RoleListResponse{
		Data: result,
	}, nil
}

func (h *InternalHandler) GetRolePermissions(ctx context.Context, req *sharepb.IdRequest) (*userpb.GetRolePermissionsResponse, error) {
	var roles []access.Role
	var err error

	switch {
	case req.Id > 0:
		role, getErr := h.roleRepo.GetByID(ctx, req.Id)
		if getErr != nil {
			return nil, getErr
		}
		if role != nil {
			roles = []access.Role{*role}
		}
	case len(req.Ids) > 0:
		roles, err = h.roleRepo.GetListByIDs(ctx, req.Ids)
		if err != nil {
			return nil, err
		}
	default:
		roles, err = h.roleRepo.GetAll(ctx)
		if err != nil {
			return nil, err
		}
	}

	roleIds := make([]uint64, 0, len(roles))
	for _, role := range roles {
		roleIds = append(roleIds, role.ID)
	}

	permMap, err := h.roleRepo.GetPermissionIdsByRoleIds(ctx, roleIds)
	if err != nil {
		return nil, err
	}

	data := make([]*userpb.RolePermissionItem, 0, len(roles))
	for _, role := range roles {
		permIds := permMap[role.ID]
		pbRole := h.roleMapper.MapRoleToPb(&role)
		if pbRole != nil {
			pbRole.PermissionIds = permIds
		}
		data = append(data, &userpb.RolePermissionItem{
			Role:          pbRole,
			PermissionIds: permIds,
		})
	}

	return &userpb.GetRolePermissionsResponse{Data: data}, nil
}

// GetAuthUsersByProfileIDs lấy thông tin AuthUser theo danh sách profile IDs (Internal API)
func (h *InternalHandler) GetAuthUsersByProfileIDs(ctx context.Context, req *userpb.GetAuthUsersByProfileIDsRequest) (*userpb.GetAuthUsersByProfileIDsResponse, error) {
	// Call usecase
	result, err := h.userInfoUsecase.GetAuthUsersByProfileIDs(ctx, req.ProfileIds)
	if err != nil {
		return nil, err
	}

	// Convert DTO response to proto response
	var users []*userpb.AuthUserStatusInfo
	for _, user := range result {
		users = append(users, h.userInfoMapper.MapAuthUserStatusInfoToPb(&user))
	}

	return &userpb.GetAuthUsersByProfileIDsResponse{
		Data: users,
	}, nil
}

// DeleteAuthMethodsByUserId xóa tất cả auth_method của 1 user
func (h *InternalHandler) DeleteAuthMethodsByUserId(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	err := h.authMethodUsecase.DeleteAllAuthMethodsByUserId(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Message: fmt.Sprintf("Đã xóa tất cả auth methods của user %d", req.Id),
	}, nil
}

// GetAuthAdminByIds lấy thông tin admin kèm với role theo danh sách IDs
func (h *InternalHandler) GetAuthAdminByIds(ctx context.Context, req *userpb.GetAuthAdminByIdsRequest) (*userpb.GetAuthAdminByIdsResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không hợp lệ")
	}
	if h.roleRepo == nil {
		return nil, _errors.ReturnError(500, "role repository chưa được cấu hình")
	}
	if h.roleMapper == nil {
		return nil, _errors.ReturnError(500, "role mapper chưa được cấu hình")
	}

	admins, err := h.roleRepo.GetAuthAdminsWithRolesByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}

	var pbAdmins []*userpb.AuthAdminWithRole
	for _, admin := range admins {
		if admin == nil {
			continue
		}
		pbAdmins = append(pbAdmins, h.roleMapper.MapAuthAdminWithRoleToPb(admin))
	}

	return &userpb.GetAuthAdminByIdsResponse{
		Data: pbAdmins,
	}, nil
}

// GetAuthDataByAuthId trả về thông tin status, pin và devices theo authId
func (h *InternalHandler) GetAuthDataByAuthId(ctx context.Context, req *sharepb.RequestV3Proto) (*sharepb.AuthDataV3Proto, error) {
	if req == nil || req.Id == 0 {
		return nil, _errors.ReturnError(400, "authId là bắt buộc")
	}

	result, err := h.authSecurityUsecase.GetAuthSecurityData(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return h.authSecurityMapper.MapToProto(result), nil
}

// GetPushTokensByProfileId lấy danh sách push token theo profile ID (Internal API cho notification service)
func (h *InternalHandler) GetPushTokensByProfileId(ctx context.Context, req *sharepb.IdRequest) (*userpb.GetPushTokensResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, _errors.ReturnError(400, "profileId là bắt buộc")
	}

	devices, err := h.deviceRepo.ListByProfileID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var pushTokens []string
	for _, device := range devices {
		if device.PushToken != "" {
			pushTokens = append(pushTokens, device.PushToken)
		}
	}

	return &userpb.GetPushTokensResponse{
		PushTokens: pushTokens,
	}, nil
}

func (h *InternalHandler) GetRoleIdsByProfileId(ctx context.Context, req *sharepb.IdRequest) (*userpb.GetRoleIdsByProfileIdResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, _errors.ReturnError(400, "profileId là bắt buộc")
	}

	roleIds, err := h.roleRepo.GetRoleIdsByProfileId(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &userpb.GetRoleIdsByProfileIdResponse{
		RoleIds: roleIds,
	}, nil
}
