package postgre

import (
	"context"
	"social/internal/domain"
	"social/internal/repo"

	"gorm.io/gorm"
)

type NewsFeedOfUserPostgres struct {
	db *gorm.DB
}

func NewNewsFeedOfUserPostgres(db *gorm.DB) repo.NewsFeedOfUserRepo {
	return &NewsFeedOfUserPostgres{db: db}
}

func (r *NewsFeedOfUserPostgres) Create(ctx context.Context, newsFeedOfUser *domain.NewsFeedOfUser) error {
	return r.db.WithContext(ctx).Create(newsFeedOfUser).Error
}

func (r *NewsFeedOfUserPostgres) GetByUserID(ctx context.Context, userID uint64, page, size int) ([]domain.NewsFeedOfUser, int64, error) {
	var newsFeeds []domain.NewsFeedOfUser
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.NewsFeedOfUser{}).Where("user_id = ?", userID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&newsFeeds).Error; err != nil {
		return nil, 0, err
	}

	return newsFeeds, total, nil
}

func (r *NewsFeedOfUserPostgres) GetByPostID(ctx context.Context, postID uint64) (*domain.NewsFeedOfUser, error) {
	var newsFeed domain.NewsFeedOfUser
	if err := r.db.WithContext(ctx).Where("post_id = ?", postID).First(&newsFeed).Error; err != nil {
		return nil, err
	}
	return &newsFeed, nil
}

func (r *NewsFeedOfUserPostgres) DeleteByPostID(ctx context.Context, postID uint64) error {
	return r.db.WithContext(ctx).Where("post_id = ?", postID).Delete(&domain.NewsFeedOfUser{}).Error
}
