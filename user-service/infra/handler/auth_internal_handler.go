package handler

import (
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"log"
	authpb "pb/types/auth"
	"strconv"
	"strings"
	"time"
	"user/internal/interface/providers"
	"user/internal/interface/repo"
	"user/internal/usecase"

	"user/infra/mapper"
)

type AuthInternalHandler struct {
	authpb.UnimplementedAuthInternalServiceServer
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

func NewAuthInternalHandler(
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
) *AuthInternalHandler {
	return &AuthInternalHandler{
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

func (h *AuthInternalHandler) RequiredRoles(ctx context.Context, req *authpb.RequiredRolesRequest) (*authpb.RequiredRolesResponse, error) {
	return &authpb.RequiredRolesResponse{}, nil
}
func (h *AuthInternalHandler) RequiredPermissions(ctx context.Context, req *authpb.RequiredPermissionsRequest) (*authpb.RequiredPermissionsResponse, error) {
	err := h.HasPermissions(ctx, req.Permissions)
	if err != nil {
		return nil, err
	}
	return &authpb.RequiredPermissionsResponse{
		Status: true,
	}, nil
}

func (h *AuthInternalHandler) HasPermissions(ctx context.Context, keys []string) error {
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
		h.redisService.Set(ctx, roleKey, strings.Join(roleIds, ","), 0)
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
			h.redisService.Set(ctx, rolePermKey, strings.Join(strIds, ","), 0)
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
			h.redisService.Set(ctx, rk, key, 0)
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

func (h *AuthInternalHandler) ResetPermission(ctx context.Context, roleId uint64) error {
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
func (h *AuthInternalHandler) ClearUserRolesCache(ctx context.Context, profileId uint64) {
	if profileId == 0 || h.redisService == nil {
		return
	}
	_ = h.redisService.DeleteAllKeysWithPrefix(ctx, "USER_ROLES_"+strconv.FormatUint(profileId, 10))
}
