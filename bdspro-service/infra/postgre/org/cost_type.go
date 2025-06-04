package org_postgre

import (
	"bdspro/infra/shared"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo/org.OrgCostTypeRepo
type PostgreOrgCostType struct {
	shared.PostgreCrud[domain.AssetCostType, *dto.AssetCostTypeSearchDTO]
}

func NewPostgreOrgCostType(db *gorm.DB) *PostgreOrgCostType {
	repo := &PostgreOrgCostType{
		PostgreCrud: shared.PostgreCrud[domain.AssetCostType, *dto.AssetCostTypeSearchDTO]{
			DB: db,
		},
	}
	repo.Repo = repo
	return repo
}

func (r *PostgreOrgCostType) QuerySearch(c context.Context, profileId uint64, dto *dto.AssetCostTypeSearchDTO) (*gorm.DB, error) {
	query := r.DB.Model(&domain.AssetCostType{}).
		Where("deleted_at IS NULL and ((owner_id = ? and owner_type = ?) or owner_id is null)", profileId, enums.EOwnerTypeOrgnization)
	if dto.Text != "" {
		query = query.Where("type_name like ?", "%"+dto.Text+"%")
	}
	if dto.Types != nil {
		query = query.Where("type in (?)", dto.Types)
	}
	return query, nil
}
