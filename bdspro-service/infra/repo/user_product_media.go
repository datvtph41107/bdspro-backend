package postgre_repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.ProductMediaRepo
type GormProductMediaRepo struct {
	DB *gorm.DB
}

func NewGormProductMediaRepo(db *gorm.DB) *GormProductMediaRepo {
	return &GormProductMediaRepo{
		DB: db,
	}
}

// ExistsByUserAndProduct kiểm tra xem UserId và ProductId có tồn tại không
func (repo *GormProductMediaRepo) ExistsByUserAndProduct(userId, productId uint64) (bool, error) {
	var count int64
	err := repo.DB.Model(&domain.ProductMedia{}).
		Where("created_by = ? AND product_id = ? and deleted_at is null", userId, productId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *GormProductMediaRepo) Create(access *domain.ProductMedia) error {
	return repo.DB.Create(&access).Error
}

func (r *GormProductMediaRepo) UpdateProductMediaItems(productID uint64, mediaItems []dto.MediaItem) error {
	// Bắt đầu transaction để đảm bảo toàn vẹn dữ liệu
	tx := r.DB.Begin()

	// Xóa các bản ghi cũ trong bảng
	if err := tx.Where("product_id = ?", productID).Delete(&domain.ProductMedia{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Nếu danh sách amenityIDs rỗng -> không cần thêm mới, chỉ cần xóa là đủ
	if len(mediaItems) == 0 {
		tx.Commit()
		return nil
	}

	// Tạo danh sách mới để chèn vào
	var productAmenities []domain.ProductMedia
	for idx, mediaItem := range mediaItems {
		productAmenities = append(productAmenities, domain.ProductMedia{
			ProductID: productID,
			MediaURL:  mediaItem.MediaURL,
			MediaType: mediaItem.MediaType,
			Order:     idx + 1,
			IsMain:    mediaItem.IsMain,
		})
	}

	// Chèn danh sách mới vào bảng nối
	if err := tx.Create(&productAmenities).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	tx.Commit()
	return nil
}
