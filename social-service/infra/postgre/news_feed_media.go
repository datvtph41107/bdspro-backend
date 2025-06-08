package postgre

import (
	"context"
	"social/internal/domain"
	"time"

	"gorm.io/gorm"
)

// @bind: social/internal/repo.NewsFeedMediaRepo
type PostgreNewsFeedMedia struct {
	db *gorm.DB
}

func NewPostgreNewsFeedMedia(db *gorm.DB) *PostgreNewsFeedMedia {
	return &PostgreNewsFeedMedia{db: db}
}

func (r *PostgreNewsFeedMedia) UpdateMedia(ctx context.Context, newsFeedID uint64, medias []*domain.NewsFeedMedia) error {
	return r.db.WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			tx = GetDB(ctx, tx)
			// 1. Xoá hết ảnh cũ
			if err := tx.Model(&domain.NewsFeedMedia{}).
				Where("news_feed_id = ?", newsFeedID).
				Update("deleted_at", time.Now()).
				Error; err != nil {
				return err
			}

			// 2. Gán post_id cho tất cả ảnh mới (nếu cần)
			for _, media := range medias {
				media.NewsFeedID = &newsFeedID
			}

			// 3. Tạo mới toàn bộ ảnh
			if len(medias) > 0 {
				if err := tx.Create(&medias).Error; err != nil {
					return err
				}
			}

			return nil
		})
}
