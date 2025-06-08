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
	err := GetDB(ctx, r.db).Create(post).Error
	return post, err
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

func (r *PostgreNewsFeed) QuerySeach(ctx context.Context, dto *dto.NewsFeedGlobalSearch) (*gorm.DB, error) {
	query := r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("deleted_at IS NULL and visibility = 10")

	if dto.Text != "" {
		query = query.Where("title LIKE ?", "%"+dto.Text+"%")
	}

	return query, nil
}

func (r *PostgreNewsFeed) GetGlobal(ctx context.Context, post *dto.NewsFeedGlobalSearch) ([]*domain.NewsFeed, int32, error) {
	var newsFeeds []*domain.NewsFeed
	var total int64
	query, err := r.QuerySeach(ctx, post)
	if err != nil {
		return nil, 0, err
	}
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Order("created_at DESC").
		Offset(post.GetOffset()).
		Limit(post.GetLimit()).
		Find(&newsFeeds).Error
	if err != nil {
		return nil, 0, err
	}

	return newsFeeds, int32(total), nil
}

func (r *PostgreNewsFeed) UpdateLikeNumber(ctx context.Context, newFeedID uint64, totalLike int) error {
	return r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("id = ? and deleted_at is null", newFeedID).
		Update("num_like", totalLike).Error
}

func (r *PostgreNewsFeed) UpdateCommentNumber(ctx context.Context, newFeedID uint64, totalComment int) error {
	return r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("id = ? and deleted_at is null", newFeedID).
		Update("num_comment", totalComment).Error
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

func (r *PostgreNewsFeed) GetLastUpdated(ctx context.Context, id uint64) (*time.Time, error) {
	var updatedTime time.Time
	err := r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Select("updated_at").
		Where("id = ? AND deleted_at IS NULL", id).
		Order("updated_at DESC").
		Limit(1).
		Pluck("updated_at", &updatedTime).Error

	if err != nil {
		return nil, err
	}

	return &updatedTime, nil
}

func (r *PostgreNewsFeed) IncrementNumberShare(ctx context.Context, newsFeedID uint64) error {
	return r.db.WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&domain.NewsFeed{}).
				Where("id = ? AND deleted_at IS NULL", newsFeedID).
				UpdateColumn("num_share", gorm.Expr("COALESCE(num_share, 0) + ?", 1)).Error
			if err != nil {
				return err
			}
			return nil
		})
}
