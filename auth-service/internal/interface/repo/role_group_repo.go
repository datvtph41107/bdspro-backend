package repo

import (
	"auth/internal/domain"
	"auth/internal/enums"
	"common/case/crud"
	"context"
)

type RoleGroupRepository interface {
	crud.IBaseRepo[domain.RoleGroup]
	GetByKey(ctx context.Context, key enums.GroupRoleKeyEnum) (*domain.RoleGroup, error)
}
