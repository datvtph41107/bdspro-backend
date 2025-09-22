package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.PostMediaRepo
type GormPostMediaRepo struct {
	DB *gorm.DB
}

func NewGormPostMediaRepo(db *gorm.DB) *GormPostMediaRepo {
	return &GormPostMediaRepo{
		DB: db,
	}
}

func (r *GormPostMediaRepo) UpdateMediaItems(c context.Context, postID uint64, mediaItems []dto.PostMediaItem) error {
	// Bắt đầu transaction để đảm bảo toàn vẹn dữ liệu
	// tx := r.DB.WithContext(c).Begin()
	tx := GetDB(c, r.DB)

	// Xóa các bản ghi cũ trong bảng
	if err := tx.Where("post_id = ?", postID).Delete(&domain.PostMediaEntity{}).Error; err != nil {
		// tx.Rollback()
		return err
	}

	// Nếu danh sách amenityIDs rỗng -> không cần thêm mới, chỉ cần xóa là đủ
	if len(mediaItems) == 0 {
		// tx.Commit()
		return nil
	}

	// Tạo danh sách mới để chèn vào
	var postMedia []domain.PostMediaEntity
	for idx, mediaItem := range mediaItems {
		postMedia = append(postMedia, domain.PostMediaEntity{
			PostID:    postID,
			MediaURL:  mediaItem.MediaURL,
			MediaType: mediaItem.MediaType,
			SortOrder: idx + 1,
			IsMain:    mediaItem.IsMain,
		})
	}

	// Chèn danh sách mới vào bảng nối
	if err := tx.Create(&postMedia).Error; err != nil {
		// tx.Rollback()
		return err
	}

	// Commit transaction
	// tx.Commit()
	return nil
}

func (r *GormPostMediaRepo) GetMediaList(c context.Context, dto dto.PostMediaGetDTO) ([]domain.PostMedia, error) {
	var postMedia []domain.PostMedia
	if err := r.DB.Where("deleted_at is null and post_id = ?", dto.PostID).Find(&postMedia).Error; err != nil {
		return nil, err
	}
	return postMedia, nil
}
