package postgres

import (
	"bdspro/internal/domain"
	property_repo "bdspro/internal/repo/property"
	"context"

	"gorm.io/gorm"
)

type AmenityRepo struct {
	DB *gorm.DB
}

func NewAmenityRepo(db *gorm.DB) property_repo.PropertyAmenityRepository {
	return &AmenityRepo{DB: db}
}

// ReplaceAmenities: DELETE cũ → INSERT mới trong cùng transaction.
// Dùng GORM Unscoped để xóa hard (bảng liên kết không có soft delete).
func (r *AmenityRepo) ReplaceAmenities(ctx context.Context, lineageID uint64, amenityIDs []uint64) error {
	db := GetDB(ctx, r.DB)

	// 1. Xóa toàn bộ liên kết cũ
	if err := db.
		Where("property_lineage_id = ?", lineageID).
		Delete(&domain.PropertyAmenity{}).Error; err != nil {
		return err
	}

	// 2. Insert list mới (nếu rỗng thì bỏ qua)
	if len(amenityIDs) == 0 {
		return nil
	}

	rows := make([]domain.PropertyAmenity, 0, len(amenityIDs))
	for _, id := range amenityIDs {
		rows = append(rows, domain.PropertyAmenity{
			PropertyLineageID: lineageID,
			AmenityItemID:     id,
		})
	}

	// CreateInBatches để tránh query quá dài khi list lớn
	return db.CreateInBatches(rows, 100).Error
}

func (r *AmenityRepo) GetAmenityIDsByLineageID(ctx context.Context, lineageID uint64) ([]uint64, error) {
	var rows []domain.PropertyAmenity
	err := GetDB(ctx, r.DB).
		Where("property_lineage_id = ?", lineageID).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.AmenityItemID)
	}
	return ids, nil
}
