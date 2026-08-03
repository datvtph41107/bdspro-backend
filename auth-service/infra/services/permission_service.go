package services

import (
	"context"
	"log"
	"sync"
	"time"

	_algorithm "common/algorithm"
	_enum "common/domain/enum"
	_errors "common/errors"
	_utils "common/utils"
	"pb/clients"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
)

type PermissionService struct {
	mu              sync.RWMutex
	rolePermissions map[uint64]*_algorithm.BitCache
	userClient      *clients.UserGrpcClient
}

func NewPermissionService(envRuntime string) *PermissionService {
	svc := &PermissionService{
		rolePermissions: make(map[uint64]*_algorithm.BitCache),
		userClient:      clients.MakeUserGrpcClient(envRuntime),
	}
	go svc.loadRolePermissions()
	return svc
}

func (s *PermissionService) loadRolePermissions() {
	for {
		if s.userClient == nil || s.userClient.InternalClient == nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		resp, err := s.userClient.GetRolePermissions(ctx, &sharepb.IdRequest{})
		cancel()
		if err != nil {
			log.Printf("[PermissionService] load role permissions failed: %v, retrying...", err)
			time.Sleep(10 * time.Second)
			continue
		}

		s.setRolePermissions(resp.GetData())
		log.Printf("[PermissionService] loaded permissions for %d roles", len(resp.GetData()))
		return
	}
}

func (s *PermissionService) setRolePermissions(items []*userpb.RolePermissionItem) {
	updated := make(map[uint64]*_algorithm.BitCache, len(items))

	for _, item := range items {
		if item == nil || item.Role == nil || item.Role.Id == 0 {
			continue
		}

		bitmap := _algorithm.NewBitCache(_enum.PermissionIDMax)
		for _, permID := range item.PermissionIds {
			bitmap.Set(int(permID), true)
		}
		updated[item.Role.Id] = bitmap
	}

	s.mu.Lock()
	s.rolePermissions = updated
	s.mu.Unlock()
}

func (s *PermissionService) HasPermissions(roleID uint64, permissionIds []uint32) bool {
	if len(permissionIds) == 0 {
		return true
	}

	log.Printf("[PermissionService] HasPermissions: roleID: %d, permissionIds: %v", roleID, permissionIds)
	s.mu.RLock()
	bitmap, ok := s.rolePermissions[roleID]
	s.mu.RUnlock()
	if !ok || bitmap == nil {
		return false
	}

	for _, permID := range permissionIds {
		if bitmap.Get(int(permID)) {
			return true
		}
	}
	return false
}

func (s *PermissionService) RequiredPermissions(ctx context.Context, permissionIds []uint32) error {
	if len(permissionIds) == 0 {
		return nil
	}

	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return _errors.ReturnError(401, "Unauthorized")
	}

	roleIds := _utils.GetRoleIdsFromContext(ctx)
	if len(roleIds) == 0 {
		return _errors.ReturnError(403, "Bạn không có quyền truy cập")
	}

	for _, roleId := range roleIds {
		if s.HasPermissions(roleId, permissionIds) {
			return nil
		}
	}

	return _errors.ReturnError(403, "Bạn không có quyền truy cập")
}
