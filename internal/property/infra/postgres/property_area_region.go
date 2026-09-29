package postgres

import (
	"bdspro/internal/domain"
	property_repo "bdspro/internal/repo/property"
	"context"

	"gorm.io/gorm"
)

type AreaRegionRepo struct {
	DB *gorm.DB
}

func NewAreaRegionRepo(db *gorm.DB) property_repo.AreaRegionRepository {
	return &AreaRegionRepo{DB: db}
}

func (r *AreaRegionRepo) ReplaceAreaRegions(ctx context.Context, lineageID uint64, regionIDs []uint64) error {
	db := GetDB(ctx, r.DB)

	if err := db.
		Where("property_lineage_id = ?", lineageID).
		Delete(&domain.PropertyAreaRegion{}).Error; err != nil {
		return err
	}

	if len(regionIDs) == 0 {
		return nil
	}

	rows := make([]domain.PropertyAreaRegion, 0, len(regionIDs))
	for _, id := range regionIDs {
		rows = append(rows, domain.PropertyAreaRegion{
			PropertyLineageID: lineageID,
			AreaRegionID:      id,
		})
	}

	// CreateInBatches để tránh query quá dài khi list lớn
	return db.CreateInBatches(rows, 100).Error
}

func (r *AreaRegionRepo) GetAreaRegionIDsByLineageID(ctx context.Context, lineageID uint64) ([]uint64, error) {
	var rows []domain.PropertyAreaRegion
	err := GetDB(ctx, r.DB).
		Where("property_id = ?", lineageID).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.AreaRegionID)
	}
	return ids, nil
}
