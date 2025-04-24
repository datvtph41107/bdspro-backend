package product_media

import (
	"bdspro/v1/data"
	"bdspro/v1/models"

	"gorm.io/gorm"
)

type ProductMediaRepo struct {
	DB *gorm.DB
}

func NewProductMediaRepo(db *gorm.DB) *ProductMediaRepo {
	return &ProductMediaRepo{
		DB: db,
	}
}

// ExistsByUserAndProduct kiểm tra xem UserId và ProductId có tồn tại không
func (repo *ProductMediaRepo) ExistsByUserAndProduct(userId, productId uint64) (bool, error) {
	var count int64
	err := repo.DB.Model(&models.ProductMedia{}).
		Where("created_by = ? AND product_id = ? and deleted_at is null", userId, productId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *ProductMediaRepo) Create(access *models.ProductMedia) error {
	return repo.DB.Create(&access).Error
}

func (r *ProductMediaRepo) UpdateProductMediaItems(productID uint64, mediaItems []data.MediaItem) error {
	// Bắt đầu transaction để đảm bảo toàn vẹn dữ liệu
	tx := r.DB.Begin()

	// Xóa các bản ghi cũ trong bảng
	if err := tx.Where("product_id = ?", productID).Delete(&models.ProductMedia{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Nếu danh sách amenityIDs rỗng -> không cần thêm mới, chỉ cần xóa là đủ
	if len(mediaItems) == 0 {
		tx.Commit()
		return nil
	}

	// Tạo danh sách mới để chèn vào
	var productAmenities []models.ProductMedia
	for idx, mediaItem := range mediaItems {
		productAmenities = append(productAmenities, models.ProductMedia{
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
