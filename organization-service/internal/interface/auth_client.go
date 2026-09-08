package iusecase

import (
	"context"
	"organization/internal/dto"
	authpb "pb/types/auth"
)

type IAuthClient interface {
	MapRoleToGroups(ctx context.Context, groups []*dto.GroupWithDetails)
	GetAdminRoleIdByGroupKey(ctx context.Context, groupKey int) (*authpb.Role, error)
}
