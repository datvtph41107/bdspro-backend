package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"

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

func (r *GormProductMediaRepo) UpdateProductMediaItems(ctx context.Context, productID uint64, mediaItems []dto.MediaItem) ([]uint64, error) {
	// Bắt đầu transaction để đảm bảo toàn vẹn dữ liệu
	tx := GetDB(ctx, r.DB)

	// Xóa các bản ghi cũ trong bảng
	if err := tx.Where("product_id = ?", productID).Delete(&domain.ProductMedia{}).Error; err != nil {
		// tx.Rollback()
		return nil, err
	}

	// Nếu danh sách amenityIDs rỗng -> không cần thêm mới, chỉ cần xóa là đủ
	if len(mediaItems) == 0 {
		// tx.Commit()
		return nil, nil
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
		// tx.Rollback()
		return nil, err
	}

	// Lấy danh sách ID đã tạo
	var mediaIDs []uint64
	for _, media := range productAmenities {
		mediaIDs = append(mediaIDs, media.ID)
	}

	// Commit transaction
	// tx.Commit()
	return mediaIDs, nil
}

func (r *GormProductMediaRepo) GetByProductID(ctx context.Context, productID uint64) ([]domain.ProductMedia, error) {
	var mediaList []domain.ProductMedia
	tx := GetDB(ctx, r.DB)
	err := tx.Where("product_id = ? AND deleted_at IS NULL", productID).
		Order("sort_order ASC, created_at ASC").
		Find(&mediaList).Error
	if err != nil {
		return nil, err
	}
	return mediaList, nil
}
