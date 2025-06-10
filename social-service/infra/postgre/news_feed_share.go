package postgre

import (
	"context"
	"social/internal/domain"
	"social/internal/enums"

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

func (r *PostgreNewsFeedShare) GetNearestUserIds(ctx context.Context, newsFeedId uint64, targetType enums.TargetType) ([]uint64, error) {
	var ids []uint64
	err := r.db.WithContext(ctx).Debug().Raw(
		`
		select user_id 
		from (
			SELECT  c.user_id, MAX(created_at) as latest
				FROM news_feed_shares c
				WHERE news_feed_id = ? AND deleted_at IS NULL
				GROUP BY c.user_id
				ORDER BY latest DESC
			)
		Limit 3`, newsFeedId).Scan(&ids).Error

	return ids, err
}
func (r *PostgreNewsFeedShare) SumUserComment(ctx context.Context, newsFeedId uint64, targetType enums.TargetType) (int, error) {
	var count int
	err := r.db.WithContext(ctx).Raw(
		`
		select count(*) from (SELECT  c.user_id, MAX(created_at) as latest
		FROM news_feed_shares c
		WHERE news_feed_id = ? AND deleted_at IS NULL
		group by c.user_id
		)
	`, newsFeedId).Scan(&count).Error

	return count, err
}
