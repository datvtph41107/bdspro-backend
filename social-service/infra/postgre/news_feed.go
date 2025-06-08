package postgre

import (
	"context"
	"social/internal/domain"
	"social/internal/dto"
	"time"

	"gorm.io/gorm"
)

// @bind: social/internal/repo.NewsFeedRepo
type PostgreNewsFeed struct {
	db *gorm.DB
}

func NewPostgreNewsFeed(db *gorm.DB) *PostgreNewsFeed {
	return &PostgreNewsFeed{db: db}
}

func (r *PostgreNewsFeed) Create(ctx context.Context, post *domain.NewsFeed) (*domain.NewsFeed, error) {
	err := r.db.WithContext(ctx).Create(post).Error
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *PostgreNewsFeed) Update(ctx context.Context, post *domain.NewsFeed) (*domain.NewsFeed, error) {
	err := r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("id = ?", post.ID).
		Update("title", post.Title).
		Update("content", post.Content).
		// Update("description", post.Description).
		Update("image", post.Image).
		Update("link", post.Link).
		Update("visibility", post.Visibility).Error
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *PostgreNewsFeed) Delete(ctx context.Context, post *domain.NewsFeed) error {
	err := r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("id = ?", post.ID).
		Update("deleted_at", time.Now()).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgreNewsFeed) GetGlobal(ctx context.Context, post *dto.NewsFeedGlobalSearch) ([]*domain.NewsFeed, int32, error) {
	var newsFeeds []*domain.NewsFeed
	var total int64
	query := r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("deleted_at IS NULL and visibility = 10")

	if post.Text != "" {
		query = query.Where("title LIKE ?", "%"+post.Text+"%")
	}

	err := query.
		Order("created_at DESC").
		Offset(post.GetOffset()).
		Limit(post.GetLimit()).
		Find(&newsFeeds).Error
	if err != nil {
		return nil, 0, err
	}
	query.Count(&total)
	return newsFeeds, int32(total), nil
}

func (r *PostgreNewsFeed) UpdateLike(ctx context.Context, newFeedID uint64, numLike int) error {
	return r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("id = ?", newFeedID).
		Update("like_number", numLike).Error
}

func (r *PostgreNewsFeed) UpdateComment(ctx context.Context, newFeedID uint64, numComment int) error {
	return nil
}

func (r *PostgreNewsFeed) UpdateShare(ctx context.Context, newFeedID uint64, numShare int) error {
	return nil
}

func (r *PostgreNewsFeed) UpdateView(ctx context.Context, newFeedID uint64, numView int) error {
	return nil
}

func (r *PostgreNewsFeed) UserAvaiableReact(ctx context.Context, profileID uint64, id uint64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("id = ? AND visibility = 10 AND deleted_at IS NULL", id).
		Count(&count).
		Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgreNewsFeed) ExistByID(ctx context.Context, id uint64) (bool, error) {
	var exists bool
	query := `
        SELECT EXISTS (
            SELECT 1 FROM news_feeds
            WHERE id = $1 AND deleted_at IS NULL
        )
    `
	err := r.db.WithContext(ctx).
		Raw(query, id).
		Scan(&exists).
		Error
	return exists, err
}

func (r *PostgreNewsFeed) GetByID(ctx context.Context, id uint64) (*domain.NewsFeed, error) {
	var newsFeed *domain.NewsFeed
	err := r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&newsFeed).
		Error
	return newsFeed, err
}
