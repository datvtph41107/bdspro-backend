package postgres

import (
	"context"
	"time"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type regionLabelRepoImpl struct {
	db *gorm.DB
}

func NewRegionLabelRepository(db *gorm.DB) repo.RegionLabelRepository {
	return &regionLabelRepoImpl{db: db}
}

func (r *regionLabelRepoImpl) AssignLabels(ctx context.Context, regionID uint64, labelIDs []uint64, assignedBy uint64, sourceFile, notes string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Xóa tất cả labels cũ
		if err := tx.Where("region_id = ?", regionID).Delete(&qh_domain.QHRegionLabel{}).Error; err != nil {
			return err
		}

		// Thêm labels mới
		if len(labelIDs) == 0 {
			return nil
		}

		records := make([]qh_domain.QHRegionLabel, 0, len(labelIDs))
		now := time.Now()
		for _, labelID := range labelIDs {
			records = append(records, qh_domain.QHRegionLabel{
				RegionID:   regionID,
				LabelID:    labelID,
				AssignedBy: assignedBy,
				AssignedAt: now,
				SourceFile: sourceFile,
				Notes:      notes,
			})
		}
		return tx.CreateInBatches(records, 100).Error
	})
}

func (r *regionLabelRepoImpl) AddLabels(ctx context.Context, regionID uint64, labelIDs []uint64, assignedBy uint64, sourceFile, notes string) error {
	if len(labelIDs) == 0 {
		return nil
	}

	// Lọc các labels đã tồn tại để tránh duplicate
	var existingIDs []uint64
	if err := r.db.WithContext(ctx).
		Model(&qh_domain.QHRegionLabel{}).
		Where("region_id = ? AND label_id IN ?", regionID, labelIDs).
		Pluck("label_id", &existingIDs).Error; err != nil {
		return err
	}

	existingMap := make(map[uint64]bool)
	for _, id := range existingIDs {
		existingMap[id] = true
	}

	newLabelIDs := make([]uint64, 0)
	for _, id := range labelIDs {
		if !existingMap[id] {
			newLabelIDs = append(newLabelIDs, id)
		}
	}

	if len(newLabelIDs) == 0 {
		return nil
	}

	records := make([]qh_domain.QHRegionLabel, 0, len(newLabelIDs))
	now := time.Now()
	for _, labelID := range newLabelIDs {
		records = append(records, qh_domain.QHRegionLabel{
			RegionID:   regionID,
			LabelID:    labelID,
			AssignedBy: assignedBy,
			AssignedAt: now,
			SourceFile: sourceFile,
			Notes:      notes,
		})
	}
	return r.db.WithContext(ctx).CreateInBatches(records, 100).Error
}

func (r *regionLabelRepoImpl) RemoveLabel(ctx context.Context, regionID, labelID uint64) error {
	return r.db.WithContext(ctx).
		Where("region_id = ? AND label_id = ?", regionID, labelID).
		Delete(&qh_domain.QHRegionLabel{}).Error
}

func (r *regionLabelRepoImpl) BatchAssignLabels(ctx context.Context, regionIDs, labelIDs []uint64, assignedBy uint64, sourceFile, notes string) (int, error) {
	if len(regionIDs) == 0 || len(labelIDs) == 0 {
		return 0, nil
	}

	records := make([]qh_domain.QHRegionLabel, 0, len(regionIDs)*len(labelIDs))
	now := time.Now()
	for _, regionID := range regionIDs {
		for _, labelID := range labelIDs {
			records = append(records, qh_domain.QHRegionLabel{
				RegionID:   regionID,
				LabelID:    labelID,
				AssignedBy: assignedBy,
				AssignedAt: now,
				SourceFile: sourceFile,
				Notes:      notes,
			})
		}
	}

	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(records, 500).Error
	return len(records), err
}

func (r *regionLabelRepoImpl) GetLabelsByRegion(ctx context.Context, regionID uint64) ([]qh_domain.QHLabel, error) {
	var labels []qh_domain.QHLabel
	err := r.db.WithContext(ctx).
		Table("qh_labels").
		Joins("INNER JOIN qh_region_labels ON qh_region_labels.label_id = qh_labels.id").
		Where("qh_region_labels.region_id = ? AND qh_labels.deleted_at IS NULL", regionID).
		Order("qh_labels.display_order ASC, qh_labels.id ASC").
		Find(&labels).Error
	return labels, err
}

func (r *regionLabelRepoImpl) GetRegionIDsByLabel(ctx context.Context, labelID uint64) ([]uint64, error) {
	var regionIDs []uint64
	err := r.db.WithContext(ctx).
		Model(&qh_domain.QHRegionLabel{}).
		Where("label_id = ?", labelID).
		Pluck("region_id", &regionIDs).Error
	return regionIDs, err
}

func (r *regionLabelRepoImpl) DeleteByRegion(ctx context.Context, regionID uint64) error {
	return r.db.WithContext(ctx).
		Where("region_id = ?", regionID).
		Delete(&qh_domain.QHRegionLabel{}).Error
}

func (r *regionLabelRepoImpl) DeleteByLabel(ctx context.Context, labelID uint64) error {
	return r.db.WithContext(ctx).
		Where("label_id = ?", labelID).
		Delete(&qh_domain.QHRegionLabel{}).Error
}
