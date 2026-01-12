package postgres

import (
	"bdspro/internal/domain"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.PropertyAmenityRepo
type PostgrePropertyAmenityRepo struct {
	DB *gorm.DB
}

func NewPostgrePropertyAmenityRepo(db *gorm.DB) *PostgrePropertyAmenityRepo {
	return &PostgrePropertyAmenityRepo{
		DB: db,
	}
}

func (r *PostgrePropertyAmenityRepo) UpdatePropertyAmenityIds(ctx context.Context, propertyID uint64, amenityIDs []uint64) error {
	// Bắt đầu transaction để đảm bảo toàn vẹn dữ liệu
	tx := GetDB(ctx, r.DB)

	// Xóa các bản ghi cũ trong bảng nối
	if err := tx.Where("property_id = ?", propertyID).Delete(&domain.PropertyAmenity{}).Error; err != nil {
		return err
	}

	// Nếu danh sách amenityIDs rỗng -> không cần thêm mới, chỉ cần xóa là đủ
	if len(amenityIDs) == 0 {
		return nil
	}

	// Tạo danh sách mới để chèn vào
	var propertyAmenities []domain.PropertyAmenity
	for _, amenityID := range amenityIDs {
		propertyAmenities = append(propertyAmenities, domain.PropertyAmenity{
			PropertyID: propertyID,
			AmenityID:  amenityID,
		})
	}

	// Chèn danh sách mới vào bảng nối
	if err := tx.Create(&propertyAmenities).Error; err != nil {
		return err
	}

	return nil
}
