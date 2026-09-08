package postgre

import (
	"context"
	"social/internal/domain"
	"social/internal/repo"

	"gorm.io/gorm"
)

type NewsFeedOfGroupPostgres struct {
	db *gorm.DB
}

func NewNewsFeedOfGroupPostgres(db *gorm.DB) repo.NewsFeedOfGroupRepo {
	return &NewsFeedOfGroupPostgres{db: db}
}

func (r *NewsFeedOfGroupPostgres) Create(ctx context.Context, newsFeedOfGroup *domain.NewsFeedOfGroup) error {
	return r.db.WithContext(ctx).Create(newsFeedOfGroup).Error
}

func (r *NewsFeedOfGroupPostgres) GetByGroupID(ctx context.Context, groupID uint64, page, size int) ([]domain.NewsFeedOfGroup, int64, error) {
	var newsFeeds []domain.NewsFeedOfGroup
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.NewsFeedOfGroup{}).Where("group_id = ?", groupID)

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

func (r *NewsFeedOfGroupPostgres) GetByPostID(ctx context.Context, postID uint64) (*domain.NewsFeedOfGroup, error) {
	var newsFeed domain.NewsFeedOfGroup
	if err := r.db.WithContext(ctx).Where("post_id = ?", postID).First(&newsFeed).Error; err != nil {
		return nil, err
	}
	return &newsFeed, nil
}

func (r *NewsFeedOfGroupPostgres) DeleteByPostID(ctx context.Context, postID uint64) error {
	return r.db.WithContext(ctx).Where("post_id = ?", postID).Delete(&domain.NewsFeedOfGroup{}).Error
}
