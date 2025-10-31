package handler

import (
	"auth/internal/interface/repo"
	"auth/internal/usecase"
	_errors "common/errors"
	_redis "common/redis"
	_utils "common/utils"
	"context"
	"fmt"
	authpb "pb/types/auth"
	sharepb "pb/types/shared"
	"strconv"
	"strings"

	"auth/infra/mapper"
)

type InternalHandler struct {
	authpb.UnimplementedAuthInternalServiceServer
	roleRepo          repo.RoleRepository
	redisService      *_redis.RedisService
	roleMapper        *mapper.RoleMapper
	authMethodMapper  *mapper.AuthMethodMapper
	authMethodUsecase *usecase.AuthAdminUsecase
	userInfoUsecase   *usecase.AuthUserProfileUsecase
	userInfoMapper    *mapper.UserInfoMapper
}

func NewInternalHandler(
	roleRepo repo.RoleRepository,
	redisService *_redis.RedisService,
	roleMapper *mapper.RoleMapper,
	authMethodMapper *mapper.AuthMethodMapper,
	authMethodUsecase *usecase.AuthAdminUsecase,
	userInfoUsecase *usecase.AuthUserProfileUsecase,
	userInfoMapper *mapper.UserInfoMapper,
) *InternalHandler {
	return &InternalHandler{
		roleRepo:          roleRepo,
		redisService:      redisService,
		roleMapper:        roleMapper,
		authMethodMapper:  authMethodMapper,
		authMethodUsecase: authMethodUsecase,
		userInfoUsecase:   userInfoUsecase,
		userInfoMapper:    userInfoMapper,
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

func (h *InternalHandler) GetRoleById(ctx context.Context, req *sharepb.IdRequest) (*authpb.Role, error) {
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

func (h *InternalHandler) GetAdminRoleIdByGroupKey(ctx context.Context, req *sharepb.IdRequest) (*authpb.Role, error) {
	role, err := h.roleRepo.GetAdminRoleIdByGroupKey(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var pbRole *authpb.Role
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
	// Gọi repository để lấy admin info kèm role
	admins, err := h.roleRepo.GetAuthAdminsWithRolesByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}

	// Convert sang proto response
	var pbAdmins []*authpb.AuthAdminWithRole
	for _, admin := range admins {
		pbAdmins = append(pbAdmins, h.roleMapper.MapAuthAdminWithRoleToPb(admin))
	}

	return &authpb.GetAuthAdminByIdsResponse{
		Data: pbAdmins,
	}, nil
}
