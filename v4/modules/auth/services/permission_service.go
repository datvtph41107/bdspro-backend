package services

import (
	"context"
	"time"

	_utils "common/utils"
	"userv4/internal/usecase"
	"userv4/modules/auth/dto"
)

// PermissionService owns caller permission-key reads migrated into modules/auth.
type PermissionService interface {
	GetMyPermissionKeys(ctx context.Context) (*dto.PermissionKeysResponse, error)
}

// RoleUsecasePermissionService adapts legacy RoleUsecase into modules/auth PermissionService.
// Permission ownership lives here; RoleUsecase remains the lookup backend until
// role_profiles / org-member resolution moves into this package.
type RoleUsecasePermissionService struct {
	role *usecase.RoleUsecase
	sync *_utils.SyncUtil
}

func NewRoleUsecasePermissionService(role *usecase.RoleUsecase, sync *_utils.SyncUtil) *RoleUsecasePermissionService {
	return &RoleUsecasePermissionService{role: role, sync: sync}
}

func (s *RoleUsecasePermissionService) GetMyPermissionKeys(ctx context.Context) (*dto.PermissionKeysResponse, error) {
	keys, err := s.role.GetPermissionKeys(ctx)
	if err != nil {
		return nil, err
	}
	now := time.Now().UnixMilli()
	if s.sync != nil {
		_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyPermissionMe, 0), now)
	}
	if keys == nil {
		keys = []string{}
	}
	return &dto.PermissionKeysResponse{
		Keys:      keys,
		Timestamp: now,
	}, nil
}
