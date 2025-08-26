package handler

import (
	"auth/internal/interface/repo"
	_errors "common/errors"
	_redis "common/redis"
	_utils "common/utils"
	"context"
	"fmt"
	authpb "pb/types/auth"
	"strconv"
	"strings"
)

type InternalHandler struct {
	authpb.UnimplementedAuthInternalServiceServer
	roleRepo     repo.RoleRepository
	redisService *_redis.RedisService
}

func NewInternalHandler(roleRepo repo.RoleRepository, redisService *_redis.RedisService) *InternalHandler {
	return &InternalHandler{
		roleRepo:     roleRepo,
		redisService: redisService,
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
	profileId := _utils.GetProfileIdWithContext(ctx)

	// -------------------
	// 1. Lấy roleIds từ Redis
	roleKey := "USER_ROLES_" + strconv.FormatUint(profileId, 10)
	roleIdsStr, err := h.redisService.Get(roleKey)
	var roleIds []string
	if err != nil || roleIdsStr == "" {
		// fallback: lấy từ DB
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
		roleIds = strings.Split(roleIdsStr, ",")
	}

	// -------------------
	// 2. Lấy permission_ids từ Redis
	permIdSet := make(map[string]struct{})
	missingRoleIds := []uint64{}
	for _, roleId := range roleIds {
		rolePermKey := "ROLE_PERMS_" + roleId
		permIdsStr, err := h.redisService.Get(rolePermKey)
		if err != nil || permIdsStr == "" {
			id, _ := strconv.ParseUint(roleId, 10, 64)
			missingRoleIds = append(missingRoleIds, id)
			continue
		}
		for _, pid := range strings.Split(permIdsStr, ",") {
			permIdSet[pid] = struct{}{}
		}
	}

	// fallback DB cho role chưa có trong Redis
	if len(missingRoleIds) > 0 {
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
	}

	// -------------------
	// 3. Lấy permission_keys từ Redis
	userPerms := make(map[string]struct{})
	missingPermIds := []uint64{}
	for pid := range permIdSet {
		permKey := "PERM_KEY_" + pid
		key, err := h.redisService.Get(permKey)
		if err != nil || key == "" {
			id, _ := strconv.ParseUint(pid, 10, 64)
			missingPermIds = append(missingPermIds, id)
			continue
		}
		userPerms[key] = struct{}{}
	}

	// fallback DB cho permission chưa có Redis
	if len(missingPermIds) > 0 {
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
		fmt.Println("missingPermIds", missingPermIds)
	}

	// -------------------
	// 4. Kiểm tra các permission cần thiết
	for _, required := range keys {
		if _, ok := userPerms[required]; ok {
			return nil
		}
	}

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
