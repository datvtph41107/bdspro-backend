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
	crud.BaseRepo[domain.RoleGroup]
}

func NewRoleGroupRepo(db *_db.DBTransaction) repo.RoleGroupRepository {
	return crud.MakeCrudRepo(&RoleGroupRepo{}, db).(*RoleGroupRepo)
}

func (r *RoleGroupRepo) GetByKey(ctx context.Context, key enums.GroupRoleKeyEnum) (*domain.RoleGroup, error) {
	var roleGroup domain.RoleGroup
	err := r.GetDB(ctx).Where("key = ?", key).First(&roleGroup).Error
	if err != nil {
		return nil, err
	}
	return &roleGroup, nil
}
