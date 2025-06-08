package postgre

import (
	"context"
	"social/internal/domain"

	"gorm.io/gorm"
)

// @bind: social/internal/repo.FriendTagRepo
type PostgreFriendTag struct {
	db *gorm.DB
}

func NewPostgreFriendTag(db *gorm.DB) *PostgreFriendTag {
	return &PostgreFriendTag{db: db}
}

func (r *PostgreFriendTag) UpdateFriendTag(ctx context.Context, newsFeedID uint64, friendTagIds []uint64) error {
	return r.db.WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			tx = GetDB(ctx, tx)
			// 1. Xoá hết ảnh cũ
			if err := tx.
				Where("news_feed_id = ?", newsFeedID).
				Delete(&domain.FriendTag{}).
				Error; err != nil {
				return err
			}

			tags := make([]*domain.FriendTag, len(friendTagIds))
			for i, friendTagId := range friendTagIds {
				tags[i] = &domain.FriendTag{
					NewsFeedID: newsFeedID,
					UserID:     friendTagId,
				}
			}

			// 3. Tạo mới toàn bộ ảnh
			if len(tags) > 0 {
				if err := tx.Create(&tags).Error; err != nil {
					return err
				}
			}

			return nil
		})
}
