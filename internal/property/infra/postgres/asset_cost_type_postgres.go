package postgres

import (
	"bdspro/infra/shared"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.AssetCostTypeRepo
type PostgreUserCostType struct {
	shared.PostgreCrud[domain.AssetCostType, *dto.AssetCostTypeSearchDTO]
}

func NewPostgreUserCostType(DB *gorm.DB) *PostgreUserCostType {
	repo := &PostgreUserCostType{
		PostgreCrud: shared.PostgreCrud[domain.AssetCostType, *dto.AssetCostTypeSearchDTO]{
			DB: DB,
		},
	}
	repo.Repo = repo
	return repo
}

func (r *PostgreUserCostType) QuerySearch(c context.Context, profileId uint64, ownerType enums.EOwnerOf, dto *dto.AssetCostTypeSearchDTO) (*gorm.DB, error) {
	query := r.DB.Model(&domain.AssetCostType{}).
		Where("deleted_at IS NULL and ((owner_id = ? and owner_type = ?) or owner_id is null)", profileId, enums.EOwnerOfMember)
	if dto.Text != "" {
		query = query.Where("type_name like ?", "%"+dto.Text+"%")
	}
	if dto.Type != nil {
		query = query.Where("type = ?", dto.Type)
	}
	query = query.Order("is_custom ASC")
	return query, nil
}
