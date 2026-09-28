package services

import (
	"context"

	"userv4/modules/auth/dto"
)

// PermissionService owns caller permission-key reads migrated into modules/auth.
type PermissionService interface {
	GetMyPermissionKeys(ctx context.Context) (*dto.PermissionKeysResponse, error)
}
