package postgre

import (
	"context"
	"social/internal/domain"

	"gorm.io/gorm"
)

// @bind: social/internal/repo.NewsFeedShareRepo
type PostgreNewsFeedShare struct {
	db *gorm.DB
}

func NewPostgreNewsFeedShare(db *gorm.DB) *PostgreNewsFeedShare {
	return &PostgreNewsFeedShare{db: db}
}

func (r *PostgreNewsFeedShare) Create(ctx context.Context, newsFeedShare *domain.NewsFeedShare) (*domain.NewsFeedShare, error) {
	err := r.db.WithContext(ctx).Create(newsFeedShare).Error
	if err != nil {
		return nil, err
	}
	return newsFeedShare, nil
}
