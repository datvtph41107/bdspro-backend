package handler

import (
	"auth/internal/interface/repo"
	"auth/internal/usecase"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"log"
	authpb "pb/types/auth"
	sharepb "pb/types/shared"
	"strconv"
	"strings"
	"time"

	"auth/infra/mapper"
	"auth/infra/redis"
)

type InternalHandler struct {
	authpb.UnimplementedAuthInternalServiceServer
	roleRepo            repo.RoleRepository
	redisService        *redis.RedisProvider
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
	roleRepo repo.RoleRepository,
	redisService *redis.RedisProvider,
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

func (h *InternalHandler) RequiredRoles(ctx context.Context, req *authpb.RequiredRolesRequest) (*authpb.RequiredRolesResponse, error) {
	return &authpb.RequiredRolesResponse{}, nil
}
func (h *InternalHandler) RequiredPermissions(ctx context.Context, req *authpb.RequiredPermissionsRequest) (*authpb.RequiredPermissionsResponse, error) {
	err := h.HasPermissions(ctx, req.Permissions)
	if err != nil {
		return nil, err
	}
	return &authpb.RequiredPermissionsResponse{
		Status: true,
	}, nil
}

func (h *InternalHandler) HasPermissions(ctx context.Context, keys []string) error {
	funcStart := time.Now()
	profileId := _utils.GetProfileIdWithContext(ctx)

	// -------------------
	// 1. Lấy roleIds từ Redis
	stage1Start := time.Now()
	roleKey := "USER_ROLES_" + strconv.FormatUint(profileId, 10)
	roleIdsStr, err := h.redisService.Get(roleKey)
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
		h.redisService.Set(roleKey, strings.Join(roleIds, ","))
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
			h.redisService.Set(rolePermKey, strings.Join(strIds, ","))
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
			h.redisService.Set(rk, key)
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
	if err := h.redisService.DeleteAllKeysWithPrefix(fmt.Sprintf("ROLE_PERMS_%d", roleId)); err != nil {
		return err
	}
	// if err := h.redisService.DeleteAllKeysWithPrefix("PERM_KEY_*"); err != nil {
	// 	return err
	// }
	return nil
}

// ClearUserRolesCache xóa cache role của 1 profile sau khi gán role mới
func (h *InternalHandler) ClearUserRolesCache(profileId uint64) {
	if profileId == 0 || h.redisService == nil {
		return
	}
	_ = h.redisService.DeleteAllKeysWithPrefix("USER_ROLES_" + strconv.FormatUint(profileId, 10))
}

func (h *InternalHandler) GetRoleById(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Role, error) {
	role, err := h.roleRepo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return h.roleMapper.MapRoleToPb(role), nil
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

func (h *InternalHandler) GetAdminRoleIdByGroupKey(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Role, error) {
	role, err := h.roleRepo.GetAdminRoleIdByGroupKey(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var pbRole *sharepb.Role
	if role != nil {
		pbRole = h.roleMapper.MapRoleToPb(role)
	}

	return pbRole, nil
}

func (h *InternalHandler) GetProfileByIds(ctx context.Context, req *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error) {
	// Lấy thông tin từ AuthMethod domain thay vì gọi user service
	authMethods, err := h.roleRepo.GetAuthMethodsByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}

	profiles := make([]*sharepb.ProfileItem, len(authMethods))
	for i, auth := range authMethods {
		profiles[i] = &sharepb.ProfileItem{
			Id:       auth.UserID,
			FullName: auth.FullName,
			Avatar:   auth.Avatar,
			Phone:    auth.Phone,
		}
	}

	return &sharepb.GetProfileByIdsResponse{
		Profiles: profiles,
	}, nil
}

func (h *InternalHandler) GetRolesByRoleKeys(ctx context.Context, req *authpb.GetRolesByRoleKeysRequest) (*authpb.RoleListResponse, error) {
	roles, err := h.roleRepo.FindByRoleKeys(ctx, req.RoleKeys)
	if err != nil {
		return nil, err
	}

	result := h.roleMapper.MapRoleListToPb(roles)
	return &authpb.RoleListResponse{
		Data: result,
	}, nil
}

func (h *InternalHandler) CreateAdmin(ctx context.Context, req *authpb.AuthMethod) (*authpb.AuthMethod, error) {
	entity := h.authMethodMapper.MapAuthMethodPbToDomain(req)

	result, err := h.authMethodUsecase.CreateAdmin(ctx, entity)
	if err != nil {
		return nil, err
	}

	// Convert DTO response to proto response
	return h.authMethodMapper.MapAuthMethodPb(result), nil
}

func (h *InternalHandler) UpdateAdmin(ctx context.Context, req *authpb.AuthMethod) (*authpb.AuthMethod, error) {
	// Convert proto request to DTO
	en := h.authMethodMapper.MapAuthMethodPbToDomain(req)

	result, err := h.authMethodUsecase.UpdateAdmin(ctx, en)
	if err != nil {
		return nil, err
	}

	// Convert DTO response to proto response
	return h.authMethodMapper.MapAuthMethodPb(result), nil
}

func (h *InternalHandler) DeleteAdmin(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	err := h.authMethodUsecase.DeleteAdmin(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Message: "Admin đã được xóa thành công",
	}, nil
}

func (h *InternalHandler) GetAdminById(ctx context.Context, req *sharepb.IdRequest) (*authpb.AuthMethod, error) {
	result, err := h.authMethodUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return h.authMethodMapper.MapAuthMethodPb(result), nil
}

// GetAuthUsersByProfileIDs lấy thông tin AuthUser theo danh sách profile IDs (Internal API)
func (h *InternalHandler) GetAuthUsersByProfileIDs(ctx context.Context, req *authpb.GetAuthUsersByProfileIDsRequest) (*authpb.GetAuthUsersByProfileIDsResponse, error) {
	// Call usecase
	result, err := h.userInfoUsecase.GetAuthUsersByProfileIDs(ctx, req.ProfileIds)
	if err != nil {
		return nil, err
	}

	// Convert DTO response to proto response
	var users []*authpb.AuthUserStatusInfo
	for _, user := range result {
		users = append(users, h.userInfoMapper.MapAuthUserStatusInfoToPb(&user))
	}

	return &authpb.GetAuthUsersByProfileIDsResponse{
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
func (h *InternalHandler) GetAuthAdminByIds(ctx context.Context, req *authpb.GetAuthAdminByIdsRequest) (*authpb.GetAuthAdminByIdsResponse, error) {
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

	var pbAdmins []*authpb.AuthAdminWithRole
	for _, admin := range admins {
		if admin == nil {
			continue
		}
		pbAdmins = append(pbAdmins, h.roleMapper.MapAuthAdminWithRoleToPb(admin))
	}

	return &authpb.GetAuthAdminByIdsResponse{
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
func (h *InternalHandler) GetPushTokensByProfileId(ctx context.Context, req *sharepb.IdRequest) (*authpb.GetPushTokensResponse, error) {
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


	return &authpb.GetPushTokensResponse{
		PushTokens: pushTokens,
	}, nil
}
