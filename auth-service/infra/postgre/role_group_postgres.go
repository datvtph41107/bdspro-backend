package postgre

import (
	"auth/internal/domain"
	"auth/internal/enums"
	"auth/internal/interface/repo"
	"common/case/crud"
	_db "common/db"
	"context"
)

type RoleGroupRepo struct {
	*crud.CrudRepo[domain.RoleGroup]
}

func NewRoleGroupRepo(db *_db.TransactionRepo) repo.RoleGroupRepository {
	x := &RoleGroupRepo{
		CrudRepo: &crud.CrudRepo[domain.RoleGroup]{},
	}
	x.CrudRepo.Init(x, db)
	return x
}

func (r *RoleGroupRepo) GetByKey(ctx context.Context, key enums.GroupRoleKeyEnum) (*domain.RoleGroup, error) {
	var roleGroup domain.RoleGroup
	err := r.GetDB(ctx).Where("key = ?", key).First(&roleGroup).Error
	if err != nil {
		return nil, err
	}
	return &roleGroup, nil
}
